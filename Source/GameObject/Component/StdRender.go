package component

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type StdRender struct {
}

func NewRender() (*StdRender, error) {
	return &StdRender{}, nil
}

// test
var _ Utils.Factory[StdRender] = NewRender

func (r StdRender) Draw() (scene *ebiten.Image) {
	panic("Not implemented!")
}
