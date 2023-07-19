package typex

import "sync"

type Lazy[T any] interface {
	Value() T
}

func ByLazy[T any](ctor func() T) Lazy[T] {
	return &lazy[T]{
		ctor: ctor,
	}
}

type lazy[T any] struct {
	o     sync.Once
	value T
	ctor  func() T
}

func (l *lazy[T]) Value() T {
	l.o.Do(func() {
		l.value = l.ctor()
	})

	return l.value
}
