package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type Player struct {
}

func (p Player) Update(active bool) (err error) {
	// なにもしない
	return nil
}

func NewPlayer() Utils.Factory[*Player] {
	return func() (*Player, error) {
		return &Player{}, nil
	}
}

// test
var _ gameobject.Component = Player{}
