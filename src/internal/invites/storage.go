package invites

import "database/sql"

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Fetch(invite string) (*FetchedInvite, error) {
	query := `
		SELECT
			created_invites.invite, 
			created_invites.author_id,
			created_invites.create_time, 
			coalesce(activated_invites.activated_by, 0),
			coalesce(activated_invites.activate_time, 0)
		FROM created_invites
		LEFT JOIN activated_invites 
			ON activated_invites.invite = created_invites.invite
		WHERE created_invites.invite = $1
		LIMIT 1
	`

	i := FetchedInvite{}

	if err := s.db.QueryRow(query, invite).Scan(&i.Invite, &i.AuthorId, &i.CreateTime, &i.ActivatedBy, &i.ActivateTime); err != nil {
		return nil, err
	}

	return &i, nil
}

func (s *Storage) Delete(invite string) error {
	_, err := s.db.Exec("DELETE FROM created_invites WHERE invite = $1", invite)
	return err
}

func (s *Storage) Activate(i *ActivatedInvite) error {
	_, err := s.db.Exec("INSERT INTO activated_invites VALUES ($1, $2, $3)", &i.Invite, &i.ActivatedBy, &i.ActivateTime)
	return err
}

func (s *Storage) Create(i *Invite) error {
	_, err := s.db.Exec("INSERT INTO created_invites VALUES ($1, $2, $3)", &i.Invite, &i.AuthorId, &i.CreateTime)
	return err
}
