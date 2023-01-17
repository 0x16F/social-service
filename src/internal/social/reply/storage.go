package reply

import "database/sql"

func NewStorage(db *sql.DB) Storager {
	return &Storage{
		db: db,
	}
}

func (s *Storage) GetAll(messageId int64) (*[]Reply, error) {
	replies := make([]Reply, 0)

	rows, err := s.db.Query("SELECT *, (SELECT users.nickname FROM users WHERE users.id = social_replies.author_id )FROM social_replies WHERE message_id = $1", messageId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		reply := Reply{}

		if err := rows.Scan(&reply.Id, &reply.ProfileId, &reply.MessageId, &reply.Content, &reply.Author, &reply.Nickname); err != nil {
			return nil, err
		}

		replies = append(replies, reply)
	}

	return &replies, nil
}

func (s *Storage) Get(id int64) (*Reply, error) {
	reply := Reply{}

	if err := s.db.QueryRow("SELECT * FROM social_replies WHERE id = $1", id).Scan(&reply.Id, &reply.ProfileId, &reply.MessageId, &reply.Content, &reply.Author); err != nil {
		return nil, err
	}

	return &reply, nil
}

func (s *Storage) Create(m *Reply) error {
	if _, err := s.db.Exec("INSERT INTO social_replies VALUES ($1, $2, $3, $4, $5)", &m.Id, &m.ProfileId, &m.MessageId, &m.Content, &m.Author); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(id int64) error {
	if _, err := s.db.Exec("DELETE FROM social_replies WHERE id = $1", &id); err != nil {
		return err
	}

	return nil
}
