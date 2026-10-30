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

type TitleScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	elements []Element
	uiDrawer UIDrawer

	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewTitleScene(worldDrawer Utils.Factory[WorldDrawer], uiDrawer Utils.Factory[UIDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TitleScene] {
	return func() (*TitleScene, error) {
		camera, err := worldDrawer()
		if err != nil {

			return &TitleScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		uiCamera, err := uiDrawer()
		if err != nil {
			return &TitleScene{}, fmt.Errorf("uiDrawer() でエラーが発生しました。 \n%w", err)
		}
		ts := TitleScene{
			entities:    []Entity{},
			worldDrawer: camera,
			elements:    []Element{},
			uiDrawer:    uiCamera,
			// text:       text,
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
	s.worldDrawer.Draw(screen, s.entities)
	s.uiDrawer.Draw(screen, s.elements)
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
			if err != nil {
				return nil, fmt.Errorf("s.initScenes.InitPlayScene() でエラーが発生しました。\n%w", err)
			}
			nextScene.Gamedata = s.gamedata
			nextScene.Gamedata.ForNewScene(game.Play_Explore)
			nextScene.initScenes = s.initScenes
			return nextScene, nil
		}, transition.NewScene(), nil
	}

	return
}

func (s *TitleScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s *TitleScene) AddUIElement(e Element) {
	s.elements = append(s.elements, e)
}

func (s TitleScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
