package lazy

import (
	"sync"
)

type Lazy[T any] struct {
	value T
	once  sync.Once
}

func (l *Lazy[T]) GetOrInit(init func() T) T {
	l.once.Do(func() {
		l.value = init()
	})
	return l.value
}
