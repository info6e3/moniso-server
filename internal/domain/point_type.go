package domain

type PointType struct {
	Id    int
	Type  PType
	Title string
	Owner int
	Min   int
	Max   int
}
