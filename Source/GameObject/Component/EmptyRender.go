package component

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// StdRender を使いまわし
func NewEmptyRender() (*StdRender, error) {
	EmptyImage := ebiten.NewImage(1, 1)
	r := StdRender{
		img: EmptyImage,
	}
	return &r, nil
}
