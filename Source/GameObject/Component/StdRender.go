package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"

	"github.com/hajimehoshi/ebiten/v2"
)

type StdRender struct {
}

func NewRender() (*StdRender, error) {
	return &StdRender{}, nil
}

func (r StdRender) Draw(obj gameobject.GameObject, scene *ebiten.Image) {

}
