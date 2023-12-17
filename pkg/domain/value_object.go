package domain

type ValueObject interface {
	SameValueObjectAs(other ValueObject) bool
}
