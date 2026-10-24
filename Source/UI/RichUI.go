package ui

import (
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type RichUI struct {
	position utils.Position
	size     image.Point
	image    *ebiten.Image
	text     *Text
}

func NewRichUI(
	position utils.Position,
	size image.Point,
	background *ebiten.Image,
	text *Text,
) (*RichUI, error) {
	if size.X < 0 || size.Y < 0 {
		return nil, fmt.Errorf("size は 0 以上である必要があります: %v", size)
	}
	if background == nil && text == nil {
		return nil, fmt.Errorf("画像またはテキストのどちらかが必要です")
	}
	if background != nil && (size.X == 0 || size.Y == 0) {
		size = background.Bounds().Size()
	}

	return &RichUI{
		position: position,
		size:     size,
		image:    background,
		text:     text,
	}, nil
}

func (r *RichUI) Update(active bool) error {
	if r.text == nil {
		return nil
	}
	return r.text.Update(active)
}

func (r *RichUI) Draw(screen *ebiten.Image) {
	if r.image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(r.position.X, r.position.Y)
		screen.DrawImage(r.image, op)
	}
	if r.text != nil {
		r.text.Draw(screen)
	}
}

func (r *RichUI) Position() utils.Position {
	return r.position
}

func (r *RichUI) Size() image.Point {
	return r.size
}

func (r *RichUI) Bounds() image.Rectangle {
	x := int(r.position.X)
	y := int(r.position.Y)
	return image.Rect(x, y, x+r.size.X, y+r.size.Y)
}

func (r *RichUI) SetPosition(position utils.Position) {
	r.position = position
}

var _ Element = (*RichUI)(nil)
