package domain

import (
	"time"
)

type Point[T PointTypes] struct {
	Id          int
	Title       string
	Owner       int
	Type        PointType[T]
	Value       T
	Date        time.Time
	Description string
}
