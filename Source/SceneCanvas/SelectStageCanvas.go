package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// SelectStageCanvas は SelectStageScene のUI構成である。ステージ選択の案内テキストだけを持つ。
type SelectStageCanvas struct {
	canvas *canvas.Canvas
}

// NewSelectStageCanvas は SelectStageCanvas を生成する。
func NewSelectStageCanvas(font *text.GoTextFaceSource) (*SelectStageCanvas, error) {
	c, err := newTextCanvas(
		utils.Position{},
		"ステージ選択\nPress 1 Key to Start Stage1\nPress 2 Key to Start RedStage\nPress 3 Key to Start BlueStage\nPress 4 Key to Start GreenStage\nPress 5 Key to Start OrangeStage\nPress 6 Key to Start PurpleStage\nスペースキーで帰宅",
		font,
		15,
	)
	if err != nil {
		return nil, fmt.Errorf("ステージ選択画面のCanvas構築に失敗しました。\n%w", err)
	}
	return &SelectStageCanvas{canvas: c}, nil
}

// Update implements [scene.SelectStageCanvasAPI].
func (c *SelectStageCanvas) Update(active bool) error {
	return c.canvas.Update(active)
}

// Draw implements [scene.SelectStageCanvasAPI].
func (c *SelectStageCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsSelectStageCanvas implements [scene.SelectStageCanvasAPI].
func (c *SelectStageCanvas) IsSelectStageCanvas() {}

var _ scene.SelectStageCanvasAPI = (*SelectStageCanvas)(nil)
