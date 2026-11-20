package scenecanvas

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	howToPlayHintTextSize = 12
	howToPlayHintMargin   = 2
	howToPlayHintMessage  = "Esc / Enter：もどる"
)

// HowToPlayCanvas は HowToPlayScene のUI構成である。遊び方の画像と、閉じ方の案内を持つ。
type HowToPlayCanvas struct {
	page *ebiten.Image
	hint *ui.Text
}

// NewHowToPlayCanvas は HowToPlayCanvas を生成する。
func NewHowToPlayCanvas(gdata *game.Gamedata, page *ebiten.Image) (*HowToPlayCanvas, error) {
	if gdata == nil || gdata.Font == nil {
		return nil, fmt.Errorf("gdata.Font が nil です")
	}
	if page == nil {
		return nil, fmt.Errorf("遊び方の画像が nil です")
	}
	hint, err := ui.NewText(utils.Position{}, color.DarkGray(), howToPlayHintMessage, ui.FontFace(gdata.Font, howToPlayHintTextSize))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() -> hint でエラーが発生しました。\n%w", err)
	}
	return &HowToPlayCanvas{page: page, hint: hint}, nil
}

// Update implements [scene.HowToPlayCanvasAPI].
func (c *HowToPlayCanvas) Update(active bool) error {
	return c.hint.Update(active)
}

// Draw implements [scene.HowToPlayCanvasAPI].
// 画像は縦横比を保ったまま画面に収まる最大の大きさで中央に描く。
func (c *HowToPlayCanvas) Draw(screen *ebiten.Image) {
	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	screen.Fill(color.White().ToRGBA())

	iw, ih := float64(c.page.Bounds().Dx()), float64(c.page.Bounds().Dy())
	scale := min(sw/iw, sh/ih)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((sw-iw*scale)/2, (sh-ih*scale)/2)
	screen.DrawImage(c.page, op)

	c.hint.SetPosition(utils.Position{
		X: sw - howToPlayHintWidth() - howToPlayHintMargin,
		Y: sh - howToPlayHintTextSize*1.5 - howToPlayHintMargin,
	})
	c.hint.Draw(screen)
}

// IsHowToPlayCanvas implements [scene.HowToPlayCanvasAPI].
func (c *HowToPlayCanvas) IsHowToPlayCanvas() {}

var _ scene.HowToPlayCanvasAPI = (*HowToPlayCanvas)(nil)

// howToPlayHintWidth は案内文のおおよその幅を返す。半角は全角の約半分として数える。
func howToPlayHintWidth() float64 {
	w := 0.0
	for _, r := range howToPlayHintMessage {
		if r < 0x80 {
			w += howToPlayHintTextSize / 2
		} else {
			w += howToPlayHintTextSize
		}
	}
	return w
}
