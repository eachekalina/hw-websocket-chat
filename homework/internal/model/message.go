package model

import "time"

type Message struct {
	SentTime time.Time `db:"sent_time"`
	Nickname string    `db:"nickname"`
	Message  string    `db:"message"`
}
