package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	ebitenutil.DebugPrint(screen, "PlayScene\nP to pause")
}

// Update implements [Scene].
func (s PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return nil, transition.None(), nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		return func() (game.Scene, error) {
			return PauseScene{}, nil
		}, transition.Push(), nil
	}
	return nil, transition.None(), nil
}
