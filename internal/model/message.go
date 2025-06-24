package model

import "time"

type Message struct {
    ID          int       `json:"id"`
    PhoneNumber string    `json:"phone_number"`
    Content     string    `json:"content"`
    Sent        bool      `json:"sent"`
    SentAt      time.Time `json:"sent_at"`
}