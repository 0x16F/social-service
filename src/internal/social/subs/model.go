package subs

import (
	"core-server/src/internal/social/profile"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

type SocialSubscribers struct {
	Subscriber int `json:"subscriber_id"`
	UserId     int `json:"user_id"`
}

type Storager interface {
	GetSubscriptions(wall, offset, limit int) (*[]profile.ProfilePreview, error)
	GetSubscribers(wall, offset, limit int) (*[]profile.ProfilePreview, error)
	GetSubscribersLike(name string, wall, offset, limit int) (*[]profile.ProfilePreview, error)
	GetSubscriptionsLike(name string, wall, offset, limit int) (*[]profile.ProfilePreview, error)
	Subscribe(subscriberId, profileId int) error
	Unsubscribe(subscriberId, profileId int) error
}
