package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type UIDrawer interface {
	Draw(screen *ebiten.Image, objects []Element)
	SetCenter(Utils.Position)
}
