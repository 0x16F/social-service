package like

import "database/sql"

type Storage struct {
	db *sql.DB
}

type Like struct {
	Id   int64 `json:"id"`
	From int   `json:"from"`
}

type Storager interface {
	Create(m *Like) error
	Delete(m *Like) error
}
