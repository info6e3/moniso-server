package domain

import "time"

type Event struct {
	Id          int
	Title       string
	Owner       int
	Description string
	Date        time.Time
}
