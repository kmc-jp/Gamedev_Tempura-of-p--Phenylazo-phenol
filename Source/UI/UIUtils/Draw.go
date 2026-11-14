// Package uiutils は、UIの描画で共通に使う図形描画のヘルパーを提供する。
//
// 描画のたびに ebiten.NewImage を呼ばないよう、1×1 の白画像を1枚だけ作って使い回し、
// 拡大と着色で矩形を描く。
package uiutils

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// SelectionBorderColor は選択中のUI（ボタン、タブ、パーツカード）を示す枠の色である。
var SelectionBorderColor = color.RGBA{R: 255, G: 165, B: 0, A: 255}

// pixel は矩形描画に使い回す 1×1 の白画像。初回の描画時に作る。
var pixel *ebiten.Image

func whitePixel() *ebiten.Image {
	if pixel == nil {
		pixel = ebiten.NewImage(1, 1)
		pixel.Fill(color.White)
	}
	return pixel
}

// FillRect は rect を c で塗りつぶす。rect が空なら何もしない。
func FillRect(dst *ebiten.Image, rect image.Rectangle, c color.Color) {
	if rect.Empty() {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(rect.Dx()), float64(rect.Dy()))
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	op.ColorScale.ScaleWithColor(c)
	dst.DrawImage(whitePixel(), op)
}

// DrawBorder は rect の内側に太さ width の枠を c で描く。
// rect の幅か高さが width*2 より小さいときは何もしない。
func DrawBorder(dst *ebiten.Image, rect image.Rectangle, width int, c color.Color) {
	w, h := rect.Dx(), rect.Dy()
	if width <= 0 || w < width*2 || h < width*2 {
		return
	}
	minX, minY, maxX, maxY := rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y
	FillRect(dst, image.Rect(minX, minY, maxX, minY+width), c) // 上
	FillRect(dst, image.Rect(minX, maxY-width, maxX, maxY), c) // 下
	FillRect(dst, image.Rect(minX, minY, minX+width, maxY), c) // 左
	FillRect(dst, image.Rect(maxX-width, minY, maxX, maxY), c) // 右
}
