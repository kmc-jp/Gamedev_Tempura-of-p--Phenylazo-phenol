package camera

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

type BackGroundCamera struct {
	asp        *goaseprite.File
	wholeImg   *ebiten.Image
	drawOption *ebiten.DrawImageOptions
	player     *goaseprite.Player
}

func NewBackGroundCamera(asp *goaseprite.File, wholeImg *ebiten.Image, firstTag string) Utils.Factory[*BackGroundCamera] {
	return func() (*BackGroundCamera, error) {
		player := asp.CreatePlayer()

		// タグ設定
		err := player.Play(firstTag)
		if err != nil {
			return nil, fmt.Errorf("goaseprite.File.CreatePlayer().Play(firstTag) でエラーが発生しました。\nfirstTag = %s\n%w", firstTag, err)
		}

		bgc := BackGroundCamera{
			asp:        asp,
			wholeImg:   wholeImg,
			drawOption: &ebiten.DrawImageOptions{},
			player:     player,
		}
		return &bgc, nil
	}
}

func (bgc BackGroundCamera) Draw(screen *ebiten.Image, gdata game.Gamedata) {
	bgc.drawOption.GeoM.Reset()
	bgc.drawOption.GeoM.Scale(5, 5)

	img := bgc.wholeImg.SubImage(image.Rect(bgc.player.CurrentFrameCoords())).(*ebiten.Image)

	screen.DrawImage(img, bgc.drawOption)
}

func (bgc *BackGroundCamera) Update() error {
	bgc.FrameUpdate()

	return nil
}

func (bgc *BackGroundCamera) SetTag(tag string) error {
	err := bgc.player.Play(tag)
	if err != nil {
		return fmt.Errorf("goaseprite.File.CreatePlayer().Play(firstTag) でエラーが発生しました。\nfirstTag = %s\n%w", tag, err)
	} else {
		return nil
	}
}

func (bgc *BackGroundCamera) FrameUpdate() {
	const spf float32 = 1 / 60.
	bgc.player.Update(spf)
}

// test
var _ scene.BackGroundDrawer = &BackGroundCamera{}
