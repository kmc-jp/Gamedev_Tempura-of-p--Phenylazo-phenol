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

	DisplayTime uint16 //メッセージログ表示時間; unit: 0.1s
	displaying  bool

	allowInterrupt bool

	timer *tick.Timer

	gdata  game.Gamedata
	msg_sn uint
}

const (
	TextSize    = 30 //WindowHeight
	WindowWidth = 500
)

var (
	TextColor       = gamecolor.Black()
	BackgroundColor = color.RGBA{200, 200, 200, 200}
)

func NewLogWindow(pos Utils.Position, gdata game.Gamedata, displaytime uint16) (*LogWindow, error) {
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
		RichUI:         rui,
		DisplayTime:    displaytime,
		displaying:     false,
		allowInterrupt: true,
		timer:          tick.NewTimer(),
		gdata:          gdata,
		msg_sn:         0,
	}, nil
}

func (lw *LogWindow) SetMessage(text string) {
	lw.RichUI.SetText(text)
}

func (lw *LogWindow) ClearMessage() {
	lw.RichUI.SetText("")
}

func (lw *LogWindow) Draw(screen *ebiten.Image) {
	lw.RichUI.Draw(screen)
}

func (lw *LogWindow) Update(active bool) error {
	if !active {
		return nil
	}
	if lw.gdata.Message == "" {
		return nil
	}
	if !lw.displaying {
		lw.displaying = true
		lw.timer.Set(lw.DisplayTime)
		lw.SetMessage(lw.gdata.Message)
		lw.msg_sn = lw.gdata.Message_sn
	} else {
		if lw.msg_sn != lw.gdata.Message_sn && lw.allowInterrupt {
			lw.timer.Set(0)
			lw.timer.Set(lw.DisplayTime)
			lw.SetMessage(lw.gdata.Message)
			lw.msg_sn = lw.gdata.Message_sn
		} else if lw.msg_sn != lw.gdata.Message_sn && !lw.allowInterrupt {
			return fmt.Errorf("*LogWindow.Update()でメッセージウィンドウのテキスト割り込みが発生しました")
		} else {
			//no process
		}
	}

	if lw.timer.Expired() {
		lw.timer.Set(0)
		lw.ClearMessage()
		lw.displaying = false
	}
	return nil
}
