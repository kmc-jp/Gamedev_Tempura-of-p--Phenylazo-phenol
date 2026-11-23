package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	selection "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Selection"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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

// TitleMenuItem はタイトル画面のメニュー項目である。
type TitleMenuItem int

const (
	TitleMenuStart TitleMenuItem = iota
	TitleMenuHowToPlay
	TitleMenuRanking
	TitleMenuSettings
)

// TitleMenuItems はメニュー項目を表示順に返す。
func TitleMenuItems() []TitleMenuItem {
	return []TitleMenuItem{TitleMenuStart, TitleMenuHowToPlay, TitleMenuRanking, TitleMenuSettings}
}

type TitleScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas TitleCanvasAPI

	bgDrawer BackGroundDrawer

	// selector はメニューのカーソル（選択中の項目）を管理する。
	selector *selection.Graph[TitleMenuItem, selection.Direction]

	// text       *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes

	bgmPlayer *audio.Player
}

// NewTitleScene は TitleScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
func NewTitleScene(worldDrawer Utils.Factory[WorldDrawer], bgDrawer Utils.Factory[BackGroundDrawer], newCanvas func(*TitleScene) (TitleCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*TitleScene] {
	return func() (*TitleScene, error) {
		camera, err := worldDrawer()
		if err != nil {

			return &TitleScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		if newCanvas == nil {
			return &TitleScene{}, fmt.Errorf("TitleSceneのcanvas生成関数が nil です")
		}
		bgCamera, err := bgDrawer()
		if err != nil {
			return &TitleScene{}, fmt.Errorf("BGDrawer() でエラーが発生しました。 \n%w", err)
		}

		// 上下に隣り合う項目へ移る。端では止まる（ループしない）。
		selector := selection.NewGraph[TitleMenuItem, selection.Direction](TitleMenuStart)
		items := TitleMenuItems()
		for i := 0; i+1 < len(items); i++ {
			selector.AddTransition(items[i], selection.Down, items[i+1], nil)
			selector.AddTransition(items[i+1], selection.Up, items[i], nil)
		}

		const bgmPath = "bgm/title.mp3"
		bgmpl, err := gdata.LoadBGM(bgmPath)
		if err != nil {
			return &TitleScene{}, fmt.Errorf("gdata.LoadBGM() でエラーが発生しました。\npath : %s\n%w", bgmPath, err)
		}
		bgmpl.Play()

		ts := TitleScene{
			entities:    []Entity{},
			worldDrawer: camera,
			bgDrawer:    bgCamera,
			selector:    selector,
			// text:       text,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
			bgmPlayer:  bgmpl,
		}
		canvas, err := newCanvas(&ts)
		if err != nil {
			return &TitleScene{}, fmt.Errorf("TitleSceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &TitleScene{}, fmt.Errorf("TitleSceneのcanvas が nil です")
		}
		ts.canvas = canvas
		return &ts, nil
	}
}

// SelectedMenuItem はカーソルが指しているメニュー項目を返す。
func (s *TitleScene) SelectedMenuItem() TitleMenuItem {
	return s.selector.Current()
}

// SelectionRevision はカーソルが動くたびに増える整数を返す。
// Canvasはこの値の変化を観測し、変化したときに表示を更新する。
func (s *TitleScene) SelectionRevision() int {
	return s.selector.Revision()
}

func (s TitleScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	s.bgDrawer.Draw(screen, *s.gamedata)
	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *TitleScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

func (s *TitleScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.gamedata.InitialProcess {
		//シーン起動時処理
		s.gamedata.InitialProcess = false

	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	// カーソル移動
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
		s.selector.Fire(selection.Up)
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
		s.selector.Fire(selection.Down)
	}

	// 決定
	if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) || s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
		switch s.selector.Current() {
		case TitleMenuStart:
			return s.startGame(), transition.NewScene(), nil
		case TitleMenuHowToPlay:
			return s.openHowToPlay(), transition.Push(), nil
		case TitleMenuRanking:
			return s.openRanking(), transition.NewScene(), nil
		case TitleMenuSettings:
			return s.openSettings(), transition.Push(), nil
		}
	}

	return
}

// openHowToPlay は遊び方画面を重ねて開く。
func (s *TitleScene) openHowToPlay() Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		howToPlayScene, err := s.initScenes.InitHowToPlayScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitHowToPlayScene() でエラーが発生しました。\n%w", err)
		}
		howToPlayScene.gamedata = s.gamedata
		howToPlayScene.initScenes = s.initScenes
		return howToPlayScene, nil
	}
}

// openRanking はランキング画面へ移る。タイトルBGMはランキング画面から戻るときに作り直される。
func (s *TitleScene) openRanking() Utils.Factory[game.Scene] {
	s.bgmPlayer.Pause()

	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitRankingScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitRankingScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.gamedata = s.gamedata
		nextScene.gamedata.ForNewScene(game.Ranking_View)
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}

// openSettings は設定画面を重ねて開く。タイトルBGMは鳴らしたままにする。
func (s *TitleScene) openSettings() Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		settingsScene, err := s.initScenes.InitSettingsScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitSettingsScene() でエラーが発生しました。\n%w", err)
		}
		settingsScene.gamedata = s.gamedata
		settingsScene.initScenes = s.initScenes
		return settingsScene, nil
	}
}

// startGame はステージ選択画面へ向かう。
func (s *TitleScene) startGame() Utils.Factory[game.Scene] {
	// bgmを止める
	s.bgmPlayer.Pause()

	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitSelectStageScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitSelectStageScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.gamedata = s.gamedata
		nextScene.gamedata.ForNewScene(game.SelectStage_Select)
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}

// backToTitleFactory は全BGMを止め、ゲームデータを初期化してタイトルへ戻る。
// 初期化の後でTitleSceneを生成するので、タイトルBGMは止められない。
func backToTitleFactory(gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		gdata.StopAllBGM()
		gdata.ResetGame()
		gdata.ForNewScene(game.Title_Default)
		nextScene, err := initScenes.InitTitleScene()
		if err != nil {
			return nil, fmt.Errorf("initScenes.InitTitleScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.gamedata = gdata
		nextScene.initScenes = initScenes
		return nextScene, nil
	}
}

func (s *TitleScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s TitleScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}

func (s TitleScene) mkStage(stage stageinfo.Location) Utils.Factory[game.Scene] {
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
