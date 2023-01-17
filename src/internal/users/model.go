package users

import (
	"database/sql"
)

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	Nickname     string `json:"nickname"`
	Role         string `json:"role"`
	Password     string `json:"password"`
	Salt         string `json:"salt"`
	Balance      int64  `json:"balance"`
	RegisterDate int64  `json:"register_date"`
}

type Storage struct {
	db *sql.DB
}

type Storager interface {
	Fetch(id int) (*User, error)
	FetchLogin(login string) (*User, error)
	Delete(id int) error
	Update(u *User) error
	Create(u *User) (*int, error)
}
