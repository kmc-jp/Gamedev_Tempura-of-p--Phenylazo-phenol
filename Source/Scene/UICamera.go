package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// UI用カメラ
type UICamera struct {
}

// Draw implements [WorldDrawer].
func (c UICamera) Draw(screen *ebiten.Image, elements []Element) {
	for _, elm := range elements {
		elm.Draw(screen)
	}
}

var _ UIDrawer = &UICamera{}

func NewUICamera() Utils.Factory[*UICamera] {
	return func() (*UICamera, error) {
		return &UICamera{}, nil
	}
}
