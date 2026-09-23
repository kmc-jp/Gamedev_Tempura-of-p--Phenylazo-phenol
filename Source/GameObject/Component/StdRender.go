package component

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

// 静止画像を映すRender
type StdRender struct {
	asp *goaseprite.File
	img *ebiten.Image
}

func NewStdRender(asp *goaseprite.File, img *ebiten.Image) Utils.Factory[*StdRender] {
	player := asp.CreatePlayer()

	startTag := asp.Tags[0]

	err := player.Play(startTag.Name)
	if err != nil {
		return func() (*StdRender, error) {
			return nil, fmt.Errorf("goaseprite.File.CreatePlayer().Play(\"\") でエラーが発生しました。\n%w", err)
		}
	}

	sub := img.SubImage(image.Rect(player.CurrentFrameCoords())).(*ebiten.Image)

	rndr := StdRender{
		asp: asp,
		img: sub,
	}
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
