package gameobject

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Render interface {
	Draw() (image *ebiten.Image, DIO *ebiten.DrawImageOptions, Layer string)
	Update() error
}
