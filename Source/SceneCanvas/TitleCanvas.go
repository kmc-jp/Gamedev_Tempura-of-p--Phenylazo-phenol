package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TitleCanvas は TitleScene のUI構成である。タイトルの案内テキストだけを持つ。
type TitleCanvas struct {
	canvas *canvas.Canvas
}

// NewTitleCanvas は TitleCanvas を生成する。
func NewTitleCanvas(font *text.GoTextFaceSource) (*TitleCanvas, error) {
	c, err := newTextCanvas(utils.Position{}, "Title\nPress Space", font, 15)
	if err != nil {
		return nil, fmt.Errorf("タイトル画面のCanvas構築に失敗しました。\n%w", err)
	}
	return &TitleCanvas{canvas: c}, nil
}

// Update implements [scene.TitleCanvasAPI].
func (c *TitleCanvas) Update(active bool) error {
	return c.canvas.Update(active)
}

// Draw implements [scene.TitleCanvasAPI].
func (c *TitleCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsTitleCanvas implements [scene.TitleCanvasAPI].
func (c *TitleCanvas) IsTitleCanvas() {}

var _ scene.TitleCanvasAPI = (*TitleCanvas)(nil)
