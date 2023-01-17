package invites

import "database/sql"

type Invite struct {
	Invite     string `json:"invite"`
	AuthorId   int    `json:"author_id"`
	CreateTime int64  `json:"create_time"`
}

type ActivatedInvite struct {
	Invite       string `json:"invite"`
	ActivatedBy  int    `json:"activated_by"`
	ActivateTime int64  `json:"activate_time"`
}

type FetchedInvite struct {
	Invite       string `json:"invite"`
	AuthorId     int    `json:"author_id"`
	CreateTime   int64  `json:"create_time"`
	ActivatedBy  int    `json:"activated_by"`
	ActivateTime int64  `json:"activate_time"`
}

type Storage struct {
	db *sql.DB
}

type Storager interface {
	Fetch(i string) (*FetchedInvite, error)
	Delete(i string) error
	Activate(i *ActivatedInvite) error
	Create(i *Invite) error
}
