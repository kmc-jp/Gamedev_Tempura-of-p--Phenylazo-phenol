package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type Tile struct {
	Data tilemap.TileData
}

// test
var _ gameobject.Component = Tile{}

func (h Tile) Update(active bool) (err error) {
	// なにもしない
	return nil
}

func NewTile(tile tilemap.TileData) Utils.Factory[Tile] {
	return func() (*Tile, error) {
		t := Tile{}
		return &t, nil
	}
}

// test
var _ Utils.Factory[HogeComponent] = NewHogeComponent(nil)
