package social

import (
	"core-server/src/internal/social/feed"
	"core-server/src/internal/social/like"
	"core-server/src/internal/social/message"
	"core-server/src/internal/social/profile"
	"core-server/src/internal/social/reply"
	"core-server/src/internal/social/subs"
	"database/sql"
)

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Subs:     subs.NewStorage(db),
		Profiles: profile.NewStorage(db),
		Messages: message.NewStorage(db),
		Likes:    like.NewStorage(db),
		Replies:  reply.NewStorage(db),
		Feed:     feed.NewStorage(db),
	}
}
