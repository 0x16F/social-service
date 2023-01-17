package web

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/handlers"
	"core-server/src/pkg/jwt"
	"core-server/src/pkg/logger"
	"core-server/src/pkg/turnstile"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(l *logger.Logger, j jwt.Servicer, s *storage.Storage, c turnstile.Servicer) *Server {
	return &Server{
		Router: fiber.New(fiber.Config{
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
		}),
		Handlers: handlers.NewHandler(l, j, s, c),
		Logger:   l,
	}
}

func (s *Server) _ConfigureRoutes() {
	// Cors
	config := cors.ConfigDefault
	config.AllowCredentials = true

	s.Router.Use(cors.New(config))
	s.Router.Use(s.Handlers.Authorization.IsAuthorized)

	// Handlers
	v0 := s.Router.Group("/v0")
	{
		auth := v0.Group("/auth")
		{
			auth.Post("/login", s.Handlers.Authorization.Login)
			auth.Post("/register", s.Handlers.Authorization.Register)
			auth.Get("/refresh", s.Handlers.Authorization.Refresh)
			auth.Get("/logout", s.Handlers.Authorization.Logout)
			auth.Get("/parse/:token", s.Handlers.Authorization.ParseToken)
		}

		invite := v0.Group("/invite")
		{
			invite.Post("", s.Handlers.Invite.Create)
			invite.Get("/:invite", s.Handlers.Invite.Use)
		}

		user := v0.Group("/user")
		{
			user.Get("", s.Handlers.User.GetMyInfo)
			user.Put("", s.Handlers.User.ChangePassword)
			user.Get("/:id", s.Handlers.User.GetProfileInfo)
		}

		social := v0.Group("/social")
		{
			messages := social.Group("/message")
			{
				messages.Get("/:id", s.Handlers.Social.Message.GetAll)
				messages.Delete("/:id", s.Handlers.Social.Message.Delete)
				messages.Post("/:id", s.Handlers.Social.Message.Like)
				messages.Post("", s.Handlers.Social.Message.Create)

				replies := messages.Group("/reply")
				{
					replies.Post("", s.Handlers.Social.Reply.Create)
					replies.Get("/:id", s.Handlers.Social.Reply.GetAll)
					replies.Delete("/:id", s.Handlers.Social.Reply.Delete)
				}
			}

			profile := social.Group("/profile")
			{
				profile.Get("", s.Handlers.Social.Profile.GetAll)
				profile.Get("/subscribers/:id", s.Handlers.Social.Profile.GetSubscribers)
				profile.Get("/subscriptions/:id", s.Handlers.Social.Profile.GetSubscriptions)
				profile.Post("/subscribe/:id", s.Handlers.User.Subscribe)

				avatar := profile.Group("/avatar")
				{
					avatar.Get("/:id", s.Handlers.User.GetUserAvatar)
					avatar.Delete("", s.Handlers.User.DeleteUserAvatar)
					avatar.Post("", s.Handlers.User.SetUserAvatar)
				}

				background := profile.Group("/background")
				{
					background.Get("/:id", s.Handlers.User.GetUserBackground)
					background.Delete("", s.Handlers.User.DeleteUserBackground)
					background.Post("", s.Handlers.User.SetUserBackground)
				}
			}

			feed := social.Group("/feed")
			{
				feed.Get("", s.Handlers.Social.Feed.Get)
			}
		}
	}
}

func (s *Server) Start(port int) {
	s._ConfigureRoutes()
	s.Logger.Error(s.Router.Listen(fmt.Sprintf(":%d", port)))
}
