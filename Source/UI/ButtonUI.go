package ui

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ButtonUI struct {
	content *RichUI
	hovered bool
	pressed bool
}

func NewButtonUI(content *RichUI) (*ButtonUI, error) {
	if content == nil {
		return nil, fmt.Errorf("content が nil です")
	}
	return &ButtonUI{content: content}, nil
}

func (b *ButtonUI) Update(active bool) error {
	b.hovered = false
	b.pressed = false
	if !active {
		return nil
	}

	x, y := ebiten.CursorPosition()
	if image.Pt(x, y).In(b.content.Bounds()) {
		b.hovered = true
		b.pressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	}
	return b.content.Update(active)
}

func (b *ButtonUI) Draw(screen *ebiten.Image) {
	b.content.Draw(screen)
}

func (b *ButtonUI) IsMouseOver() bool {
	return b.hovered
}

func (b *ButtonUI) IsButtonPressed() bool {
	return b.pressed
}

func (b *ButtonUI) RichUI() *RichUI {
	return b.content
}

var _ Element = (*ButtonUI)(nil)
