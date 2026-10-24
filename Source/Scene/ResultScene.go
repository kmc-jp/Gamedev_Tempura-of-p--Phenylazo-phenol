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

type ResultScene struct {
	entities   []Entity
	drawer     Drawer
	text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewResultScene(NewDrawer Utils.Factory[Drawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*ResultScene] {
	return func() (*ResultScene, error) {
		camera, err := NewDrawer()
		if err != nil {

			return &ResultScene{}, fmt.Errorf("NewDrawer() でエラーが発生しました。 \n%w", err)
		}
		text, err := ui.NewText(
			Utils.Position{},
			color.Orange(),
			"Result\nEnter キーでタイトルに戻る",
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		ts := ResultScene{
			entities:   []Entity{},
			drawer:     camera,
			text:       text,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
		}
		// なにもしない
		return &ts, nil
	}
}

func (s ResultScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *ResultScene) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen, s.entities)
	s.text.Draw(screen)
}

func (s *ResultScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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
	if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
		return func() (game.Scene, error) {
			nextScene, err := s.initScenes.InitTitleScene()
			nextScene.gamedata = s.gamedata
			nextScene.gamedata.ForNewScene(game.Title_Default)
			nextScene.initScenes = s.initScenes
			if err != nil {
				return nil, err
			}
			return nextScene, nil
		}, transition.NewScene(), nil
	}

	return
}

func (s ResultScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s ResultScene) SetCameraCenter(center Utils.Position) {
	s.drawer.SetCenter(center)
}
