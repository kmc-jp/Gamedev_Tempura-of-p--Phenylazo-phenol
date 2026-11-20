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

// SettingsCanvasAPI は SettingsScene が利用するCanvasの契約である。
// 具体的な実装はscenecanvas packageに置き、constructから注入する。
type SettingsCanvasAPI interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
	// IsSettingsCanvas は識別用メソッドで、処理を持たない。
	// メソッド集合が同じ他SceneのCanvasを誤って注入しないためのもの。
	IsSettingsCanvas()
}

// SettingsMenuItem は設定画面の項目である。
type SettingsMenuItem int

const (
	SettingsBGM SettingsMenuItem = iota
	SettingsSE
	SettingsBack
)

// SettingsMenuItems は項目を表示順（上から）に返す。
func SettingsMenuItems() []SettingsMenuItem {
	return []SettingsMenuItem{SettingsBGM, SettingsSE, SettingsBack}
}

// 音量変更時の確認用SE
const settingsSEPreviewPath = "爆発1.mp3"

type SettingsScene struct {
	entities    []Entity
	worldDrawer WorldDrawer

	canvas SettingsCanvasAPI

	selector *selection.Graph[SettingsMenuItem, selection.Direction]

	// dimImage は背後の画面を暗くする半透明の画像。画面サイズが変わったときだけ作り直す。
	dimImage *ebiten.Image

	Input      *input.Handler
	gamedata   *game.Gamedata
	initScenes *InitializedScenes
}

// NewSettingsScene は SettingsScene を生成する。
// CanvasはSceneを参照するため、生成関数 newCanvas を受け取り、Sceneの生成直後に呼んで注入する。
func NewSettingsScene(worldDrawer Utils.Factory[WorldDrawer], newCanvas func(*SettingsScene) (SettingsCanvasAPI, error), Input *input.Handler, gdata *game.Gamedata, initScenes *InitializedScenes) Utils.Factory[*SettingsScene] {
	return func() (*SettingsScene, error) {
		camera, err := worldDrawer()
		if err != nil {
			return &SettingsScene{}, fmt.Errorf("worldDrawer() でエラーが発生しました。 \n%w", err)
		}
		if newCanvas == nil {
			return &SettingsScene{}, fmt.Errorf("SettingsSceneのcanvas生成関数が nil です")
		}

		// 上下に隣り合う項目へ移る。端では止まる（ループしない）。
		selector := selection.NewGraph[SettingsMenuItem, selection.Direction](SettingsBGM)
		items := SettingsMenuItems()
		for i := 0; i+1 < len(items); i++ {
			selector.AddTransition(items[i], selection.Down, items[i+1], nil)
			selector.AddTransition(items[i+1], selection.Up, items[i], nil)
		}

		ss := SettingsScene{
			entities:    []Entity{},
			worldDrawer: camera,
			selector:    selector,
			Input:       Input,
			gamedata:    gdata,
			initScenes:  initScenes,
		}
		canvas, err := newCanvas(&ss)
		if err != nil {
			return &SettingsScene{}, fmt.Errorf("SettingsSceneのcanvas生成に失敗しました。\n%w", err)
		}
		if canvas == nil {
			return &SettingsScene{}, fmt.Errorf("SettingsSceneのcanvas が nil です")
		}
		ss.canvas = canvas
		return &ss, nil
	}
}

// SelectedMenuItem はカーソルが指している項目を返す。
func (s *SettingsScene) SelectedMenuItem() SettingsMenuItem {
	return s.selector.Current()
}

// SelectionRevision はカーソルが動くたびに増える整数を返す。
func (s *SettingsScene) SelectionRevision() int {
	return s.selector.Revision()
}

// BGMVolume は現在のBGM音量を返す。Canvasが表示更新のために観測する。
func (s *SettingsScene) BGMVolume() int { return s.gamedata.BGMVolume() }

// SEVolume は現在のSE音量を返す。Canvasが表示更新のために観測する。
func (s *SettingsScene) SEVolume() int { return s.gamedata.SEVolume() }

func (s SettingsScene) Name() string {
	return fmt.Sprintf("%T", s)
}

func (s *SettingsScene) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if s.dimImage == nil || s.dimImage.Bounds().Dx() != w || s.dimImage.Bounds().Dy() != h {
		if s.dimImage != nil {
			s.dimImage.Deallocate()
		}
		s.dimImage = ebiten.NewImage(w, h)
		s.dimImage.Fill(stdcolor.RGBA{60, 60, 60, 220})
	}
	screen.DrawImage(s.dimImage, &ebiten.DrawImageOptions{})

	s.worldDrawer.Draw(screen, s.entities)
	s.canvas.Draw(screen)
}

// Update はSceneの処理の後でCanvasを更新し、同じフレームの状態変化をCanvasが観測できるようにする。
func (s *SettingsScene) Update(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
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

func (s *SettingsScene) updateScene(active bool) (nextScene Utils.Factory[game.Scene], transitionType transition.Type, err error) {
	for _, entity := range s.entities {
		entity.Update(active)
	}

	if s.Input.ActionIsJustPressed(inputconfig.ActionEscape) {
		return nil, transition.Pop(), nil
	}

	// カーソル移動
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveUp) {
		s.selector.Fire(selection.Up)
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveDown) {
		s.selector.Fire(selection.Down)
	}

	// 音量変更
	delta := 0
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveLeft) {
		delta--
	}
	if s.Input.ActionIsJustPressed(inputconfig.ActionMoveRight) {
		delta++
	}
	if delta != 0 {
		if err := s.changeVolume(delta); err != nil {
			return nil, transition.None(), err
		}
	}

	// 決定
	if s.Input.ActionIsJustPressed(inputconfig.ActionSpace) || s.Input.ActionIsJustPressed(inputconfig.ActionEnter) {
		if s.selector.Current() == SettingsBack {
			return nil, transition.Pop(), nil
		}
	}

	return
}

// changeVolume は選択中の項目の音量を delta だけ変える。音量項目以外では何もしない。
func (s *SettingsScene) changeVolume(delta int) error {
	switch s.selector.Current() {
	case SettingsBGM:
		s.gamedata.SetBGMVolume(s.gamedata.BGMVolume() + delta)
	case SettingsSE:
		before := s.gamedata.SEVolume()
		s.gamedata.SetSEVolume(before + delta)
		if s.gamedata.SEVolume() == before {
			return nil // 端で変わらなかったときは鳴らさない
		}
		// 変更後の音量を確認できるよう、SEを鳴らす
		pl, err := s.gamedata.LoadSE(settingsSEPreviewPath)
		if err != nil {
			return fmt.Errorf("確認用SEの読み込みに失敗しました。\n%w", err)
		}
		pl.Play()
	}
	return nil
}

func (s *SettingsScene) AddEntity(e Entity) {
	s.entities = append(s.entities, e)
}

func (s SettingsScene) SetCameraCenter(center Utils.Position) {
	s.worldDrawer.SetCenter(center)
}
