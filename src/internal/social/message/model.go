package message

import "database/sql"

type Storage struct {
	db *sql.DB
}

type Message struct {
	Id       int64  `json:"id"`
	AuthorId int    `json:"author_id"`
	UserId   int    `json:"user_id"`
	Content  string `json:"content"`
}

type FetchedMessage struct {
	Id       int64  `json:"id"`
	AuthorId int    `json:"author_id"`
	UserId   int    `json:"user_id"`
	Content  string `json:"content"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Likes    int    `json:"likes"`
	IsLiked  bool   `json:"is_liked"`
}

type Storager interface {
	Get(id int64) (*Message, error)
	GetAll(wall, requester, offset, limit int) (*[]FetchedMessage, error)
	GetCount(wall int) (*int, error)
	GetWithInfo(id int64, requester int) (*FetchedMessage, error)
	Create(m *Message) error
	Delete(id int64) error
}
