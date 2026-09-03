package scene

import "github.com/hajimehoshi/ebiten/v2"

type Entity interface {
	Name() string
	Update(active bool) (err error)
	Draw(screen *ebiten.Image)
}
