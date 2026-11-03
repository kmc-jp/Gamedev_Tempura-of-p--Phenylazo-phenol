package component

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type ShelfTile struct {
}

func (s ShelfTile) Interact() error {
	fmt.Printf("Interact with Shelf.\n")
	return nil
}

func (s ShelfTile) Update(active bool) error {
	// なにもしない
	return nil
}

func NewShelfTile() Utils.Factory[*ShelfTile] {
	return func() (*ShelfTile, error) {
		return &ShelfTile{}, nil
	}
}

// test
var _ InteractableTile = ShelfTile{}
