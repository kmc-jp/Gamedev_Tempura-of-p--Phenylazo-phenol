package game

import (
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Name() string
	Update(active bool) (nextScene Utils.Factory[Scene], transitionType transition.Type, err error)
	Draw(screen *ebiten.Image)
}
