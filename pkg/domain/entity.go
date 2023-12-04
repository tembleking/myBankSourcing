package domain

type Entity interface {
	ID() string
	SameEntityAs(other Entity) bool
}
