package model

import (
	"database/sql"
)

type Message struct {
	ID       int64  `json:"id"`
	Body     string `json:"body"`
	Username string `json:"username"`
}

func MessagesAll(db *sql.DB) ([]*Message, error) {

	rows, err := db.Query(`select id, body, username from message`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []*Message
	for rows.Next() {
		m := &Message{}
		if err := rows.Scan(&m.ID, &m.Body, &m.Username); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

func MessageByID(db *sql.DB, id string) (*Message, error) {
	m := &Message{}

	if err := db.QueryRow(`select id, body, username from message where id = ?`, id).Scan(&m.ID, &m.Body, &m.Username); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Message) Insert(db *sql.DB) (*Message, error) {
	res, err := db.Exec(`insert into message (body, username) values (?, ?)`, m.Body, m.Username)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:   id,
		Body: m.Body,
		Username: m.Username,
	}, nil
}

func (m *Message) Update(db *sql.DB) (*Message, error) {
	_, err := db.Exec(`UPDATE message SET body = ? WHERE id = ?`, m.Body, m.ID)
	if err != nil {
		return nil, err
	}

	return m, nil
}
