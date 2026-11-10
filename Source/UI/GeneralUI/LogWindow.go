package generalui

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type LogWindow struct {
	*ui.RichUI

	DisplayTime uint //メッセージログ表示時間; unit: 0.1s

	timer *tick.Timer
}

const (
	TextSize    = 30 //WindowHeight
	WindowWidth = 500
)

var (
	TextColor       = gamecolor.Black()
	BackgroundColor = color.RGBA{200, 200, 200, 200}
)

func NewLogWindow(pos Utils.Position, gdata game.Gamedata, displaytime uint) (*LogWindow, error) {
	t, err := ui.NewText(pos, TextColor, "", ui.FontFace(gdata.Font, TextSize))
	if err != nil {
		return nil, fmt.Errorf("NewLogWindow()のNewText()でエラーが発生しました\n%w\n", err)
	}
	bg := ebiten.NewImage(WindowWidth, TextSize)
	bg.Fill(BackgroundColor)
	rui, err := ui.NewRichUI(pos, image.Point{X: WindowWidth, Y: TextSize}, bg, t)
	if err != nil {
		return nil, fmt.Errorf("NewLogWindow()のNewRichUI()でエラーが発生しました\n%w\n", err)
	}
	return &LogWindow{
		RichUI:      rui,
		DisplayTime: displaytime,
		timer:       tick.NewTimer(),
	}, nil
}
