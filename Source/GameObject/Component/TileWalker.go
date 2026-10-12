package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TileWalker struct {
	Data     tilemap.TileData
	Position tilemap.TilePosition
	tf       gameobject.Transform
	tm       *TileMap
}

// test
var _ gameobject.Component = Tile{}

func (h TileWalker) Update(active bool) (err error) {
	// なにもしない
	return nil
}

func (h *TileWalker) Move(Direction Utils.Vec2Int) {
	h.Position.Move(Direction)
}

func NewTileWalker(tiledata tilemap.TileData, pos tilemap.TilePosition, tf gameobject.Transform, Map *TileMap) Utils.Factory[*TileWalker] {
	return func() (*TileWalker, error) {
		t := TileWalker{
			Data:     tiledata,
			Position: pos,
		}
		return &t, nil
	}
}
