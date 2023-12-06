package model

import "time"

type LastChat struct {
	From int       `json:"from" db:"from"`
	To   int       `json:"to" db:"to"`
	Time time.Time `json:"time" db:"time"`
}
