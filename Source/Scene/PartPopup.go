package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	partPopupDuration = 10
	partPopupOffset   = 30
)

type partPopup struct {
	image    *ebiten.Image
	position Utils.Position
	timer    *tick.Timer
	expired  bool
}

func NewPartPopup(image *ebiten.Image, position Utils.Position) Entity {
	bounds := image.Bounds()
	popupImage := ebiten.NewImage((bounds.Dx()+3)/2, (bounds.Dy()+3)/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	popupImage.DrawImage(image, op)

	position.X += (rand.Float64()*2 - 1) * partPopupOffset
	position.Y += (rand.Float64()*2 - 1) * partPopupOffset

	timer := tick.NewTimer()
	_ = timer.Set(partPopupDuration)
	return &partPopup{
		image:    popupImage,
		position: position,
		timer:    timer,
	}
}

func (p *partPopup) Name() string {
	return "PartPopup"
}

func (p *partPopup) Update(active bool) error {
	p.position.Y -= 1
	if active && p.timer.Expired() {
		p.expired = true
	}
	return nil
}

func (p *partPopup) Draw() (*ebiten.Image, string) {
	return p.image, "Popup"
}

func (p *partPopup) Position() Utils.Position {
	return p.position
}

func (p *partPopup) IsExpired() bool {
	return p.expired
}

var _ Entity = (*partPopup)(nil)
