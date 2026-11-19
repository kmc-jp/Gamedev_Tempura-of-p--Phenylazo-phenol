package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	commonui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/CommonUI"
	pauseui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/PauseUI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// 以下はポーズ画面の見た目の調整値。
const (
	pauseTitleTextSize    = 24
	pauseTitleMarginY     = 30
	pauseMenuTextSize     = 20
	pauseMenuButtonWidth  = 200
	pauseMenuButtonHeight = 30
	pauseMenuSpacing      = 10
	pauseConfirmWidth     = 280
	pauseConfirmHeight    = 110
	pauseConfirmPadding   = 12
	pauseConfirmBtnWidth  = 90
	pauseConfirmBtnHeight = 30
)

// pauseMenuLabels はメニュー項目ごとのボタンの文言。
var pauseMenuLabels = map[scene.PauseMenuItem]string{
	scene.PauseMenuResume:    "ゲームにもどる",
	scene.PauseMenuTitle:     "タイトルにもどる",
	scene.PauseMenuSettings:  "せってい",
	scene.PauseMenuHowToPlay: "あそびかた",
}

// PauseCanvas は PauseScene のUI構成である。見出し、縦に並ぶメニューボタン、タイトルへ戻る確認を持つ。
// 背景の半透明の塗りつぶし（グレー化）は、従来どおり PauseScene.Draw が行う。
// 選択中の項目はSceneが持つため、Sceneのポインタを持って観測する。
type PauseCanvas struct {
	scene        *scene.PauseScene
	canvas       *canvas.Canvas
	title        *ui.Text
	menu         *pauseui.PauseMenu
	titleConfirm *commonui.Confirm

	items         []scene.PauseMenuItem
	shownRevision int
}

// NewPauseCanvas は PauseCanvas を生成する。
func NewPauseCanvas(s *scene.PauseScene, gdata *game.Gamedata) (*PauseCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("PauseScene が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}
	if gdata.Font == nil {
		return nil, fmt.Errorf("gdata.Font が nil です")
	}

	title, err := ui.NewText(utils.Position{}, color.Orange(), "ポーズ", ui.FontFace(gdata.Font, pauseTitleTextSize))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
	}

	items := scene.PauseMenuItems()
	labels := make([]string, 0, len(items))
	for _, item := range items {
		label, ok := pauseMenuLabels[item]
		if !ok {
			return nil, fmt.Errorf("メニュー項目 %d の文言が pauseMenuLabels にありません", item)
		}
		labels = append(labels, label)
	}
	menu, err := pauseui.NewPauseMenu(pauseui.PauseMenuConfig{
		Labels:                labels,
		Font:                  ui.FontFace(gdata.Font, pauseMenuTextSize),
		TextColor:             color.White(),
		ButtonBackgroundColor: color.DarkGray(),
		ButtonSize:            image.Point{X: pauseMenuButtonWidth, Y: pauseMenuButtonHeight},
		Spacing:               pauseMenuSpacing,
	})
	if err != nil {
		return nil, fmt.Errorf("pauseui.NewPauseMenu() でエラーが発生しました。\n%w", err)
	}

	titleConfirm, err := commonui.NewConfirm(commonui.ConfirmConfig{
		Size:                  image.Point{X: pauseConfirmWidth, Y: pauseConfirmHeight},
		Font:                  ui.FontFace(gdata.Font, pauseMenuTextSize),
		TextColor:             color.White(),
		BackgroundColor:       color.LightGray(),
		ButtonBackgroundColor: color.DarkGray(),
		Padding:               image.Point{X: pauseConfirmPadding, Y: pauseConfirmPadding},
		ButtonSize:            image.Point{X: pauseConfirmBtnWidth, Y: pauseConfirmBtnHeight},
		Message:               "タイトルにもどりますか？",
		YesLabel:              "はい",
		NoLabel:               "いいえ",
	})
	if err != nil {
		return nil, fmt.Errorf("commonui.NewConfirm() -> titleConfirm でエラーが発生しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{title, menu, titleConfirm} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("ポーズ画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	pc := &PauseCanvas{scene: s, canvas: c, title: title, menu: menu, titleConfirm: titleConfirm, items: items}
	pc.syncWithScene()
	return pc, nil
}

// syncWithScene は Scene の選択中の項目をメニューへ反映する。
func (c *PauseCanvas) syncWithScene() {
	selected := c.scene.SelectedMenuItem()
	for i, item := range c.items {
		if item == selected {
			c.menu.SetSelected(i)
			break
		}
	}
	c.shownRevision = c.scene.SelectionRevision()
}

// Update implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) Update(active bool) error {
	if c.scene.SelectionRevision() != c.shownRevision {
		c.syncWithScene()
	}
	return c.canvas.Update(active)
}

// Draw implements [scene.PauseCanvasAPI].
// 画面サイズに合わせて見出し・メニュー・確認を中央へ配置する。
func (c *PauseCanvas) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	menuSize := c.menu.Size()
	menuY := float64(bounds.Dy()-menuSize.Y) / 2
	c.menu.SetPosition(utils.Position{
		X: float64(bounds.Dx()-menuSize.X) / 2,
		Y: menuY,
	})
	c.title.SetPosition(utils.Position{
		X: float64(bounds.Dx()/2) - pauseTitleTextSize*1.5,
		Y: menuY - pauseTitleMarginY - pauseTitleTextSize,
	})
	confirmSize := c.titleConfirm.Size()
	c.titleConfirm.SetPosition(utils.Position{
		X: float64(bounds.Dx()-confirmSize.X) / 2,
		Y: float64(bounds.Dy()-confirmSize.Y) / 2,
	})
	c.canvas.Draw(screen)
}

// IsPauseCanvas implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) IsPauseCanvas() {}

// ShowTitleConfirmation implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) ShowTitleConfirmation() { c.titleConfirm.Show() }

// MoveTitleConfirmationLeft implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) MoveTitleConfirmationLeft() { c.titleConfirm.MoveLeft() }

// MoveTitleConfirmationRight implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) MoveTitleConfirmationRight() { c.titleConfirm.MoveRight() }

// ConfirmTitle implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) ConfirmTitle() { c.titleConfirm.Confirm() }

// HideTitleConfirmation implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) HideTitleConfirmation() { c.titleConfirm.Hide() }

// IsTitleConfirmationOpen implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) IsTitleConfirmationOpen() bool { return c.titleConfirm.Visible() }

// TitleConfirmationResult implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) TitleConfirmationResult() (bool, bool) { return c.titleConfirm.Result() }

var _ scene.PauseCanvasAPI = (*PauseCanvas)(nil)
