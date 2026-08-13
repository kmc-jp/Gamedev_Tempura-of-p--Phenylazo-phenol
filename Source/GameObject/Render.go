package gameobject

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Render struct {
}

func NewRender() (*Render, error) {
	return &Render{}, nil
}

func (r Render) Draw(obj GameObject, scene *ebiten.Image) {

}
