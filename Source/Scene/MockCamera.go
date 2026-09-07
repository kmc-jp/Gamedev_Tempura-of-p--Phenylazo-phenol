package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// 仮作成のカメラ(本実装は後から作る)
type MockCamera struct {
}

// Draw implements [Drawer].
func (m MockCamera) Draw(screen *ebiten.Image, objects []Entity, center Utils.Position) {
	panic("unimplemented")
}

func test() {
	var _ Drawer = MockCamera{}
}
