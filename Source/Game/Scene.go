package game

import (
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Name() string
	Update(active bool) (nextScene SceneFactory, transitionType transition.Type, err error)
	Draw(screen *ebiten.Image)
}

type SceneFactory func() (Scene, error)
