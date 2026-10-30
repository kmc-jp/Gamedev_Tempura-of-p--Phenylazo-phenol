package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

// タイル上を歩くコンポーネント
type TileWalker struct {
	Data     tilemap.TileData
	Position tilemap.TilePosition
	tf       gameobject.Transform
	tm       *TileMap
	anim     *AnimRender
	walktime Utils.ReadOnly[int] //フレーム数

	isWalking bool
	elapsed   int
	direction Utils.Vec2Int
}

// test
var _ gameobject.Component = &TileWalker{}

func (w *TileWalker) Update(active bool) (err error) {
	if !active {
		return nil
	}

	if w.isWalking {
		w.elapsed++
		if w.elapsed > w.walktime.Value() {
			w.isWalking = false
			w.elapsed = 0
			w.Position.Offset = Utils.NewZeroVec2()
			w.Position.Position = w.Position.Position.Add(w.direction)
		} else {
			diff := w.tm.TileSize.ToVec2().HadamardProd(w.direction.ToVec2()).ScalarProd(1 / float64(w.walktime.Value()))
			w.Position.Offset = w.Position.Offset.Add(diff)
		}
	}

	w.tf.SetPosition(w.tm.TilePositionToWorldPosition(w.Position, w.Data))
	return nil
}

func (w *TileWalker) MoveStart(Direction Utils.Vec2Int) {
	w.isWalking = true
	w.direction = Direction
}

func NewTileWalker(tiledata tilemap.TileData, pos tilemap.TilePosition, tf gameobject.Transform, Map *TileMap, anim *AnimRender, walktime int) Utils.Factory[*TileWalker] {
	return func() (*TileWalker, error) {
		t := TileWalker{
			Data:     tiledata,
			Position: pos,
			tf:       tf,
			tm:       Map,
			anim:     anim,
			walktime: Utils.NewReadOnly(walktime),
		}
		return &t, nil
	}
}
