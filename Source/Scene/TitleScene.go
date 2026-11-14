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

// TitleCanvasAPI は TitleScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
type TitleCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsTitleCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsTitleCanvas()
}

type TitleScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas TitleCanvasAPI

	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

func NewTitleScene(worldDrawer Utils.Factory[WorldDrawer], canvas TitleCanvasAPI, Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TitleScene] {
	return func() (*TitleScene, error) {
		camera, err := worldDrawer()
		if err != nil {

			return &TitleScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		if canvas == nil {
			return &TitleScene{}, fmt.Errorf("TitleSceneのcanvas が nil です")
		}
		ts := TitleScene{
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

func (s TitleScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
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

	if err := s.canvas.Update(active); err != nil {
		return nil, transition.None(), fmt.Errorf("Canvasの更新に失敗しました。\n%w", err)
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) {
		return func() (game.Scene, error) {
			nextScene, err := s.initScenes.InitSelectStageScene()
			if err != nil {
				return nil, fmt.Errorf("s.initScenes.InitSelectStageScene() でエラーが発生しました。\n%w", err)
			}
			nextScene.gamedata = s.gamedata
			nextScene.gamedata.ForNewScene(game.SelectStage_Default)
			nextScene.initScenes = s.initScenes
			return nextScene, nil
		}, transition.NewScene(), nil
	}

	return
}

func (s *TitleScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s TitleScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}

func (s TitleScene) mkStage(stage Stage) Utils.Factory[game.Scene] {
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
