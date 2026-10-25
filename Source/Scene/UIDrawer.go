package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UIDrawer interface {
	Draw(screen *ebiten.Image, objects []Element)
}
