package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// メインの遊べるシーン
type PlayScene struct {
	entities []Entity
	drawer   Drawer
	Input    *input.Handler
	// テスト用
	cameraCenter Utils.Position
}

func NewPlayScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler) Utils.Factory[*PlayScene] {
	camera, err := NewDrawer()
	if err != nil {
		return func() (*PlayScene, error) {
			return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
	}
	ps := PlayScene{
		entities: []Entity{},
		drawer:   camera,
		Input:    Input,
	}
	// なにもしない
	return func() (*PlayScene, error) {
		return &ps, nil
	}
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
}

// Update implements [Scene].
func (s *PlayScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	// 何か軽いテストはここでやりましょう
	for _, action := range inputconfig.ActionConstList() {
		if s.Input.ActionIsPressed(action) {
			switch action {
			case inputconfig.ActionMoveLeft:
				s.cameraCenter.Move(Utils.NewLeftVec2())
			case inputconfig.ActionMoveRight:
				s.cameraCenter.Move(Utils.NewRightVec2())
			case inputconfig.ActionMoveUp:
				s.cameraCenter.Move(Utils.NewUpVec2())
			case inputconfig.ActionMoveDown:
				s.cameraCenter.Move(Utils.NewDownVec2())
			default: // 何もしない
			}
		}
	}
	s.SetCameraCenter(s.cameraCenter)
	return
}

func (s *PlayScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s PlayScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
