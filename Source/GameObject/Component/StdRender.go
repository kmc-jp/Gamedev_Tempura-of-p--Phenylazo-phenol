package component

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type StdRender struct {
}

func NewRender() (*StdRender, error) {
	return &StdRender{}, nil
}

func (r StdRender) Draw() (scene *ebiten.Image) {
	panic("Not implemented!")
}
