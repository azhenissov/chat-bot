package domain

import (
	"time"
)

type Chat struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        int64     `json:"id"`
	ChatID    string    `json:"chat_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
