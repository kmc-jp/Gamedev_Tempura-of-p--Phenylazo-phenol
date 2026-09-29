package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// 仮作成のカメラ(本実装作成済み)
// (0,0) が左上で、ワールド座標 = 描画座標
type MockCamera struct {
}

// Draw implements [Drawer].
func (m MockCamera) Draw(screen *ebiten.Image, objects []Entity) {
	for _, obj := range objects {
		image, _ := obj.Draw()
		pos := obj.Position()

		posx, posy := pos.X, pos.Y
		sizex, sizey := image.Bounds().Dx(), image.Bounds().Dy()
		iposx, iposy := math.Floor(posx-float64(sizex)/2), math.Floor(posy-float64(sizey)/2)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(iposx, iposy)
		screen.DrawImage(image, op)
	}
}

func (m MockCamera) SetCenter(Utils.Position) {
	// 何もしない
}

func test() {
	var _ Drawer = MockCamera{}
}

func NewMockCamera() (*MockCamera, error) {
	return &MockCamera{}, nil
}

// test
var _ Utils.Factory[*MockCamera] = NewMockCamera
