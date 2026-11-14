package ui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	"fmt"

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
	uiutils.DrawBorder(screen, b.content.Bounds(), buttonBorderWidth, uiutils.SelectionBorderColor)
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
