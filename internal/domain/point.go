package domain

import "time"

type Point struct {
	Id          int
	Type        int
	Value       int
	Description string
	Date        time.Time
}
