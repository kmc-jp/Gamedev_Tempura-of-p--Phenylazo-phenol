package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	input "github.com/quasilyte/ebitengine-input"
)

type TrainScene struct {
	entities   []Entity
	drawer     Drawer
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewTrainScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TrainScene] {
	camera, err := NewDrawer()
	if err != nil {
		return func() (*TrainScene, error) {
			return &TrainScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
	}
	trs := TrainScene{
		entities:   []Entity{},
		drawer:     camera,
		Input:      Input,
		gamedata:   gdata,
		initScenes: initScenes,
	}
	// なにもしない
	return func() (*TrainScene, error) {
		return &trs, nil
	}
}

func (s TrainScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *TrainScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
	ui.DrawText(Utils.Position{}, color.Orange(),
		"Train\nESC to Open Pause Menu\n",
		ui.FontFace(s.gamedata.Font, 30), screen)
}

func (s *TrainScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	if s.gamedata.InitialProcess {
		//シーン起動時処理
		s.gamedata.InitialProcess = false

	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		pauseScene, err := s.initScenes.InitPauseScene()
		if err != nil {
			return nil, transition.None(), err
		}
		pauseScene.gamedata = s.gamedata
		pauseScene.initScenes = s.initScenes
		return func() (game.Scene, error) { return pauseScene, nil }, transition.Push(), nil
	}

	return
}

func (s TrainScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s TrainScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
