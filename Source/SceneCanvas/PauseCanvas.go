package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// PauseCanvas は PauseScene のUI構成である。ポーズ中の案内テキストだけを持つ。
// 背景の半透明の塗りつぶしは、従来どおり PauseScene.Draw が行う。
type PauseCanvas struct {
	canvas *canvas.Canvas
}

// NewPauseCanvas は PauseCanvas を生成する。
func NewPauseCanvas(font *text.GoTextFaceSource) (*PauseCanvas, error) {
	c, err := newTextCanvas(
		utils.Position{X: 0, Y: 50},
		"Pause\nPress ESC Key to Return the Game",
		font,
		15,
	)
	if err != nil {
		return nil, fmt.Errorf("ポーズ画面のCanvas構築に失敗しました。\n%w", err)
	}
	return &PauseCanvas{canvas: c}, nil
}

// Update implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) Update(active bool) error {
	return c.canvas.Update(active)
}

// Draw implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsPauseCanvas implements [scene.PauseCanvasAPI].
func (c *PauseCanvas) IsPauseCanvas() {}

var _ scene.PauseCanvasAPI = (*PauseCanvas)(nil)
