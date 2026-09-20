package component

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

type StdRender struct {
	asp *goaseprite.File
	img *ebiten.Image
}

func NewStdRender(asp *goaseprite.File, img *ebiten.Image) Utils.Factory[*StdRender] {
	sub := img.SubImage(image.Rect(asp.CreatePlayer().CurrentFrameCoords()))
	rndr := StdRender{
		asp: asp,
		img: sub.(*ebiten.Image),
	}
	return func() (*StdRender, error) {
		return &rndr, nil
	}
}

func NewEmptyRender() (*StdRender, error) {
	r := StdRender{
		img: ebiten.NewImage(1, 1),
	}
	return &r, nil
}

func (r StdRender) Draw() (image *ebiten.Image) {
	return r.img
}
