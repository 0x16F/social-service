package authorization

import (
	"core-server/src/controller/storage"
	"core-server/src/pkg/jwt"
	"core-server/src/pkg/logger"
	"core-server/src/pkg/turnstile"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Logger  *logger.Logger
	Storage *storage.Storage
	JWT     jwt.Servicer
	CF      turnstile.Servicer
}

type Servicer interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	ParseToken(c *fiber.Ctx) error
	IsAuthorized(c *fiber.Ctx) error
}

type Whitelist struct {
	Path    string
	Methods []string
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Solution string `json:"solution"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Solution string `json:"solution"`
}
