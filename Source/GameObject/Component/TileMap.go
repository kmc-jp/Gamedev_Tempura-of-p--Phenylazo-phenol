package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TileMap struct {
	Tiles    []*Tile
	TileSize Utils.Vec2
}

// test
var _ gameobject.Component = TileMap{}

func (h TileMap) Update(active bool) (err error) {
	// なにもしない
	return nil
}

func NewTileMap(mapdata tilemap.TileMapData) Utils.Factory[TileMap] {
	t := TileMap{}
	return func() (*TileMap, error) {
		return &t, nil
	}
}

// test
var _ Utils.Factory[HogeComponent] = NewHogeComponent(nil)
