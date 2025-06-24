package db

import (
    "database/sql"
	"time"

    "messenger/internal/model"
)

type MessageRepository struct {
    DB *sql.DB
}

func (r *MessageRepository) GetUnsentMessages(limit int) ([]model.Message, error) {
    rows, err := r.DB.Query("SELECT id, phone_number, content FROM messages WHERE sent = false LIMIT $1", limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []model.Message
    for rows.Next() {
        var msg model.Message
        if err := rows.Scan(&msg.ID, &msg.PhoneNumber, &msg.Content); err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }
    return messages, nil
}

func (r *MessageRepository) MarkMessageAsSent(id int) error {
    _, err := r.DB.Exec("UPDATE messages SET sent = true, sent_at = $1 WHERE id = $2", time.Now(), id)
    return err
}

func (r *MessageRepository) GetSentMessages() ([]model.Message, error) {
    rows, err := r.DB.Query("SELECT id, phone_number, content, sent, sent_at FROM messages WHERE sent = true ORDER BY sent_at DESC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []model.Message
    for rows.Next() {
        var msg model.Message
        if err := rows.Scan(&msg.ID, &msg.PhoneNumber, &msg.Content, &msg.Sent, &msg.SentAt); err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }
    return messages, nil
}
