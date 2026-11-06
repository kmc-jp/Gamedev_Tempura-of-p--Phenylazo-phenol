package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	Name() string
	Update(active bool) (err error)
	Draw() (*ebiten.Image, *ebiten.DrawImageOptions, string)
	Position() Utils.Position
}
