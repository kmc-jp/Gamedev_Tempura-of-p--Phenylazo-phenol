package component

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type StdRender struct {
}

func NewStdRender() (*StdRender, error) {
	return &StdRender{}, nil
}

// test
var _ Utils.Factory[*StdRender] = NewStdRender

func (r StdRender) Draw() (scene *ebiten.Image) {
	panic("Not implemented!")
}
