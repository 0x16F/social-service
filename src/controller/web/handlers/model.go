package handlers

import (
	"core-server/src/controller/web/handlers/authorization"
	"core-server/src/controller/web/handlers/invite"
	"core-server/src/controller/web/handlers/social"
	"core-server/src/controller/web/handlers/user"
)

type Handlers struct {
	Authorization authorization.Servicer
	Invite        invite.Servicer
	User          *user.Handler
	Social        *social.Handler
}
