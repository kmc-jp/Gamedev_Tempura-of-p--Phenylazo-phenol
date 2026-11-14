package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// TrainCanvasAPI は TrainScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructからコンストラクタ経由で注入する。
// 所持パーツ一覧の表示は、Canvasが TrainScene.PartsListRevision を観測して行う。
type TrainCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsTrainCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsTrainCanvas()

	// ---- 入力用（一覧が非表示のときは何もしない） ----

	// ScrollPartsListDown は所持パーツ一覧を下へスクロールする。
	ScrollPartsListDown()
	// ScrollPartsListUp は所持パーツ一覧を上へスクロールする。
	ScrollPartsListUp()
}

type TrainScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas TrainCanvasAPI
	// partsListRevision は所持パーツ一覧を表示する時点に達するたびに増える。Canvasが観測する。
	// 0 の間は一覧を表示しない。
	partsListRevision int

	bgDrawer BackGroundDrawer

	// text       *ui.Text
	Input      *input.Handler
	Gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

// NewTrainScene は TrainScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
func NewTrainScene(worldDrawer Utils.Factory[WorldDrawer], bgDrawer Utils.Factory[BackGroundDrawer], newCanvas func(*TrainScene) (TrainCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TrainScene] {
	return func() (*TrainScene, error) {
		if newCanvas == nil {
			return &TrainScene{}, fmt.Errorf("TrainSceneのcanvas生成関数が nil です")
		}
		camera, err := worldDrawer()
		if err != nil {
			return &TrainScene{}, fmt.Errorf("worldDrawer() でエラーが発生しました。 \n%w", err)
		}
		bgCamera, err := bgDrawer()
		if err != nil {
			return &TrainScene{}, fmt.Errorf("BGDrawer() でエラーが発生しました。 \n%w", err)
		}
		trs := TrainScene{
			entities:    []Entity{},
			worldDrawer: camera,
			bgDrawer:    bgCamera,
			// text:       text,
			Input:      Input,
			Gamedata:   gdata,
			initScenes: initScenes,
		}
		canvas, err := newCanvas(&trs)
		if err != nil {
			return &TrainScene{}, fmt.Errorf("TrainSceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &TrainScene{}, fmt.Errorf("TrainSceneのcanvas が nil です")
		}
		trs.canvas = canvas
		return &trs, nil
	}
}

// PartsListRevision は所持パーツ一覧を表示する時点に達した回数を返す。
// 到着演出のタイマー満了で 1 になる。0 の間は一覧を表示しない。
func (s *TrainScene) PartsListRevision() int {
	return s.partsListRevision
}

func (s TrainScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *TrainScene) Draw(screen *ebiten.Image) {
	s.bgDrawer.Draw(screen, *s.Gamedata)
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *TrainScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if !active {
		return
	}

	nextScene, transitionType, err = s.updateScene(active)
	if err != nil {
		return nil, transition.None(), err
	}
	if err := s.canvas.Update(active); err != nil {
		return nil, transition.None(), fmt.Errorf("Canvasの更新に失敗しました。\n%w", err)
	}
	return
}

func (s *TrainScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.Gamedata.InitialProcess {
		//シーン起動時処理
		s.Gamedata.InitialProcess = false

		if s.Gamedata.Visited == 0 {
			return
		}

		s.Gamedata.Timer = tick.NewTimer()
		s.Gamedata.Timer.Set(50)
	}

	// 到着演出のタイマー満了時に1回だけ、一覧を表示する時点に達したことをCanvasへ伝える
	if s.Gamedata.Timer.Expired() && s.partsListRevision == 0 {
		s.partsListRevision++
	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	err = s.bgDrawer.Update()
	if err != nil {
		return nil, transition.None(), fmt.Errorf("s.bgDrawer.Update() でエラーが発生しました。\n%w", err)
	}

	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
		s.canvas.ScrollPartsListUp()
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
		s.canvas.ScrollPartsListDown()
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
		s.Gamedata.Visited++
		return s.Gamedata.NextStage(), transition.NewScene(), nil
	}

	return
}

func (s *TrainScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s TrainScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
