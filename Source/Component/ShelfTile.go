package component

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type ShelfTile struct {
	scene *scene.PlayScene
}

func (s ShelfTile) Interact(player *Player) error {
	fmt.Printf("Interact with Shelf.\n")
	return nil
}

func (s ShelfTile) Update(active bool) error {
	// なにもしない
	return nil
}

func NewShelfTile(scene *scene.PlayScene) Utils.Factory[*ShelfTile] {
	return func() (*ShelfTile, error) {
		return &ShelfTile{
			scene: scene,
		}, nil
	}
}

// test
var _ InteractableTile = ShelfTile{}
