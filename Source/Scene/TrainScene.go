package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type TrainScene struct {
	entities   []Entity
	drawer     Drawer
	text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewTrainScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TrainScene] {
	return func() (*TrainScene, error) {
		camera, err := NewDrawer()
		if err != nil {
			return &TrainScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		text, err := ui.NewText(
			Utils.Position{},
			color.Orange(),
			"電車で移動中\nESCキーでメニューを開く\nEnterキー で Result へ",
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		trs := TrainScene{
			entities:   []Entity{},
			drawer:     camera,
			text:       text,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
		}
		// なにもしない
		return &trs, nil
	}
}

func (s TrainScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *TrainScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
	s.text.Draw(screen)
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

	// if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
	if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
		return func() (game.Scene, error) {
			pauseScene, err := s.initScenes.InitPauseScene()
			if err != nil {
				return nil, err
			}
			pauseScene.gamedata = s.gamedata
			pauseScene.gamedata.InitialProcess = true
			pauseScene.initScenes = s.initScenes
			return pauseScene, nil
		}, transition.Push(), nil
	}

	if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
		return func() (game.Scene, error) {
			resultScene, err := s.initScenes.InitResultScene()
			if err != nil {
				return nil, err
			}
			resultScene.gamedata = s.gamedata
			resultScene.gamedata.ForNewScene(game.Result_Default)
			resultScene.initScenes = s.initScenes
			return resultScene, nil
		}, transition.NewScene(), nil
	}

	return
}

func (s TrainScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s TrainScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
