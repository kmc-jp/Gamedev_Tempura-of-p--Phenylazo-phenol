package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// カメラの本実装
// 場所を指定したらその場所を中心として映す
type Camera struct {
	center Utils.Position
}

// Draw implements [Drawer].
func (c Camera) Draw(screen *ebiten.Image, objects []Entity) {
	scSize := screen.Bounds().Size()
	for _, obj := range objects {
		image, _ := obj.Draw()
		pos := obj.Position().Offset(c.center)

		posx, posy := pos.X+float64(scSize.X)/2, pos.Y+float64(scSize.Y)/2
		sizex, sizey := image.Bounds().Dx(), image.Bounds().Dy()
		iposx, iposy := math.Floor(posx-float64(sizex)/2), math.Floor(posy-float64(sizey)/2)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(iposx, iposy)
		screen.DrawImage(image, op)
	}
}

func (c *Camera) SetCenter(center Utils.Position) {
	c.center = center
}

var _ Drawer = &Camera{}

func NewCamera() (*Camera, error) {
	return &Camera{
		center: Utils.Position{},
	}, nil
}

// test
var _ Utils.Factory[*Camera] = NewCamera
