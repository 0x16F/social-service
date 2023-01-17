package invites

import (
	"core-server/src/pkg/crypto"
	"time"
)

func NewInvite(authorId int) *Invite {
	return &Invite{
		Invite:     crypto.GenerateString(5),
		AuthorId:   authorId,
		CreateTime: time.Now().UnixMilli(),
	}
}
