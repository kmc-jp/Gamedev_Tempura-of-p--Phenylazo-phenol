package commonui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// LogWindow は指定したテキストを一定時間だけ表示するUIである。
// ゲーム内メッセージやメッセージウィンドウの状態は管理しない。
type LogWindow struct {
	*ui.Text

	DisplayTime uint16 // テキスト表示時間; unit: 0.1s
	displaying  bool
	timer       *tick.Timer
	font        *text.GoTextFace
	value       string
	background  *gamecolor.Color
	padding     image.Point
}

// NewLogWindow は指定した表示位置・フォント・表示時間でLogWindowを生成する。
func NewLogWindow(pos utils.Position, font *text.GoTextFaceSource, textSize float64, textColor *gamecolor.Color, displayTime uint16) (*LogWindow, error) {
	if font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	if displayTime == 0 {
		return nil, fmt.Errorf("displayTime は 1 以上である必要があります")
	}

	t, err := ui.NewText(pos, textColor, "", ui.FontFace(font, textSize))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました: %w", err)
	}

	return &LogWindow{
		Text:        t,
		DisplayTime: displayTime,
		timer:       tick.NewTimer(),
		font:        ui.FontFace(font, textSize),
	}, nil
}

// SetText はテキストを設定し、表示時間の計測を開始する。
func (lw *LogWindow) SetText(value string) {
	lw.value = value
	lw.Text.SetText(value)
	lw.timer.Set(0)
	lw.displaying = value != ""
	if lw.displaying {
		lw.timer.Set(lw.DisplayTime)
	}
}

// ClearText はテキストを消去する。
func (lw *LogWindow) ClearText() {
	lw.value = ""
	lw.Text.SetText("")
	lw.timer.Set(0)
	lw.displaying = false
}

// Update は表示時間が経過したテキストを消去する。
func (lw *LogWindow) Update(active bool) error {
	if !active || !lw.displaying {
		return nil
	}
	if lw.timer.Expired() {
		lw.ClearText()
	}
	return nil
}

func (lw *LogWindow) Draw(screen *ebiten.Image) {
	if lw.displaying && lw.background != nil && lw.value != "" {
		width, height := text.Measure(lw.value, lw.font, lw.font.Size*1.5)
		bounds := image.Rect(
			int(lw.Position().X)-lw.padding.X,
			int(lw.Position().Y)-lw.padding.Y,
			int(lw.Position().X+width)+lw.padding.X,
			int(lw.Position().Y+height)+lw.padding.Y,
		)
		uiutils.FillRect(screen, bounds, lw.background.ToRGBA())
	}
	lw.Text.Draw(screen)
}

// SetBackground はLogWindowの背景色と余白を設定する。
func (lw *LogWindow) SetBackground(background *gamecolor.Color, padding image.Point) {
	lw.background = background
	lw.padding = padding
}

// SetCenter は現在のテキストを画面中央に配置する。
func (lw *LogWindow) SetCenter(centerX, centerY float64) {
	width, height := 0.0, 0.0
	if lw.value != "" {
		width, height = text.Measure(lw.value, lw.font, lw.font.Size*1.5)
	}
	lw.SetPosition(utils.Position{
		X: centerX - width/2,
		Y: centerY - height/2,
	})
}

var _ canvas.CanvasItem = (*LogWindow)(nil)
