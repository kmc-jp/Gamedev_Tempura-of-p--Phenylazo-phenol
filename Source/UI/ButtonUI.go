package ui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	"fmt"
	stdcolor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const buttonBorderWidth = 3

// ButtonUI は RichUI をラップし、選択状態の表示だけを担当する。
// 入力の判定は Scene・Component 側が input.Handler で行い、結果を SetSelected で渡す。
type ButtonUI struct {
	content  *RichUI
	selected bool
	// 入力はScene・Component側（input.Handler）へ集約するため無効化。Scene側の入力移行は完了したが、再利用の可能性があるため残している。
	// hovered bool
	// pressed bool
}

func NewButtonUI(content *RichUI) (*ButtonUI, error) {
	if content == nil {
		return nil, fmt.Errorf("content が nil です")
	}
	return &ButtonUI{content: content}, nil
}

func (b *ButtonUI) Update(active bool) error {
	// 入力はScene・Component側（input.Handler）へ集約するため無効化。Scene側の入力移行は完了したが、再利用の可能性があるため残している。
	// b.hovered = false
	// b.pressed = false
	// if !active {
	// 	return nil
	// }
	//
	// x, y := ebiten.CursorPosition()
	// if image.Pt(x, y).In(b.content.Bounds()) {
	// 	b.hovered = true
	// 	b.pressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	// }
	return b.content.Update(active)
}

func (b *ButtonUI) Draw(screen *ebiten.Image) {
	b.content.Draw(screen)
	if !b.selected {
		return
	}
	bounds := b.content.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	bw := buttonBorderWidth
	if w < bw*2 || h < bw*2 {
		return
	}
	borderColor := stdcolor.RGBA{R: 255, G: 165, B: 0, A: 255}
	for _, edge := range []struct{ x, y, w, h int }{
		{bounds.Min.X, bounds.Min.Y, w, bw},
		{bounds.Min.X, bounds.Max.Y - bw, w, bw},
		{bounds.Min.X, bounds.Min.Y, bw, h},
		{bounds.Max.X - bw, bounds.Min.Y, bw, h},
	} {
		line := ebiten.NewImage(edge.w, edge.h)
		line.Fill(borderColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(edge.x), float64(edge.y))
		screen.DrawImage(line, op)
	}
}

func (b *ButtonUI) SetSelected(selected bool) {
	b.selected = selected
}

func (b *ButtonUI) Selected() bool {
	return b.selected
}

// 入力はScene・Component側（input.Handler）へ集約するため無効化。Scene側の入力移行は完了したが、再利用の可能性があるため残している。
// func (b *ButtonUI) IsMouseOver() bool {
// 	return b.hovered
// }
//
// func (b *ButtonUI) IsButtonPressed() bool {
// 	return b.pressed
// }

func (b *ButtonUI) RichUI() *RichUI {
	return b.content
}

var _ canvas.CanvasItem = (*ButtonUI)(nil)
