package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	Name() string
	Update(active bool) (err error)
	Draw() (image *ebiten.Image)
	Position() Utils.Position
}
