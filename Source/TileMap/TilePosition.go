package tilemap

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TilePosition struct {
	Position Utils.Vec2Int
	Layer    string
	Offset   Utils.Vec2Int
}

func (t TilePosition) ToPosition() Utils.Position {
	return Utils.Position{
		X: float64(t.Position.X),
		Y: float64(t.Position.Y),
	}
}

func (t *TilePosition) Move(Direction Utils.Vec2Int) {
	t.Position = t.Position.Add(Direction)
}
