package gameobject

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Render interface {
	Draw(obj GameObject, scene *ebiten.Image)
}
