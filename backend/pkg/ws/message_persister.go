package ws

import (
	"database/sql"
	"strconv"
	"time"
)

// DBMessagePersister implements MessagePersister using SQLite.
type DBMessagePersister struct {
	DB *sql.DB
}

func NewDBMessagePersister(db *sql.DB) *DBMessagePersister {
	return &DBMessagePersister{DB: db}
}

func (p *DBMessagePersister) SaveMessage(senderID int64, msg *Message) error {
	switch msg.Type {
	case "private":
		_, err := p.DB.Exec(`
			INSERT INTO Messages (sender_id, receiver_id, content, created_at)
			VALUES (?, ?, ?, ?)
		`, senderID, msg.To, msg.Content, time.Now().UTC())
		return err
	case "group":
		groupID, err := strconv.Atoi(msg.GroupID)
		if err != nil {
			return err
		}
		_, err = p.DB.Exec(`
			INSERT INTO Messages (sender_id, group_id, content, created_at)
			VALUES (?, ?, ?, ?)
		`, senderID, groupID, msg.Content, time.Now().UTC())
		return err
	default:
		return nil // Ignore broadcast messages for now
	}
}

func (p *DBMessagePersister) FetchPrivateMessages(userA, userB int64, limit int) ([]Message, error) {
	rows, err := p.DB.Query(`
		SELECT sender_id, receiver_id, content, strftime('%s', created_at)
		FROM Messages
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at DESC
		LIMIT ?
	`, userA, userB, userB, userA, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var senderID, receiverID int64
		var tsStr string
		err := rows.Scan(&senderID, &receiverID, &m.Content, &tsStr)
		if err != nil {
			continue
		}
		m.Type = "private"
		m.To = receiverID
		ts, _ := strconv.ParseInt(tsStr, 10, 64)
		m.Timestamp = ts
		messages = append(messages, m)
	}
	return messages, nil
}

func (p *DBMessagePersister) FetchGroupMessages(groupID int64, limit int) ([]Message, error) {
	rows, err := p.DB.Query(`
		SELECT sender_id, content, strftime('%s', created_at)
		FROM Messages
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, groupID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var senderID int64
		var tsStr string
		err := rows.Scan(&senderID, &m.Content, &tsStr)
		if err != nil {
			continue
		}
		m.Type = "group"
		m.Timestamp = time.Now().Unix()
		m.GroupID = strconv.FormatInt(groupID, 10)
		ts, _ := strconv.ParseInt(tsStr, 10, 64)
		m.Timestamp = ts
		messages = append(messages, m)
	}
	return messages, nil
}

func (p *DBMessagePersister) FetchPrivateMessagesPaginated(userA, userB int64, limit, offset int) ([]Message, error) {
	rows, err := p.DB.Query(`
		SELECT sender_id, receiver_id, content, strftime('%s', created_at)
		FROM Messages
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userA, userB, userB, userA, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var senderID, receiverID int64
		var tsStr string
		err := rows.Scan(&senderID, &receiverID, &m.Content, &tsStr)
		if err != nil {
			continue
		}
		m.Type = "private"
		m.To = receiverID
		ts, _ := strconv.ParseInt(tsStr, 10, 64)
		m.Timestamp = ts
		messages = append(messages, m)
	}
	return messages, nil
}

func (p *DBMessagePersister) FetchGroupMessagesPaginated(groupID int64, limit, offset int) ([]Message, error) {
	rows, err := p.DB.Query(`
		SELECT sender_id, content, strftime('%s', created_at)
		FROM Messages
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var senderID int64
		var tsStr string
		err := rows.Scan(&senderID, &m.Content, &tsStr)
		if err != nil {
			continue
		}
		m.Type = "group"
		m.GroupID = strconv.FormatInt(groupID, 10)
		ts, _ := strconv.ParseInt(tsStr, 10, 64)
		m.Timestamp = ts
		messages = append(messages, m)
	}
	return messages, nil
}
