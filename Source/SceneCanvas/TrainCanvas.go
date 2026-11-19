package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	commonui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/CommonUI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// TrainCanvas は TrainScene のUI構成である。案内テキストと所持パーツ一覧を持つ。
// 一覧を表示する時点はSceneが決めるため、Sceneのポインタを持って観測する。
type TrainCanvas struct {
	scene       *scene.TrainScene
	gdata       *game.Gamedata
	canvas      *canvas.Canvas
	partsList   *ui.PartsListUI
	sellConfirm *commonui.Confirm

	// shownListRevision は反映済みの TrainScene.PartsListRevision。
	shownListRevision int
	// shownPartsRevision は一覧に反映済みの Gamedata.PartsRevision。
	shownPartsRevision int
}

// NewTrainCanvas は TrainCanvas を生成する。
func NewTrainCanvas(s *scene.TrainScene, gdata *game.Gamedata) (*TrainCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("TrainScene が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}

	guide, err := ui.NewText(
		utils.Position{Y: 25},
		color.Orange(),
		"電車で移動中\n上下キーでパーツ選択\nSpaceキーで売却確認\nEnterキーで店に向かう\nESCキーでメニューを開く",
		ui.FontFace(gdata.Font, 15),
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
	}

	worldTimeText, err := commonui.NewClock(
		utils.Position{},
		color.White(),
		gdata,
		ui.FontFace(gdata.Font, 15),
	)
	if err != nil {
		return nil, fmt.Errorf("commonui.NewClock() でエラーが発生しました。\n%w", err)
	}

	partsList, err := newPartsList(gdata)
	if err != nil {
		return nil, err
	}

	sellConfirm, err := commonui.NewConfirm(commonui.ConfirmConfig{
		Size:                  image.Point{X: 280, Y: 110},
		Font:                  ui.FontFace(gdata.Font, 20),
		TextColor:             color.White(),
		BackgroundColor:       color.SemiTransparentBlack(),
		ButtonBackgroundColor: color.DarkGray(),
		Padding:               image.Point{X: 12, Y: 12},
		ButtonSize:            image.Point{X: 90, Y: 30},
		Message:               "このパーツを売りますか？",
		YesLabel:              "はい",
		NoLabel:               "いいえ",
	})
	if err != nil {
		return nil, fmt.Errorf("commonui.NewConfirm() でエラーが発生しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{worldTimeText, guide, partsList, sellConfirm} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("電車画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	return &TrainCanvas{
		scene:       s,
		gdata:       gdata,
		canvas:      c,
		partsList:   partsList,
		sellConfirm: sellConfirm,
	}, nil
}

// Update implements [scene.TrainCanvasAPI].
// Sceneが一覧を表示する時点に達したら（PartsListRevision の変化）、所持パーツで一覧を構築して表示する。
// 表示後に所持パーツが変化した（PartsRevision の変化）ときも作り直す。
func (c *TrainCanvas) Update(active bool) error {
	listRevision := c.scene.PartsListRevision()
	partsRevision := c.gdata.PartsRevision()
	if listRevision > 0 && (listRevision != c.shownListRevision || partsRevision != c.shownPartsRevision) {
		if err := c.partsList.ShowSellParts(c.gdata.OwnedParts()); err != nil {
			return fmt.Errorf("パーツ一覧の構築に失敗しました。\n%w", err)
		}
		c.partsList.Show()
		c.shownListRevision = listRevision
		c.shownPartsRevision = partsRevision
	}
	return c.canvas.Update(active)
}

// Draw implements [scene.TrainCanvasAPI].
func (c *TrainCanvas) Draw(screen *ebiten.Image) {
	if c.sellConfirm.Visible() {
		bounds := screen.Bounds()
		confirmSize := c.sellConfirm.Size()
		c.sellConfirm.SetPosition(utils.Position{
			X: float64((bounds.Dx() - confirmSize.X) / 2),
			Y: float64((bounds.Dy() - confirmSize.Y) / 2),
		})
	}
	c.canvas.Draw(screen)
}

// IsTrainCanvas implements [scene.TrainCanvasAPI].
func (c *TrainCanvas) IsTrainCanvas() {}

// ScrollPartsListDown implements [scene.TrainCanvasAPI].
// 選択枠を1つ下へ動かし、スクロールを追従させる。選択結果はどこにも使わない。
// 一覧の非表示中は PartsListUI が何もしない。
func (c *TrainCanvas) ScrollPartsListDown() {
	c.partsList.SelectNext()
}

// ScrollPartsListUp implements [scene.TrainCanvasAPI].
func (c *TrainCanvas) ScrollPartsListUp() {
	c.partsList.SelectPrevious()
}

func (c *TrainCanvas) ShowNextPartsTab() error {
	return c.partsList.NextTab()
}

func (c *TrainCanvas) ShowPreviousPartsTab() error {
	return c.partsList.PreviousTab()
}

func (c *TrainCanvas) IsNoSellTabSelected() bool {
	return c.partsList.Visible() && c.partsList.ActiveTab() == ui.PartsListNoSell
}

func (c *TrainCanvas) ShowSellPartConfirmation() bool {
	if c.partsList.SelectedPart() == nil {
		return false
	}
	c.sellConfirm.Show()
	return true
}

func (c *TrainCanvas) MoveSellPartConfirmationLeft() {
	c.sellConfirm.MoveLeft()
}

func (c *TrainCanvas) MoveSellPartConfirmationRight() {
	c.sellConfirm.MoveRight()
}

func (c *TrainCanvas) ConfirmSellPart() {
	c.sellConfirm.Confirm()
}

func (c *TrainCanvas) IsSellPartConfirmationOpen() bool {
	return c.sellConfirm.Visible()
}

func (c *TrainCanvas) SellPartConfirmationResult() (bool, bool) {
	return c.sellConfirm.Result()
}

func (c *TrainCanvas) HideSellPartConfirmation() {
	c.sellConfirm.Hide()
}

func (c *TrainCanvas) SelectedPart() *game.Part {
	return c.partsList.SelectedPart()
}

var _ scene.TrainCanvasAPI = (*TrainCanvas)(nil)
