package model

import "time"

type Message struct {
	ID      int       `json:"id" db:"id"`
	From    int       `json:"from" db:"from"`
	To      int       `json:"to" db:"to"`
	Msg     string    `json:"message" db:"message"`
	Time    time.Time `json:"time" db:"time"`
	HasSent bool      `json:"has_sent" db:"has_sent"`
}
