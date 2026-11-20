package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	registerBalloonText      = "お会計"
	registerBalloonFontSize  = 16
	registerBalloonPadX      = 8
	registerBalloonPadY      = 5
	registerBalloonTailH     = 7
	registerBalloonTailHalfW = 6
	registerBalloonBorder    = 2
	// 浮遊アニメーションの周期(フレーム)と振幅(px)
	registerBalloonPeriod    = 90
	registerBalloonAmplitude = 4
)

// registerBalloon はレジ台の上に浮遊する「お会計」吹き出し。
type registerBalloon struct {
	image *ebiten.Image
	base  Utils.Position
	frame int
	DIO   *ebiten.DrawImageOptions
}

// NewRegisterBalloon は base を中心とする吹き出しを生成する。
func NewRegisterBalloon(fontSrc *text.GoTextFaceSource, base Utils.Position) (Entity, error) {
	if fontSrc == nil {
		return nil, fmt.Errorf("吹き出しのフォントが nil です")
	}
	face := &text.GoTextFace{Source: fontSrc, Size: registerBalloonFontSize}
	tw, th := text.Measure(registerBalloonText, face, 0)

	bodyW := float32(math.Ceil(tw)) + registerBalloonPadX*2
	bodyH := float32(math.Ceil(th)) + registerBalloonPadY*2
	w := int(bodyW)
	h := int(bodyH) + registerBalloonTailH

	img := ebiten.NewImage(w, h)
	border := color.RGBA{R: 60, G: 40, B: 30, A: 255}
	fill := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	b := float32(registerBalloonBorder)

	// 本体（枠線 → 塗り）
	vector.FillRect(img, 0, 0, bodyW, bodyH, border, false)
	vector.FillRect(img, b, b, bodyW-b*2, bodyH-b*2, fill, false)

	// しっぽ（下向き三角。枠線 → 塗り）
	cx := bodyW / 2
	drawTail := func(inset float32, c color.RGBA) {
		var p vector.Path
		p.MoveTo(cx-registerBalloonTailHalfW+inset, bodyH-b)
		p.LineTo(cx+registerBalloonTailHalfW-inset, bodyH-b)
		p.LineTo(cx, float32(h)-inset)
		p.Close()
		dop := &vector.DrawPathOptions{AntiAlias: true}
		dop.ColorScale.ScaleWithColor(c)
		vector.FillPath(img, &p, &vector.FillOptions{}, dop)
	}
	drawTail(0, border)
	drawTail(float32(registerBalloonBorder), fill)
	// しっぽの付け根にある本体下辺の枠線を消してつなげる
	vector.FillRect(img, cx-registerBalloonTailHalfW+b*2, bodyH-b, registerBalloonTailHalfW*2-b*4, b, fill, false)

	// 文字
	op := &text.DrawOptions{}
	op.GeoM.Translate(registerBalloonPadX, registerBalloonPadY)
	op.ColorScale.ScaleWithColor(border)
	text.Draw(img, registerBalloonText, face, op)

	return &registerBalloon{
		image: img,
		base:  base,
		DIO:   &ebiten.DrawImageOptions{},
	}, nil
}

// 画像の高さ（呼び出し側が配置の計算に使う）
func RegisterBalloonHeight(fontSrc *text.GoTextFaceSource) float64 {
	face := &text.GoTextFace{Source: fontSrc, Size: registerBalloonFontSize}
	_, th := text.Measure(registerBalloonText, face, 0)
	return math.Ceil(th) + registerBalloonPadY*2 + registerBalloonTailH
}

func (b *registerBalloon) Name() string { return "RegisterBalloon" }

func (b *registerBalloon) Update(active bool) error {
	b.frame++
	return nil
}

func (b *registerBalloon) Draw() (*ebiten.Image, *ebiten.DrawImageOptions, string) {
	return b.image, b.DIO, "Popup"
}

func (b *registerBalloon) Position() Utils.Position {
	p := b.base
	p.Y += math.Sin(float64(b.frame)*2*math.Pi/registerBalloonPeriod) * registerBalloonAmplitude
	return p
}

var _ Entity = (*registerBalloon)(nil)
