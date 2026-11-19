package component

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type RegisterTile struct {
	scene *scene.PlayScene
}

func (s RegisterTile) Interact() error {
	fmt.Printf("Interact with Regester.\n")
	return nil
}

func (s RegisterTile) Update(active bool) error {
	// なにもしない
	return nil
}

func NewRegisterTile(scene *scene.PlayScene) Utils.Factory[*RegisterTile] {
	return func() (*RegisterTile, error) {
		return &RegisterTile{
			scene: scene,
		}, nil
	}
}

// test
var _ InteractableTile = RegisterTile{}
