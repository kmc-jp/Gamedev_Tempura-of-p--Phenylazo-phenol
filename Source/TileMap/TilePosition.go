package tilemap

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TilePosition struct {
	X int
	Y int
}

func (t TilePosition) ToPosition() Utils.Position {
	return Utils.Position{
		X: float64(t.X),
		Y: float64(t.Y),
	}
}
