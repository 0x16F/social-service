package reply

import "database/sql"

type Storage struct {
	db *sql.DB
}

type Reply struct {
	Id        int64  `json:"id"`
	ProfileId int    `json:"profile_id"`
	MessageId int64  `json:"message_id"`
	Content   string `json:"content"`
	Author    int    `json:"author_id"`
	Nickname  string `json:"nickname"`
}

type Storager interface {
	Get(id int64) (*Reply, error)
	GetAll(messageId int64) (*[]Reply, error)
	Create(m *Reply) error
	Delete(id int64) error
}
