package handlers

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/handlers/authorization"
	"core-server/src/controller/web/handlers/invite"
	"core-server/src/controller/web/handlers/social"
	"core-server/src/controller/web/handlers/user"
	"core-server/src/pkg/jwt"
	"core-server/src/pkg/logger"
	"core-server/src/pkg/turnstile"
)

func NewHandler(l *logger.Logger, j jwt.Servicer, s *storage.Storage, c turnstile.Servicer) *Handlers {
	return &Handlers{
		Authorization: authorization.NewHandler(l, j, s, c),
		Invite:        invite.NewHandler(l, s),
		User:          user.NewHandler(l, s),
		Social:        social.NewHandler(l, s),
	}
}
