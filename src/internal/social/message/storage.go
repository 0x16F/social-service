package message

import "database/sql"

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) GetCount(wall int) (*int, error) {
	total := 0

	if err := s.db.QueryRow("SELECT COUNT(*) FROM social_messages WHERE user_id = $1", wall).Scan(&total); err != nil {
		return nil, nil
	}

	return &total, nil
}

func (s *Storage) GetAll(wall, likeFrom, offset, limit int) (*[]FetchedMessage, error) {
	messages := make([]FetchedMessage, 0)

	query := `
		SELECT
			*,
			(SELECT users.nickname FROM users WHERE social_messages.author_id = users.id LIMIT 1),
			(SELECT users.role FROM users WHERE social_messages.author_id = users.id LIMIT 1),
			(SELECT COUNT(id) FROM social_likes WHERE social_likes.id = social_messages.id) as likes,
			COUNT((SELECT id FROM social_likes WHERE social_likes.id = social_messages.id AND social_likes.from = $2 LIMIT 1)) as is_liked
		FROM social_messages
		WHERE user_id = $1
		GROUP BY id
		ORDER BY id DESC
		OFFSET $3
		LIMIT $4
	`

	rows, err := s.db.Query(query, wall, likeFrom, offset, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		message := FetchedMessage{}

		if err := rows.Scan(&message.Id, &message.AuthorId, &message.UserId, &message.Content, &message.Nickname, &message.Role, &message.Likes, &message.IsLiked); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return &messages, nil
}

func (s *Storage) Get(id int64) (*Message, error) {
	message := Message{}

	if err := s.db.QueryRow("SELECT * FROM social_messages WHERE id = $1", id).Scan(&message.Id, &message.AuthorId, &message.UserId, &message.Content); err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Storage) GetWithInfo(id int64, likeFrom int) (*FetchedMessage, error) {
	message := FetchedMessage{}

	query := `
		SELECT
			*,
			(SELECT users.nickname FROM users WHERE social_messages.author_id = users.id),
			(SELECT users.role FROM users WHERE social_messages.author_id = users.id LIMIT 1),
			(SELECT COUNT(id) FROM social_likes WHERE social_likes.id = social_messages.id) as likes,
			COUNT((SELECT id FROM social_likes WHERE social_likes.id = social_messages.id AND social_likes.from = $2 LIMIT 1)) as is_liked
		FROM social_messages
		WHERE id = $1
		GROUP BY id
		ORDER BY id DESC
	`

	if err := s.db.QueryRow(query, id, likeFrom).Scan(&message.Id, &message.AuthorId, &message.UserId, &message.Content, &message.Nickname, &message.Role, &message.Likes, &message.IsLiked); err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Storage) Create(m *Message) error {
	if _, err := s.db.Exec("INSERT INTO social_messages VALUES ($1, $2, $3, $4)", &m.Id, &m.AuthorId, &m.UserId, &m.Content); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(id int64) error {
	if _, err := s.db.Exec("DELETE FROM social_messages WHERE id = $1", &id); err != nil {
		return err
	}

	return nil
}
