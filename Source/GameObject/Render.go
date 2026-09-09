package gameobject

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Render interface {
	Draw() (image *ebiten.Image)
}
