// Package titleui は、TitleScene で使うUI部品を提供する。
package titleui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// CursorStyle はカーソルの見た目。人が調整する値はここに集約する。
// 背景画像と同じ解像度のドット絵として描くため、長さの単位は背景画像の1ピクセル
// （描画時に scale 倍の正方形になる）。三角形も揺れも1ピクセル単位で動く。
type CursorStyle struct {
	Color         color.Color
	Width, Height int     // 三角形の幅・高さ（ピクセル）。Height は奇数にすると上下対称になる
	SwayAmplitude float64 // 左右の揺れ幅（ピクセル）。描画時に整数へ丸める
	SwayPeriod    int     // 揺れ1周期のフレーム数
}

// DefaultCursorStyle は既定の見た目を返す。
func DefaultCursorStyle() CursorStyle {
	return CursorStyle{
		Color:         uiutils.SelectionBorderColor,
		Width:         3,
		Height:        5,
		SwayAmplitude: 2,
		SwayPeriod:    60,
	}
}

// MenuCursor は選択中の項目を指す右向きの三角形カーソルである。
// 入力は扱わず、選択中の項目は SetSelected で外から渡す。
type MenuCursor struct {
	anchors  []utils.Position
	scale    float64
	style    CursorStyle
	selected int
	frame    int
}

// NewMenuCursor は MenuCursor を生成する。
// anchors は各項目で三角形の先端（右端・縦中央）が指す位置（背景画像のピクセル座標）。
func NewMenuCursor(anchors []utils.Position, scale float64, style CursorStyle) (*MenuCursor, error) {
	if len(anchors) == 0 {
		return nil, fmt.Errorf("anchors が空です")
	}
	if scale <= 0 {
		return nil, fmt.Errorf("scale は正である必要があります: %v", scale)
	}
	if style.Width <= 0 || style.Height <= 0 {
		return nil, fmt.Errorf("Width・Height は正である必要があります: %d, %d", style.Width, style.Height)
	}
	if style.SwayPeriod <= 0 {
		return nil, fmt.Errorf("SwayPeriod は正である必要があります: %d", style.SwayPeriod)
	}
	return &MenuCursor{anchors: anchors, scale: scale, style: style}, nil
}

// SetSelected は選択中の項目の番号を設定する。範囲外なら何もしない。
func (c *MenuCursor) SetSelected(i int) {
	if i < 0 || i >= len(c.anchors) {
		return
	}
	c.selected = i
}

// Update は揺れのためにフレームを進める。
func (c *MenuCursor) Update(active bool) error {
	if active {
		c.frame++
	}
	return nil
}

// Draw は選択中の項目の左に三角形を描く。画像1ピクセルを scale×scale の正方形として塗る。
func (c *MenuCursor) Draw(screen *ebiten.Image) {
	a := c.anchors[c.selected]
	sway := math.Round(c.style.SwayAmplitude * math.Sin(2*math.Pi*float64(c.frame)/float64(c.style.SwayPeriod)))
	tipX := int(math.Round(a.X + sway)) // 先端の右隣の列
	centerY := int(math.Round(a.Y))     // 中央の行

	// 行ごとに左端をそろえた長さで塗る。中央の行が最長で、上下へ離れるほど短くなる。
	half := (c.style.Height - 1) / 2
	s := int(math.Round(c.scale))
	for row := 0; row < c.style.Height; row++ {
		d := row - half
		if d < 0 {
			d = -d
		}
		length := int(math.Round(float64(c.style.Width) * (1 - float64(d)/float64(half+1))))
		if length < 1 {
			length = 1
		}
		x0 := tipX - c.style.Width
		y := centerY - half + row
		rect := image.Rect(x0*s, y*s, (x0+length)*s, (y+1)*s)
		uiutils.FillRect(screen, rect, c.style.Color)
	}
}

var _ canvas.CanvasItem = (*MenuCursor)(nil)
