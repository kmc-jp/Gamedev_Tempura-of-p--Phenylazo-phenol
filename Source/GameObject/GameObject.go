package gameobject

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"fmt"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject struct {
	name string
	tf   Transform
	rd   Render
	cmps []Component
}

var _ scene.GameObject = GameObject{}

func (g GameObject) GetComponent(target reflect.Type) (component Component, err error) {
	for _, c := range g.components() {
		if target == reflect.TypeOf(c) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%s に %s が存在しません。", g.Name(), target.Name())
}

func (g GameObject) Update(active bool) (err error) {
	for _, c := range g.components() {
		if err := c.Update(g, active); err != nil {
			return fmt.Errorf("コンポーネント %s でエラーが発生しました。\n%w", reflect.TypeOf(c), err)
		}
	}
	return nil
}

func (g GameObject) Draw(screen *ebiten.Image) {
	g.render().Draw(g, screen)
}

func NewHogeObject() (*GameObject, error) {
	t, err := NewTransform()
	if err != nil {
		return nil, fmt.Errorf("NewTransform() でエラーが発生しました。\n %w", err)
	}
	r, err := NewRender()
	if err != nil {
		return nil, fmt.Errorf("NewRender() でエラーが発生しました。\n %w", err)
	}
	h := GameObject{
		tf: *t,
		rd: *r,
	}
	return &h, nil
}

func (g GameObject) transform() Transform    { return g.tf }
func (g GameObject) render() Render          { return g.rd }
func (g GameObject) components() []Component { return g.cmps }
func (g GameObject) Name() string            { return g.name }
