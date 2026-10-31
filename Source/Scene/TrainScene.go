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

type TrainScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	elements  []Element
	uiDrawer  UIDrawer
	partsList partsListUI

	// text       *ui.Text
	Input      *input.Handler
	Gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewTrainScene(worldDrawer Utils.Factory[WorldDrawer], uiDrawer Utils.Factory[UIDrawer], Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TrainScene] {
	return func() (*TrainScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &TrainScene{}, fmt.Errorf("worldDrawer() でエラーが発生しました。 \n%w", err)
		}
		uiCamera, err := uiDrawer()
		if err != nil {
			return &TrainScene{}, fmt.Errorf("uiDrawer() でエラーが発生しました。 \n%w", err)
		}
		trs := TrainScene{
			entities:    []Entity{},
			worldDrawer: camera,
			elements:    []Element{},
			uiDrawer:    uiCamera,
			// text:       text,
			Input:      Input,
			Gamedata:   gdata,
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
	s.worldDrawer.Draw(screen, s.entities)
	s.uiDrawer.Draw(screen, s.elements)
}

func (s *TrainScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	if s.Gamedata.InitialProcess {
		//シーン起動時処理
		s.Gamedata.InitialProcess = false

		if err := s.partsList.SetParts(s.Gamedata.Parts); err != nil {
			return nil, transition.None(), fmt.Errorf("パーツ一覧の構築に失敗しました。\n%w", err)
		}
		s.partsList.SetVisible(true)
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
			pauseScene.gamedata = s.Gamedata
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
			resultScene.gamedata = s.Gamedata
			resultScene.gamedata.ForNewScene(game.Result_SelectParts)
			resultScene.initScenes = s.initScenes
			return resultScene, nil
		}, transition.NewScene(), nil
	}

	return
}

func (s *TrainScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s *TrainScene) AddUIElement(e Element) {
	s.elements = append(s.elements, e)
}

func (s *TrainScene) SetPartsList(partsList partsListUI) {
	s.partsList = partsList
	s.AddUIElement(partsList)
}

func (s TrainScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
