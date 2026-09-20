package Utils

type Factory[T any] func() (T, error)
