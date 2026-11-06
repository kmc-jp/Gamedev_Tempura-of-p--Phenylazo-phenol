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

	isMoving  bool
	elapsed   int
	direction Utils.Vec2Int
	dirTag    map[Utils.Vec2Int]string
}

// test
var _ gameobject.Component = &TileWalker{}

func (w *TileWalker) Update(active bool) (err error) {
	if !active {
		return nil
	}

	if w.isMoving {
		w.elapsed++
		if w.elapsed > w.walktime.Value() {
			w.moveStop()
		} else {
			diff := w.tm.TileSize.ToVec2().HadamardProd(w.direction.ToVec2()).ScalarProd(1 / float64(w.walktime.Value()))
			w.Position.Offset = w.Position.Offset.Add(diff)
		}
	}

	w.tf.SetPosition(w.tm.TilePositionToWorldPosition(w.Position, w.Data))
	return nil
}

func (w *TileWalker) TryMove(Direction Utils.Vec2Int) bool {
	moveto := w.Position.Position.Add(Direction)
	if enterable, ok := w.tm.Enterable[moveto]; (enterable || !ok) && !w.isMoving {
		w.isMoving = true
		w.direction = Direction
		w.anim.SetTag(w.dirTag[Direction])
		w.anim.SetPlaying(true)
		return true
	} else {
		// 向きを変えるだけ
		w.direction = Direction
		w.anim.SetTag(w.dirTag[Direction])
		return false
	}
}

func (w *TileWalker) moveStop() {
	w.isMoving = false
	w.elapsed = 0
	w.Position.Offset = Utils.NewZeroVec2()
	w.Position.Position = w.Position.Position.Add(w.direction)
	w.anim.SetPlaying(false)
}

func (w TileWalker) IsMoving() bool {
	return w.isMoving
}

func NewTileWalker(tiledata tilemap.TileData, pos tilemap.TilePosition, tf gameobject.Transform, Map *TileMap, anim *AnimRender, walktime int, dirTag map[Utils.Vec2Int]string) Utils.Factory[*TileWalker] {
	return func() (*TileWalker, error) {
		t := TileWalker{
			Data:     tiledata,
			Position: pos,
			tf:       tf,
			tm:       Map,
			anim:     anim,
			walktime: Utils.NewReadOnly(walktime),
			dirTag:   dirTag,
		}
		return &t, nil
	}
}
