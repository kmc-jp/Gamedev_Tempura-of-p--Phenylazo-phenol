package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	input "github.com/quasilyte/ebitengine-input"
)

// メインの遊べるシーン
type PlayScene struct {
	entities   []Entity
	drawer     Drawer
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
	// テスト用
	cameraCenter Utils.Position
}

func NewPlayScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*PlayScene] {
	camera, err := NewDrawer()
	if err != nil {
		return func() (*PlayScene, error) {
			return &PlayScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
	}
	ps := PlayScene{
		entities:   []Entity{},
		drawer:     camera,
		Input:      Input,
		gamedata:   gdata,
		initScenes: initScenes,
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
	if !active {
		return
	}

	if s.gamedata.InitialProcess {
		s.gamedata.InitialProcess = false

	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	// 何か軽いテストはここでやりましょう

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		pauseScene, err := s.initScenes.InitPauseScene()
		if err != nil {
			return nil, transition.None(), err
		}
		pauseScene.gamedata = s.gamedata
		pauseScene.initScenes = s.initScenes
		return func() (game.Scene, error) { return pauseScene, nil }, transition.Push(), nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		trainScene, err := s.initScenes.InitTrainScene()
		if err != nil {
			return nil, transition.None(), err
		}
		trainScene.gamedata = s.gamedata
		trainScene.initScenes = s.initScenes
		return func() (game.Scene, error) { return trainScene, nil }, transition.NewScene(), nil
	}

	return
}

func (s *PlayScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s PlayScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
