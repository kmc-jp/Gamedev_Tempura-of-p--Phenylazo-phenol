package component

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

type StdRender struct {
	asp *goaseprite.File
	img *ebiten.Image
}

func NewStdRender(asp *goaseprite.File, img *ebiten.Image) Utils.Factory[*StdRender] {
	println("NewStdRender()")
	sub := img.SubImage(image.Rect(asp.CreatePlayer().CurrentFrameCoords()))
	rndr := StdRender{
		asp: asp,
		img: sub.(*ebiten.Image),
	}
	fmt.Printf("image : %s\n", img.Bounds().Max)
	return func() (*StdRender, error) {
		return &rndr, nil
	}
}

func NewEmptyRender() (*StdRender, error) {
	EmptyImage := ebiten.NewImage(1, 1)
	r := StdRender{
		img: EmptyImage,
	}
	return &r, nil
}

func (r StdRender) Draw() (image *ebiten.Image) {
	return r.img
}
