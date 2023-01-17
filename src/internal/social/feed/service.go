package feed

import (
	"core-server/src/controller/web/errors"
	"core-server/src/internal/social/message"
	"database/sql"
	goErrors "errors"
)

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Get(feedType string, userId, limit, offset int) (*[]message.FetchedMessage, error) {
	query := ``

	switch feedType {
	case FeedTypeSubscriptions:
		query = `
			SELECT
				social_messages.*,
				(SELECT users.nickname FROM users WHERE social_messages.author_id = users.id LIMIT 1),
				(SELECT users.role FROM users WHERE social_messages.author_id = users.id LIMIT 1),
				(SELECT COUNT(id) FROM social_likes WHERE social_likes.id = social_messages.id) as likes,
				COUNT((SELECT id FROM social_likes WHERE social_likes.id = social_messages.id AND social_likes.from = $1 LIMIT 1)) as is_liked			
			FROM social_subscribers
			JOIN social_messages ON social_subscribers.user_id = social_messages.author_id AND social_messages.author_id = social_messages.user_id
			WHERE subscriber_id = $1
			GROUP BY social_messages.id
			ORDER BY id DESC
			LIMIT $2
			OFFSET $3
		`
	case FeedTypeAll:
		query = `
			SELECT
				*,
				(SELECT users.nickname FROM users WHERE social_messages.author_id = users.id LIMIT 1),
				(SELECT users.role FROM users WHERE social_messages.author_id = users.id LIMIT 1),
				(SELECT COUNT(id) FROM social_likes WHERE social_likes.id = social_messages.id) as likes,
				COUNT((SELECT id FROM social_likes WHERE social_likes.id = social_messages.id AND social_likes.from = $1 LIMIT 1)) as is_liked
			FROM social_messages
			WHERE author_id != $1
			GROUP BY id
			ORDER BY id DESC
			LIMIT $2
			OFFSET $3
		`
	}

	if query == "" {
		return nil, goErrors.New(errors.ErrIncorrectFeedType)
	}

	messages := make([]message.FetchedMessage, 0)

	rows, err := s.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		message := message.FetchedMessage{}

		if err := rows.Scan(&message.Id, &message.AuthorId, &message.UserId, &message.Content, &message.Nickname, &message.Role, &message.Likes, &message.IsLiked); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return &messages, nil
}

func (s *Storage) GetCount(feedType string, userId int) (*int, error) {
	query := ""

	switch feedType {
	case FeedTypeSubscriptions:
		query = `
			SELECT
				COUNT(social_messages.*)
			FROM social_subscribers
			JOIN social_messages ON social_subscribers.user_id = social_messages.author_id AND social_messages.author_id = social_messages.user_id
			WHERE subscriber_id = $1
		`
	case FeedTypeAll:
		query = `
			SELECT
				COUNT(*)
			FROM social_messages
			WHERE author_id != $1
		`
	}

	count := 0

	if err := s.db.QueryRow(query, userId).Scan(&count); err != nil {
		return nil, err
	}

	return &count, nil
}
