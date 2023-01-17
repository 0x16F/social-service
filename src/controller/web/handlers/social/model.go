package social

import (
	"core-server/src/controller/web/handlers/social/feed"
	"core-server/src/controller/web/handlers/social/message"
	"core-server/src/controller/web/handlers/social/profile"
	"core-server/src/controller/web/handlers/social/reply"
)

type Handler struct {
	Message message.Servicer
	Profile profile.Servicer
	Reply   reply.Servicer
	Feed    feed.Servicer
}
