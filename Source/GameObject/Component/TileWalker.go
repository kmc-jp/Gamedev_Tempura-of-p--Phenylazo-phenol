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
var _ gameobject.Component = &TileWalker{}

func (w *TileWalker) Update(active bool) (err error) {
	if !active {
		return nil
	}

	w.tf.SetPosition(w.tm.TilePositionToWorldPosition(w.Position))

	return nil
}

func (w *TileWalker) Move(Direction Utils.Vec2Int) {
	moveto := w.Position.Position.Add(Direction)
	if enterable, ok := w.tm.Enterable[moveto]; enterable || !ok {
		w.Position.Move(Direction)
	}
}

func NewTileWalker(tiledata tilemap.TileData, pos tilemap.TilePosition, tf gameobject.Transform, Map *TileMap) Utils.Factory[*TileWalker] {
	return func() (*TileWalker, error) {
		t := TileWalker{
			Data:     tiledata,
			Position: pos,
			tf:       tf,
			tm:       Map,
		}
		return &t, nil
	}
}
