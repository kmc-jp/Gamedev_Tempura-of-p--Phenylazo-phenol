package gameobject

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type HogeObject struct {
	GameObjectUtils
	tf   Transform
	rd   Render
	cmps []Component
}

// 仮実装
func (h HogeObject) Draw(screen *ebiten.Image) {
	panic("unimplemented")
}

// 仮実装
func (h HogeObject) Update(active bool) (err error) {
	panic("unimplemented")
}

var _ scene.GameObject = HogeObject{}

func NewHogeObject() (*HogeObject, error) {
	t, err := NewTransform()
	if err != nil {
		return nil, fmt.Errorf("NewTransform() でエラーが発生しました。\n %w", err)
	}
	r, err := NewRender()
	if err != nil {
		return nil, fmt.Errorf("NewRender() でエラーが発生しました。\n %w", err)
	}
	h := HogeObject{
		tf: *t,
		rd: *r,
	}
	h.GameObjectUtils.gameObject = h
	return &h, nil
}

func (h HogeObject) transform() Transform    { return h.tf }
func (h HogeObject) render() Render          { return h.rd }
func (h HogeObject) components() []Component { return h.cmps }
func (h HogeObject) Name() string            { return "Hoge" }
