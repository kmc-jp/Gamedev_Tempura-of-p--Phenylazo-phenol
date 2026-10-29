package ui

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameDataText struct {
	position utils.Position
	color    *color.Color
	gdata    *game.Gamedata
	text     func(game.Gamedata) string
	font     *text.GoTextFace
}

func NewGameDataText(
	position utils.Position,
	color *color.Color,
	gdata *game.Gamedata,
	text func(game.Gamedata) string,
	font *text.GoTextFace,
) (*GameDataText, error) {
	if color == nil {
		return nil, fmt.Errorf("color が nil です")
	}
	if font == nil {
		return nil, fmt.Errorf("font が nil です")
	}

	return &GameDataText{
		position: position,
		color:    color,
		gdata:    gdata,
		text:     text,
		font:     font,
	}, nil
}

func (t *GameDataText) Update(active bool) error {
	return nil
}

func (t *GameDataText) Draw(screen *ebiten.Image) {
	DrawText(t.position, t.color, t.text(*t.gdata), t.font, screen)
}

func (t *GameDataText) SetText(value func(game.Gamedata) string) {
	t.text = value
}

var _ scene.Element = (*Text)(nil)
