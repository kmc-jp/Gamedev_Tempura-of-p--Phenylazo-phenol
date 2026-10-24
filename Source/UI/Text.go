package ui

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Text struct {
	position utils.Position
	color    *color.Color
	text     string
	font     *text.GoTextFace
}

func NewText(
	position utils.Position,
	color *color.Color,
	value string,
	font *text.GoTextFace,
) (*Text, error) {
	if color == nil {
		return nil, fmt.Errorf("color が nil です")
	}
	if font == nil {
		return nil, fmt.Errorf("font が nil です")
	}

	return &Text{
		position: position,
		color:    color,
		text:     value,
		font:     font,
	}, nil
}

func (t *Text) Update(active bool) error {
	return nil
}

func (t *Text) Draw(screen *ebiten.Image) {
	DrawText(t.position, t.color, t.text, t.font, screen)
}

func (t *Text) SetText(value string) {
	t.text = value
}

var _ scene.Element = (*Text)(nil)
