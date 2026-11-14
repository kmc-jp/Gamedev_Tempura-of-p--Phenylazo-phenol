package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// TrainCanvas は TrainScene のUI構成である。案内テキストと所持パーツ一覧を持つ。
// 一覧を表示する時点はSceneが決めるため、Sceneのポインタを持って観測する。
type TrainCanvas struct {
	scene     *scene.TrainScene
	gdata     *game.Gamedata
	canvas    *canvas.Canvas
	partsList *ui.PartsListUI

	// revealed は所持パーツ一覧を構築・表示済みか。
	revealed bool
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
		utils.Position{},
		color.Orange(),
		"電車で移動中\nESCキーでメニューを開く\nEnterキー で 店に向かう",
		ui.FontFace(gdata.Font, 15),
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
	}

	partsList, err := newPartsList(gdata)
	if err != nil {
		return nil, err
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{guide, partsList} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("電車画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	return &TrainCanvas{scene: s, gdata: gdata, canvas: c, partsList: partsList}, nil
}

// Update implements [scene.TrainCanvasAPI].
// Sceneが一覧を表示する時点に達したら、所持パーツで一覧を構築して表示する。
func (c *TrainCanvas) Update(active bool) error {
	if !c.revealed && c.scene.IsPartsListRevealed() {
		if err := c.partsList.ShowParts(c.gdata.Parts); err != nil {
			return fmt.Errorf("パーツ一覧の構築に失敗しました。\n%w", err)
		}
		c.partsList.Show()
		c.revealed = true
	}
	return c.canvas.Update(active)
}

// Draw implements [scene.TrainCanvasAPI].
func (c *TrainCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsTrainCanvas implements [scene.TrainCanvasAPI].
func (c *TrainCanvas) IsTrainCanvas() {}

// ScrollPartsListDown implements [scene.TrainCanvasAPI].
// 選択枠を1つ下へ動かし、スクロールを追従させる。選択結果はどこにも使わない。
func (c *TrainCanvas) ScrollPartsListDown() {
	if c.partsList.Visible() {
		c.partsList.SelectNext()
	}
}

// ScrollPartsListUp implements [scene.TrainCanvasAPI].
func (c *TrainCanvas) ScrollPartsListUp() {
	if c.partsList.Visible() {
		c.partsList.SelectPrevious()
	}
}

var _ scene.TrainCanvasAPI = (*TrainCanvas)(nil)
