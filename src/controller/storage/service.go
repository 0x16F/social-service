package storage

import (
	"core-server/src/internal/invites"
	"core-server/src/internal/social"
	"core-server/src/internal/users"
	"core-server/src/pkg/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewStorage() Storager {
	return &Storage{}
}

func (s *Storage) Connect(cfg *config.Database) (*Storage, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Schema, cfg.SSL))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{
		db:      db,
		Users:   users.NewStorage(db),
		Invites: invites.NewStorage(db),
		Socials: social.NewStorage(db),
	}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}
