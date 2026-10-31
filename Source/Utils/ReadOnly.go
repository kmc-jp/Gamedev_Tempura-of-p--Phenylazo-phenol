package Utils

type ReadOnly[T any] interface {
	Value() T
}

type readonly[T any] struct {
	value T
}

func (i readonly[T]) Value() T {
	return i.value
}

func NewReadOnly[T any](value T) ReadOnly[T] {
	return readonly[T]{value: value}
}
