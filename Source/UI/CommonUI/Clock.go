package commonui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"
	stdcolor "image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Clock はGamedataのWorldTimerを時計として表示する専用UIである。
type Clock struct {
	position utils.Position
	color    *gamecolor.Color
	gdata    *game.Gamedata
	font     *text.GoTextFace
}

// NewClock は世界時刻を表示するClockを生成する。
func NewClock(
	position utils.Position,
	color *gamecolor.Color,
	gdata *game.Gamedata,
	font *text.GoTextFace,
) (*Clock, error) {
	if color == nil {
		return nil, fmt.Errorf("color が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}
	if gdata.WorldTimer == nil {
		return nil, fmt.Errorf("gdata.WorldTimer が nil です")
	}
	if font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	return &Clock{
		position: position,
		color:    color,
		gdata:    gdata,
		font:     font,
	}, nil
}

// Update は時計表示に必要な内部状態を持たないため何もしない。
func (c *Clock) Update(active bool) error {
	return nil
}

// Draw はGamedataのWorldTimerから取得した時刻を描画する。
func (c *Clock) Draw(screen *ebiten.Image) {
	value := c.gdata.WorldTimer.GetTimeText()
	textWidth, textHeight := text.Measure(value, c.font, c.font.Size*1.5)
	padding := 5.0
	left := int(math.Floor(c.position.X - padding))
	top := int(math.Floor(c.position.Y - padding))
	right := int(math.Ceil(c.position.X + textWidth + padding))
	bottom := int(math.Ceil(c.position.Y + math.Max(textHeight, c.font.Size) + padding))
	uiutils.FillRect(screen, image.Rect(left, top, right, bottom), stdcolor.White)
	ui.DrawText(c.position, c.color, value, c.font, screen)
}

// Position は時計の左上位置を返す。
func (c *Clock) Position() utils.Position {
	return c.position
}

// SetPosition は時計の位置を設定する。
func (c *Clock) SetPosition(position utils.Position) {
	c.position = position
}

var _ canvas.CanvasItem = (*Clock)(nil)
var _ canvas.PositionItem = (*Clock)(nil)
