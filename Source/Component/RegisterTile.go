package component

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type RegisterTile struct {
	scene *scene.PlayScene
}

func (s RegisterTile) Interact(player *Player) error {
	fmt.Printf("Interact with Regester.\n")
	if s.scene.Gamedata.State == game.Play_Explore {
		s.scene.Gamedata.State = game.Play_Select
	} else {
		return fmt.Errorf("Interactした時点のstateが不正な値です\n")
	}
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
