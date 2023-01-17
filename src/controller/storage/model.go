package storage

import (
	"core-server/src/internal/invites"
	"core-server/src/internal/social"
	"core-server/src/internal/users"
	"core-server/src/pkg/config"
	"database/sql"
)

type Storage struct {
	db      *sql.DB
	Users   users.Storager
	Invites invites.Storager
	Socials social.Storage
}

type Storager interface {
	Connect(cfg *config.Database) (*Storage, error)
	Close()
}
