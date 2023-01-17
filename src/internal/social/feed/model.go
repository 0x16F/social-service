package feed

import (
	"core-server/src/internal/social/message"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

const (
	FeedTypeSubscriptions = "subscriptions"
	FeedTypeAll           = "all"
)

type Storager interface {
	Get(feedType string, userId, limit, offset int) (*[]message.FetchedMessage, error)
	GetCount(feedType string, userId int) (*int, error)
}
