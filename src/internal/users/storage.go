package users

import "database/sql"

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Fetch(id int) (*User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id)

	u := User{}

	if err := row.Scan(&u.Id, &u.Login, &u.Nickname, &u.Role, &u.Password, &u.Salt, &u.Balance, &u.RegisterDate); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Storage) FetchLogin(login string) (*User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE login = lower($1)", &login)

	u := User{}

	if err := row.Scan(&u.Id, &u.Login, &u.Nickname, &u.Role, &u.Password, &u.Salt, &u.Balance, &u.RegisterDate); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Storage) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (s *Storage) Update(u *User) error {
	_, err := s.db.Exec("UPDATE users SET password = $1, role = $2, balance = $3 WHERE login = lower($4)", &u.Password, &u.Role, &u.Balance, &u.Login)
	return err
}

func (s *Storage) Create(u *User) (*int, error) {
	id := 0

	if err := s.db.QueryRow("INSERT INTO users (login, nickname, password, salt, register_date) VALUES (lower($1), $2, $3, $4, $5) RETURNING id", &u.Login, &u.Nickname, &u.Password, &u.Salt, &u.RegisterDate).Scan(&id); err != nil {
		return nil, err
	}

	if _, err := s.db.Exec("INSERT INTO social_profiles (\"user_id\") VALUES ((SELECT id FROM users WHERE login = lower($1)))", &u.Login); err != nil {
		return nil, err
	}

	return &id, nil
}
