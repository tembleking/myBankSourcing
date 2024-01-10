package domain

type Command interface {
	SameCommandAs(other Command) bool
}
