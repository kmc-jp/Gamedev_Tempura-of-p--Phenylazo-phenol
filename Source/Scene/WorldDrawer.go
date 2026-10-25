package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type WorldDrawer interface {
	Draw(screen *ebiten.Image, objects []Entity)
	SetCenter(Utils.Position)
}
