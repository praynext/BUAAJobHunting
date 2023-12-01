package model

import "time"

type Reminder struct {
	ID        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	Message   string    `json:"message" db:"message"`
	Time      time.Time `json:"time" db:"time"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	HasSent   bool      `json:"has_sent" db:"has_sent"`
}
