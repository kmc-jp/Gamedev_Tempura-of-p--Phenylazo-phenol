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

type TitleScene struct {
	entities   []Entity
	drawer     Drawer
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewTitleScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TitleScene] {
	return func() (*TitleScene, error) {
		camera, err := NewDrawer()
		if err != nil {

			return &TitleScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		ts := TitleScene{
			entities:   []Entity{},
			drawer:     camera,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
		}
		// なにもしない
		return &ts, nil
	}
}

func (s TitleScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
	ui.DrawText(Utils.Position{}, color.Orange(),
		"Title\nPress Space Key to Start",
		ui.FontFace(s.gamedata.Font, 50), screen)
}

func (s *TitleScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

	// if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
		return func() (game.Scene, error) {
			nextScene, err := s.initScenes.InitPlayScene()
			nextScene.gamedata = s.gamedata
			nextScene.gamedata.InitialProcess = true
			nextScene.initScenes = s.initScenes
			if err != nil {
				return nil, err
			}
			return nextScene, nil
		}, transition.NewScene(), nil
	}

	return
}

func (s TitleScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s TitleScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
