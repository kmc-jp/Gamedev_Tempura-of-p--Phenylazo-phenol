package scene

import "github.com/hajimehoshi/ebiten/v2"

type Element interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
}
