package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	selection "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Selection"
	"fmt"
	stdcolor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

// PauseCanvasAPI は PauseScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
type PauseCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsPauseCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsPauseCanvas()

	// ---- 入力用 ----

	// ShowTitleConfirmation はタイトルへ戻る確認を表示する。
	ShowTitleConfirmation()
	MoveTitleConfirmationLeft()
	MoveTitleConfirmationRight()
	ConfirmTitle()
	HideTitleConfirmation()

	// ---- 問い合わせ ----

	IsTitleConfirmationOpen() bool
	// TitleConfirmationResult は (結果, 確定済み) を返す。
	TitleConfirmationResult() (bool, bool)
}

// PauseMenuItem はポーズ画面のメニュー項目である。
type PauseMenuItem int

const (
	PauseMenuResume PauseMenuItem = iota
	PauseMenuTitle
	PauseMenuSettings
	PauseMenuHowToPlay
)

// PauseMenuItems はメニュー項目を表示順（上から）に返す。
func PauseMenuItems() []PauseMenuItem {
	return []PauseMenuItem{PauseMenuResume, PauseMenuTitle, PauseMenuSettings, PauseMenuHowToPlay}
}

type PauseScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas PauseCanvasAPI

	// selector はメニューのカーソル（選択中の項目）を管理する。
	selector *selection.Graph[PauseMenuItem, selection.Direction]

	// dimImage は背後の画面を暗くする半透明の画像。毎フレーム作らないようキャッシュし、画面サイズが変わったときだけ作り直す。
	dimImage *ebiten.Image

	// text        *ui.Text
	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

// NewPauseScene は PauseScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
func NewPauseScene(worldDrawer Utils.Factory[WorldDrawer], newCanvas func(*PauseScene) (PauseCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*PauseScene] {
	return func() (*PauseScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &PauseScene{}, fmt.Errorf("NewMockCamera() でエラーが発生しました。 \n%w", err)
		}
		if newCanvas == nil {
			return &PauseScene{}, fmt.Errorf("PauseSceneのcanvas生成関数が nil です")
		}

		// 上下に隣り合う項目へ移る。端では止まる（ループしない）。
		selector := selection.NewGraph[PauseMenuItem, selection.Direction](PauseMenuResume)
		items := PauseMenuItems()
		for i := 0; i+1 < len(items); i++ {
			selector.AddTransition(items[i], selection.Down, items[i+1], nil)
			selector.AddTransition(items[i+1], selection.Up, items[i], nil)
		}

		ps := PauseScene{
			entities:    []Entity{},
			worldDrawer: camera,
			selector:    selector,
			// text:        text,
			Input:      Input,
			gamedata:   gdata,
			initScenes: initScenes,
		}
		canvas, err := newCanvas(&ps)
		if err != nil {
			return &PauseScene{}, fmt.Errorf("PauseSceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &PauseScene{}, fmt.Errorf("PauseSceneのcanvas が nil です")
		}
		ps.canvas = canvas
		return &ps, nil
	}
}

// SelectedMenuItem はカーソルが指しているメニュー項目を返す。
func (s *PauseScene) SelectedMenuItem() PauseMenuItem {
	return s.selector.Current()
}

// SelectionRevision はカーソルが動くたびに増える整数を返す。
// Canvasはこの値の変化を観測し、変化したときに表示を更新する。
func (s *PauseScene) SelectionRevision() int {
	return s.selector.Revision()
}

func (s PauseScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *PauseScene) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if s.dimImage == nil || s.dimImage.Bounds().Dx() != w || s.dimImage.Bounds().Dy() != h {
		if s.dimImage != nil {
			s.dimImage.Deallocate()
		}
		s.dimImage = ebiten.NewImage(w, h)
		s.dimImage.Fill(stdcolor.RGBA{100, 100, 100, 200})
	}
	screen.DrawImage(s.dimImage, &ebiten.DrawImageOptions{})

	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *PauseScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

func (s *PauseScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	if s.gamedata.InitialProcess {
		//シーン起動時処理
		s.gamedata.InitialProcess = false

	}

	for _, entity := range s.entities {
		entity.Update(active)
	}

	// タイトルへ戻る確認の表示中は、確認だけを操作する
	if s.canvas.IsTitleConfirmationOpen() {
		if s.Input.ActionIsJustPressed(inputconfig.ActionMoveLeft) {
			s.canvas.MoveTitleConfirmationLeft()
		}
		if s.Input.ActionIsJustPressed(inputconfig.ActionMoveRight) {
			s.canvas.MoveTitleConfirmationRight()
		}
		if s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
			s.canvas.ConfirmTitle()
		}
		if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
			s.canvas.HideTitleConfirmation()
			return nil, transition.None(), nil
		}
		accepted, confirmed := s.canvas.TitleConfirmationResult()
		if confirmed {
			if accepted {
				return s.backToTitle(), transition.NewScene(), nil
			}
			s.canvas.HideTitleConfirmation()
		}
		return nil, transition.None(), nil
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
	if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
		return s.resume()
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
		case PauseMenuResume:
			return s.resume()
		case PauseMenuTitle:
			s.canvas.ShowTitleConfirmation()
		case PauseMenuSettings:
			// TODO: せってい画面ができたら遷移する
			fmt.Println("せってい が押されました")
		case PauseMenuHowToPlay:
			// TODO: あそびかた画面ができたら遷移する
			fmt.Println("あそびかた が押されました")
		}
	}

	return
}

// resume はポーズを解除してゲームへ戻る。
func (s *PauseScene) resume() (Utils.Factory[game.Scene], transition.Type, error) {
	s.gamedata.EnableWorldTimeProgressing()
	return nil, transition.Pop(), nil
}

// backToTitle はゲームの進行をリセットしてタイトルへ戻る。
func (s *PauseScene) backToTitle() Utils.Factory[game.Scene] {
	return func() (game.Scene, error) {
		nextScene, err := s.initScenes.InitTitleScene()
		if err != nil {
			return nil, fmt.Errorf("s.initScenes.InitTitleScene() でエラーが発生しました。\n%w", err)
		}
		nextScene.gamedata = s.gamedata
		nextScene.gamedata.ForNewScene(game.Title_Default)
		nextScene.gamedata.ResetGame()
		nextScene.initScenes = s.initScenes
		return nextScene, nil
	}
}

func (s *PauseScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s PauseScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
