package component

import "fmt"

type ShelfTile struct {
}

func (s ShelfTile) Interact() error {
	fmt.Printf("Interact with Shelf.")
	return nil
}

func (s ShelfTile) Update(active bool) error {
	// なにもしない
	return nil
}

// test
var _ InteractableTile = ShelfTile{}
