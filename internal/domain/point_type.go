package domain

type PointTypes interface {
	int | bool
}

type PointType[T PointTypes] struct {
	Id    int
	Owner int
	Min   T
	Max   T
}
