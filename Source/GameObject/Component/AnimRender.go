package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

// アニメーション画像を映すRender
type AnimRender struct {
	asp      *goaseprite.File
	wholeImg *ebiten.Image
	layer    string
	player   *goaseprite.Player
}

// アニメーションが可能な Render をつくる
// 使いまわし可能
func NewAnimRender(asp *goaseprite.File, img *ebiten.Image, layer string, firstTag string) Utils.Factory[*AnimRender] {
	return func() (*AnimRender, error) {
		if asp == nil {
			return nil, fmt.Errorf("asp が nil です。")
		}

		player := asp.CreatePlayer()

		// タグ設定
		err := player.Play(firstTag)
		if err != nil {
			return nil, fmt.Errorf("goaseprite.File.CreatePlayer().Play(firstTag) でエラーが発生しました。\nfirstTag = %s\n%w", firstTag, err)
		}

		rndr := AnimRender{
			asp:      asp,
			wholeImg: img,
			layer:    layer,
			player:   player,
		}

		return &rndr, nil
	}
}

func (r *AnimRender) FrameUpdate() {
	const spf float32 = 1 / 60.
	r.player.Update(spf)
}

func (r *AnimRender) SetTag(tag string) error {
	err := r.player.Play(tag)
	if err != nil {
		return fmt.Errorf("goaseprite.File.CreatePlayer().Play(firstTag) でエラーが発生しました。\nfirstTag = %s\n%w", tag, err)
	} else {
		return nil
	}
}

func (r AnimRender) Draw() (*ebiten.Image, string) {
	img := r.wholeImg.SubImage(image.Rect(r.player.CurrentFrameCoords())).(*ebiten.Image)

	return img, r.layer
}

func (r *AnimRender) Update() error {
	r.FrameUpdate()

	return nil
}

var _ gameobject.Render = &AnimRender{}
