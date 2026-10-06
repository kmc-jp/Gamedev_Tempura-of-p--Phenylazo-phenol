package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TileWalker struct {
	Data     tilemap.TileData
	Position tilemap.TilePosition
}

// test
var _ gameobject.Component = Tile{}

func (h TileWalker) Update(active bool) (err error) {
	// なにもしない
	return nil
}

func (h *TileWalker) Move()

func NewTileWalker(tiledata tilemap.TileData, pos tilemap.TilePosition) Utils.Factory[*TileWalker] {
	return func() (*TileWalker, error) {
		t := TileWalker{
			Data:     tiledata,
			Position: pos,
		}
		return &t, nil
	}
}
