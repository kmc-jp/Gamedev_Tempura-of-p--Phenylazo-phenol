package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type StdTransform struct {
	position Utils.Position
}

// Position implements [gameobject.Transform].
func (s StdTransform) Position() Utils.Position {
	return s.position
}

func NewStdTransform(position Utils.Position) Utils.Factory[StdTransform] {
	return func() (*StdTransform, error) {
		return &StdTransform{position: position}, nil
	}
}

// test
var _ Utils.Factory[StdTransform] = NewStdTransform(Utils.Position{})

func test() {
	var _ gameobject.Transform = StdTransform{}
}
