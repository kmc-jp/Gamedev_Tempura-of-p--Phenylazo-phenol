package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

// タイル上を瞬間移動しながら動くもの
type TileTeleporter struct {
	Data     tilemap.TileData
	Position tilemap.TilePosition
	tf       gameobject.Transform
	tm       *TileMap
	anim     *AnimRender
}

// test
var _ gameobject.Component = &TileTeleporter{}

func (w *TileTeleporter) Update(active bool) (err error) {
	if !active {
		return nil
	}

	w.tf.SetPosition(w.tm.TilePositionToWorldPosition(w.Position, w.Data))

	return nil
}

func (w *TileTeleporter) Move(Direction Utils.Vec2Int) {
	moveto := w.Position.Position.Add(Direction)
	if enterable, ok := w.tm.Enterable[moveto]; enterable || !ok {
		w.Position.Move(Direction)
	}
}

func NewTileTeleporter(tiledata tilemap.TileData, pos tilemap.TilePosition, tf gameobject.Transform, Map *TileMap, anim *AnimRender) Utils.Factory[*TileTeleporter] {
	return func() (*TileTeleporter, error) {
		t := TileTeleporter{
			Data:     tiledata,
			Position: pos,
			tf:       tf,
			tm:       Map,
			anim:     anim,
		}
		return &t, nil
	}
}
