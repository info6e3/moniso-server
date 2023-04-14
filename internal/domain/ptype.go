package domain

type PType uint8

const (
	Bool PType = iota
	Number
)

func (pt PType) Valid() bool {
	switch pt {
	case Bool:
		return true
	case Number:
		return true
	default:
		return false
	}
}
