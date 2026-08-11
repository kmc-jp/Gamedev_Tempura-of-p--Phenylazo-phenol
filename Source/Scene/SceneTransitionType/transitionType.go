package transition

import (
	"fmt"
)

const (
	none = iota
	newScene
	push
	pop
	sentinel
)

type Type struct {
	val int
}

func newTransitionType(v int) (*Type, error) {
	if v <= 0 || sentinel <= v {
		return nil, fmt.Errorf("newTransitionType に不適切な引数 %d が渡されました", v)
	}
	return &Type{val: v}, nil
}

func NewScene() Type {
	tt, _ := newTransitionType(newScene)
	return *tt
}

func Push() Type {
	tt, _ := newTransitionType(push)
	return *tt
}

func Pop() Type {
	tt, _ := newTransitionType(pop)
	return *tt
}

func None() Type {
	tt, _ := newTransitionType(none)
	return *tt
}
