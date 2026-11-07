package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Element interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
}

type PositionElement interface {
	Element
	Position() Utils.Position
	SetPosition(Utils.Position)
}
