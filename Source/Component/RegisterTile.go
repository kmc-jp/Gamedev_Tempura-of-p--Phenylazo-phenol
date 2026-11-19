package component

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type RegisterTile struct {
}

func (s RegisterTile) Interact() error {
	fmt.Printf("Interact with Regester.\n")
	return nil
}

func (s RegisterTile) Update(active bool) error {
	// なにもしない
	return nil
}

func NewRegisterTile() Utils.Factory[*RegisterTile] {
	return func() (*RegisterTile, error) {
		return &RegisterTile{}, nil
	}
}

// test
var _ InteractableTile = RegisterTile{}
