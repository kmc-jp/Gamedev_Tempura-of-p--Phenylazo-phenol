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
	objects []Entity
	drawer  Drawer
}

func NewPlayScene() (*PlayScene, error) {
	camera, err := NewMockCamera()
	if err != nil {

		return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
	}
	ps := PlayScene{
		objects: []Entity{},
		drawer:  camera,
	}
	// なにもしない
	return &ps, nil
}

// test
var _ Utils.Factory[PlayScene] = NewPlayScene

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.objects, Utils.Position{})
}

// Update implements [Scene].
func (s PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	// なにもしない
	return
}
