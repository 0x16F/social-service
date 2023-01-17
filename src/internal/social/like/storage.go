package like

import "database/sql"

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(m *Like) error {
	if _, err := s.db.Exec("INSERT INTO social_likes VALUES ($1, $2)", &m.Id, &m.From); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(m *Like) error {
	if _, err := s.db.Exec("DELETE FROM social_likes WHERE id = $1 AND \"from\" = $2", m.Id, m.From); err != nil {
		return err
	}

	return nil
}
