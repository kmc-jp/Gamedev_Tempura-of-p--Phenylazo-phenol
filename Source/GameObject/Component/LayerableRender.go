package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

// 静止画像を映すRender
type LayerableRender struct {
	asp   *goaseprite.File
	img   *ebiten.Image
	layer string
}

func NewLayerableRender(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[*LayerableRender] {
	player := asp.CreatePlayer()

	startTag := asp.Tags[0]

	err := player.Play(startTag.Name)
	if err != nil {
		return func() (*LayerableRender, error) {
			return nil, fmt.Errorf("goaseprite.File.CreatePlayer().Play(\"\") でエラーが発生しました。\n%w", err)
		}
	}

	sub := img.SubImage(image.Rect(player.CurrentFrameCoords())).(*ebiten.Image)

	rndr := LayerableRender{
		asp:   asp,
		img:   sub,
		layer: layer,
	}
	return func() (*LayerableRender, error) {
		return &rndr, nil
	}
}

func (r LayerableRender) Draw() (image *ebiten.Image, Layer string) {
	return r.img, r.layer
}

var _ gameobject.Render = LayerableRender{}
