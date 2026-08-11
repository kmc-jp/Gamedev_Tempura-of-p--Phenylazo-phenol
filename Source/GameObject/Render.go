package gameobject

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"

	"github.com/hajimehoshi/ebiten/v2"
)

type Render struct {
}

func NewRender() (*Render, error) {
	return &Render{}, nil
}

func (r Render) Draw(obj scene.GameObject, scene *ebiten.Image) {

}
