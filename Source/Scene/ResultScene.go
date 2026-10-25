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

type ResultScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	elements []Element
	uiDrawer UIDrawer
	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewResultScene(NewWorldDrawer Utils.Factory[WorldDrawer], uiDrawer Utils.Factory[UIDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*ResultScene] {
	return func() (*ResultScene, error) {
		camera, err := NewWorldDrawer()
		if err != nil {

			return &ResultScene{}, fmt.Errorf("NewDrawer() でエラーが発生しました。 \n%w", err)
		}
		ts := ResultScene{
			entities:    []Entity{},
			worldDrawer: camera,
			// text:        text,
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
	s.worldDrawer.Draw(screen, s.entities)
	s.uiDrawer.Draw(screen, s.elements)
	// s.text.Draw(screen)
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

func (s *ResultScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s *ResultScene) AddUIElement(e Element) {
	s.elements = append(s.elements, e)
}

func (s ResultScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
