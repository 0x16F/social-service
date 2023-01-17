package social

import (
	"core-server/src/internal/social/feed"
	"core-server/src/internal/social/like"
	"core-server/src/internal/social/message"
	"core-server/src/internal/social/profile"
	"core-server/src/internal/social/reply"
	"core-server/src/internal/social/subs"
)

type Storage struct {
	Subs     subs.Storager
	Profiles profile.Storager
	Messages message.Storager
	Likes    like.Storager
	Replies  reply.Storager
	Feed     feed.Storager
}
