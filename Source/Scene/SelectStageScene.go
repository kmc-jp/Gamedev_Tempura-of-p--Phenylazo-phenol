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

type SelectStageScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	elements []Element
	uiDrawer UIDrawer

	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewSelectStageScene(worldDrawer Utils.Factory[WorldDrawer], uiDrawer Utils.Factory[UIDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*SelectStageScene] {
	return func() (*SelectStageScene, error) {
		camera, err := worldDrawer()
		if err != nil {

			return &SelectStageScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		uiCamera, err := uiDrawer()
		if err != nil {
			return &SelectStageScene{}, fmt.Errorf("uiDrawer() でエラーが発生しました。 \n%w", err)
		}
		ts := SelectStageScene{
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

func (s SelectStageScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *SelectStageScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.uiDrawer.Draw(screen, s.elements)
}

func (s *SelectStageScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

	if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
		return func() (game.Scene, error) {
			nextScene, err := s.initScenes.InitResultScene()
			if err != nil {
				return nil, fmt.Errorf("s.initScenes.InitResultScene() でエラーが発生しました。\n%w", err)
			}
			nextScene.gamedata = s.gamedata
			nextScene.gamedata.ForNewScene(game.Result_SelectParts)
			nextScene.initScenes = s.initScenes
			return nextScene, nil
		}, transition.NewScene(), nil
	}

	if s.Input.ActionIsJustPressed(inputconfig.Action1) {
		return s.mkStage(Stage1), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionRed) {
		return s.mkStage(Red), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionBlue) {
		return s.mkStage(Blue), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionGreen) {
		return s.mkStage(Green), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionOrange) {
		return s.mkStage(Orange), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionPurple) {
		return s.mkStage(Purple), transition.NewScene(), nil
	}

	return
}

func (s *SelectStageScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s *SelectStageScene) AddUIElement(e Element) {
	s.elements = append(s.elements, e)
}

func (s SelectStageScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}

func (s SelectStageScene) mkStage(stage Stage) Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitPlayScene[stage]()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitPlayScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.Gamedata = s.gamedata
		nextScene.Gamedata.ForNewScene(game.Play_Explore)
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}
