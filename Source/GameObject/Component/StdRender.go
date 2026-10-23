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
type StdRender struct {
	asp   *goaseprite.File
	img   *ebiten.Image
	layer string
}

func NewLayerableRender(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[*StdRender] {
	return func() (*StdRender, error) {
		if asp == nil {
			return nil, fmt.Errorf("asp が nil です。")
		}

		player := asp.CreatePlayer()

		if asp.Tags == nil || asp.HasTag("") {
			// タグ設定がない場合
			err := player.Play("")
			if err != nil {
				return nil, fmt.Errorf("goaseprite.File.CreatePlayer().Play(\"\") でエラーが発生しました。\n%w", err)
			}
		} else {
			// タグ設定があれば、はじめのタグを使う
			startTag := asp.Tags[0]
			err := player.Play(startTag.Name)
			if err != nil {
				return nil, fmt.Errorf("goaseprite.File.CreatePlayer().Play(\"\") でエラーが発生しました。\n%w", err)
			}
		}

		sub := img.SubImage(image.Rect(player.CurrentFrameCoords())).(*ebiten.Image)

		rndr := StdRender{
			asp:   asp,
			img:   sub,
			layer: layer,
		}
		return &rndr, nil
	}
}

func (r StdRender) Draw() (image *ebiten.Image, Layer string) {
	return r.img, r.layer
}

var _ gameobject.Render = StdRender{}
