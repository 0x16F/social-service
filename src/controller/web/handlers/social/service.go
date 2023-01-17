package social

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/handlers/social/feed"
	"core-server/src/controller/web/handlers/social/message"
	"core-server/src/controller/web/handlers/social/profile"
	"core-server/src/controller/web/handlers/social/reply"
	"core-server/src/pkg/logger"
)

func NewHandler(l *logger.Logger, s *storage.Storage) *Handler {
	return &Handler{
		Message: message.NewHandler(l, s),
		Profile: profile.NewHandler(l, s),
		Reply:   reply.NewHandler(l, s),
		Feed:    feed.NewHandler(l, s),
	}
}
