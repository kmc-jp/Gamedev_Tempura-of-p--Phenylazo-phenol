package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// カメラの本実装
// 場所を指定したらその場所を中心として映す
type UICamera struct {
}

// Draw implements [WorldDrawer].
func (c UICamera) Draw(screen *ebiten.Image, objects []Element) {
	for _, elm := range objects {
		elm.Draw(screen)
	}
}

var _ UIDrawer = &UICamera{}

func NewUICamera() Utils.Factory[*UICamera] {
	return func() (*UICamera, error) {
		return &UICamera{}, nil
	}
}
