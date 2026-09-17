package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type StdTransform struct {
}

// Position implements [gameobject.Transform].
func (s StdTransform) Position() Utils.Position {
	panic("unimplemented")
}

func NewTransform() (*StdTransform, error) {
	return &StdTransform{}, nil
}

// test
var _ Utils.Factory[StdTransform] = NewTransform

func test() {
	var _ gameobject.Transform = StdTransform{}
}
