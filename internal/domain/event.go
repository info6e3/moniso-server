package domain

import "time"

type Event struct {
	Id          int
	Title       string
	Date        time.Time
	Description string
}
