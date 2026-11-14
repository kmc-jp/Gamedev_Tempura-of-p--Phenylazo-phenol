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

// SelectStageCanvasAPI は SelectStageScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
type SelectStageCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsSelectStageCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsSelectStageCanvas()
}

type SelectStageScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas SelectStageCanvasAPI

	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewSelectStageScene(worldDrawer Utils.Factory[WorldDrawer], canvas SelectStageCanvasAPI, Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*SelectStageScene] {
	return func() (*SelectStageScene, error) {
		camera, err := worldDrawer()
		if err != nil {

			return &SelectStageScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		if canvas == nil {
			return &SelectStageScene{}, fmt.Errorf("SelectStageSceneのcanvas が nil です")
		}
		ts := SelectStageScene{
			entities:    []Entity{},
			worldDrawer: camera,
			canvas:      canvas,
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
	s.canvas.Draw(screen)
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

	if err := s.canvas.Update(active); err != nil {
		return nil, transition.None(), fmt.Errorf("Canvasの更新に失敗しました。\n%w", err)
	}

	if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
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
		s.gamedata.SetNextStage(s.mkStage(Stage1))
		return s.mkTrainScene(), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionRed) {
		s.gamedata.SetNextStage(s.mkStage(Red))
		return s.mkTrainScene(), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionBlue) {
		s.gamedata.SetNextStage(s.mkStage(Blue))
		return s.mkTrainScene(), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionGreen) {
		s.gamedata.SetNextStage(s.mkStage(Green))
		return s.mkTrainScene(), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionOrange) {
		s.gamedata.SetNextStage(s.mkStage(Orange))
		return s.mkTrainScene(), transition.NewScene(), nil
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionPurple) {
		s.gamedata.SetNextStage(s.mkStage(Purple))
		return s.mkTrainScene(), transition.NewScene(), nil
	}

	return
}

func (s *SelectStageScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
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

func (s SelectStageScene) mkTrainScene() Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitTrainScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitTrainScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.Gamedata = s.gamedata
		nextScene.Gamedata.ForNewScene(game.Train_Default)
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}
