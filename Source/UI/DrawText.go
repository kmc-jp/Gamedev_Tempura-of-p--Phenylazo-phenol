package ui

import (
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DrawText(pos utils.Position, color *color.Color, txt string, font *text.GoTextFace, screen *ebiten.Image) {
	textop := &text.DrawOptions{}
	textop.GeoM.Translate(pos.X, pos.Y)
	textop.ColorScale.Scale(color.R, color.G, color.B, color.A)
	textop.LineSpacing = font.Size * 1.5

	text.Draw(screen, txt, font, textop)
}
