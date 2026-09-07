package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// 仮作成のカメラ(本実装は後から作る)
type MockCamera struct {
}

// Draw implements [Drawer].
func (m MockCamera) Draw(screen *ebiten.Image, objects []Entity, center Utils.Position) {
	for _, obj := range objects {
		image := obj.Draw()
		pos := obj.Position()

		posx, posy := pos.X, pos.Y
		sizex, sizey := image.Bounds().Dx(), image.Bounds().Dy()
		iposx, iposy := math.Floor(posx-float64(sizex)/2), math.Floor(posy-float64(sizey))

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(iposx, iposy)
		screen.DrawImage(image, op)
	}
}

func test() {
	var _ Drawer = MockCamera{}
}
