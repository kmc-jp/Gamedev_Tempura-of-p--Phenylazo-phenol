package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// メインの遊べるシーン
type PlayScene struct {
}

func NewPlayScene() (PlayScene, error) {
	ps := PlayScene{}
	// なにもしない
	return ps, nil
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	// なにもしない
}

// Update implements [Scene].
func (s PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	// なにもしない
	return
}
