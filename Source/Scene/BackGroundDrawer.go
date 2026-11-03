package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

	"github.com/hajimehoshi/ebiten/v2"
)

type BackGroundDrawer interface {
	Draw(screen *ebiten.Image, gdata game.Gamedata)
	Update() error
}
