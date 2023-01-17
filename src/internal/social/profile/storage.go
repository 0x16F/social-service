package profile

import "database/sql"

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) GetAllLike(name string, offset, limit int) (*[]ProfilePreview, error) {
	profiles := make([]ProfilePreview, 0)

	query := `
		SELECT
			id, 
			nickname, 
			role
		FROM users
		WHERE LOWER(nickname) LIKE LOWER($1)
		GROUP BY id
		ORDER BY id ASC
		OFFSET $2
		LIMIT $3
	`

	rows, err := s.db.Query(query, name+"%", offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		profile := ProfilePreview{}

		if err := rows.Scan(&profile.UserId, &profile.Nickname, &profile.Role); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return &profiles, nil
}

func (s *Storage) Get(from, profileId int) (*ProfileData, error) {
	profile := ProfileData{}

	query := `
		SELECT
			social_profiles.*,
			users.nickname,
			users.register_date,
			users.role,
			(SELECT COUNT(*) FROM social_subscribers WHERE user_id = $2) AS subscribers,
			(SELECT COUNT(*) FROM social_subscribers WHERE subscriber_id = $2) AS subscriptions,
			(SELECT COUNT(*) FROM social_subscribers WHERE user_id = $2 AND subscriber_id = $1 LIMIT 1) AS is_subscriptions
		FROM social_profiles
		LEFT JOIN users ON social_profiles.user_id = users.id
		WHERE "id" = $2
	`

	if err := s.db.QueryRow(query, from, profileId).Scan(&profile.UserId, &profile.BIO, &profile.Telegram, &profile.Discord, &profile.IsWallClosed, &profile.IsCommentsClosed, &profile.Nickname, &profile.RegisterDate, &profile.Role, &profile.Subscribers, &profile.Subscriptions, &profile.IsSubscribed); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *Storage) GetCount() (*int, error) {
	total := 0

	if err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total); err != nil {
		return nil, nil
	}

	return &total, nil
}

func (s *Storage) GetAll(offset, limit int) (*[]ProfilePreview, error) {
	profiles := make([]ProfilePreview, 0)

	query := `
		SELECT 
			id, 
			nickname, 
			role
		FROM users
		GROUP BY id
		ORDER BY id ASC
		OFFSET $1
		LIMIT $2
	`

	rows, err := s.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		profile := ProfilePreview{}

		if err := rows.Scan(&profile.UserId, &profile.Nickname, &profile.Role); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return &profiles, nil
}

func (s *Storage) Update(m *Profile) error {
	if _, err := s.db.Exec("UPDATE social_profiles SET bio = $4, telegram = $5, discord = $6, is_wall_closed = $7 WHERE user = $1", &m.UserId, &m.BIO, &m.Telegram, &m.Discord, &m.IsWallClosed); err != nil {
		return err
	}

	return nil
}
