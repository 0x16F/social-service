package profile

import "database/sql"

type Storage struct {
	db *sql.DB
}

type Profile struct {
	UserId           int    `json:"id"`
	BIO              string `json:"bio"`
	Telegram         string `json:"telegram"`
	Discord          string `json:"discord"`
	IsWallClosed     bool   `json:"is_wall_closed"`
	IsCommentsClosed bool   `json:"is_comments_closed"`
}

type ProfileData struct {
	UserId           int    `json:"id"`
	Nickname         string `json:"nickname"`
	BIO              string `json:"bio"`
	Telegram         string `json:"telegram"`
	Discord          string `json:"discord"`
	IsWallClosed     bool   `json:"is_wall_closed"`
	IsCommentsClosed bool   `json:"is_comments_closed"`
	RegisterDate     int64  `json:"register_date"`
	Role             string `json:"role"`
	Subscribers      int64  `json:"subscribers"`
	Subscriptions    int64  `json:"subscriptions"`
	IsSubscribed     bool   `json:"is_subscribed"`
}

type ProfilePreview struct {
	UserId   int    `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type Storager interface {
	Get(from, profile int) (*ProfileData, error)
	GetAll(offset, limit int) (*[]ProfilePreview, error)
	GetAllLike(name string, offset, limit int) (*[]ProfilePreview, error)
	GetCount() (*int, error)
	Update(m *Profile) error
}
