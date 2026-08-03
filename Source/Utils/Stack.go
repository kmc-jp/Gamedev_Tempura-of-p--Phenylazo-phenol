package Utils

import "iter"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	} else {
		v := s.items[len(s.items)-1]
		return v, true
	}
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Len: 要素数を返す
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// All: range 句で回すためのイテレータ関数 (LIFO 順)
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		// スタックの末尾（一番上）から逆順にたどる
		for i := len(s.items) - 1; i >= 0; i-- {
			// yield が false を返したら range ループが break されたことを意味する
			if !yield(s.items[i]) {
				return
			}
		}
	}
}

// Rev: range 句で古い順（スタックの下→上）に回すためのイテレータ関数
func (s *Stack[T]) Rev() iter.Seq[T] {
	return func(yield func(T) bool) {
		// スタックの先頭（一番下）からたどる
		for i := 0; i < len(s.items); i++ {
			// yield が false を返したら range ループが break されたことを意味する
			if !yield(s.items[i]) {
				return
			}
		}
	}
}

func (s *Stack[T]) Clear() {
	s.items = nil
}
