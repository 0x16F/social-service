package subs

import (
	"core-server/src/internal/social/profile"
	"database/sql"
)

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Subscribe(subscriberId, profileId int) error {
	_, err := s.db.Exec("INSERT INTO social_subscribers VALUES($1, $2)", subscriberId, profileId)
	return err
}

func (s *Storage) Unsubscribe(subscriberId, profileId int) error {
	_, err := s.db.Exec("DELETE FROM social_subscribers WHERE subscriber_id = $1 and user_id = $2", subscriberId, profileId)
	return err
}

func (s *Storage) GetSubscribers(wall, offset, limit int) (*[]profile.ProfilePreview, error) {
	profiles := make([]profile.ProfilePreview, 0)

	query := `
		SELECT 
			subscriber_id, nickname, role
		FROM social_subscribers
		LEFT JOIN users 
			ON social_subscribers.subscriber_id = users.id
		WHERE user_id = $1
		ORDER BY id ASC
		OFFSET $2
		LIMIT $3
	`

	rows, err := s.db.Query(query, wall, offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		profile := profile.ProfilePreview{}

		if err := rows.Scan(&profile.UserId, &profile.Nickname, &profile.Role); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return &profiles, nil
}

func (s *Storage) GetSubscriptions(wall, offset, limit int) (*[]profile.ProfilePreview, error) {
	profiles := make([]profile.ProfilePreview, 0)

	query := `
		SELECT 
			user_id, nickname, role
		FROM social_subscribers
		LEFT JOIN users 
			ON social_subscribers.user_id = users.id
		WHERE subscriber_id = $1
		ORDER BY id ASC
		OFFSET $2
		LIMIT $3
	`

	rows, err := s.db.Query(query, wall, offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		profile := profile.ProfilePreview{}

		if err := rows.Scan(&profile.UserId, &profile.Nickname, &profile.Role); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return &profiles, nil
}

func (s *Storage) GetSubscribersLike(name string, wall, offset, limit int) (*[]profile.ProfilePreview, error) {
	profiles := make([]profile.ProfilePreview, 0)

	query := `
		SELECT 
			subscriber_id, nickname, role
		FROM social_subscribers
		LEFT JOIN users 
			ON social_subscribers.subscriber_id = users.id
		WHERE user_id = $1 AND LOWER(nickname) LIKE LOWER($2)
		ORDER BY id ASC
		OFFSET $3
		LIMIT $4
	`

	rows, err := s.db.Query(query, wall, name+"%", offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		profile := profile.ProfilePreview{}

		if err := rows.Scan(&profile.UserId, &profile.Nickname, &profile.Role); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return &profiles, nil
}

func (s *Storage) GetSubscriptionsLike(name string, wall, offset, limit int) (*[]profile.ProfilePreview, error) {
	profiles := make([]profile.ProfilePreview, 0)

	query := `
		SELECT 
			user_id, nickname, role
		FROM social_subscribers
		LEFT JOIN users 
			ON social_subscribers.user_id = users.id
		WHERE subscriber_id = $1 AND LOWER(nickname) LIKE LOWER($2)
		ORDER BY id ASC
		OFFSET $3
		LIMIT $4
	`

	rows, err := s.db.Query(query, wall, name+"%", offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		profile := profile.ProfilePreview{}

		if err := rows.Scan(&profile.UserId, &profile.Nickname, &profile.Role); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return &profiles, nil
}
