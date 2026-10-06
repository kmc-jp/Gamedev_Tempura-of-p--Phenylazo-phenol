package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"reflect"

	input "github.com/quasilyte/ebitengine-input"
)

type Player struct {
	ipt *input.Handler
	tw  *TileWalker
}

func (p Player) Update(active bool) (err error) {
	// なにもしない
	if !active {
		return nil
	}

	for _, action := range inputconfig.ActionConstList() {
		if p.ipt.ActionIsJustPressed(action) {
			switch action {
			case inputconfig.ActionMoveLeft:
				p.tw.Move(Utils.NewVecIntLeft())
			case inputconfig.ActionMoveRight:
				p.tw.Move(Utils.NewVecIntRight())
			case inputconfig.ActionMoveUp:
				p.tw.Move(Utils.NewVecIntUp())
			case inputconfig.ActionMoveDown:
				p.tw.Move(Utils.NewVecIntDown())
			default: // 何もしない
			}
		}
	}

	return nil
}

func NewPlayer(scene *scene.PlayScene, obj *gameobject.GameObject) Utils.Factory[*Player] {
	// TileWalker を取得
	tw, err := obj.GetComponent(reflect.TypeFor[TileWalker]())
	if err != nil {
		return func() (*Player, error) {
			return nil, fmt.Errorf("obj.GetComponent(reflect.TypeFor[TileWalker]()) で失敗しました。\n%w", err)
		}
	}

	pl := &Player{
		tw:  tw.(*TileWalker),
		ipt: scene.Input,
	}

	return func() (*Player, error) {
		return pl, nil
	}
}

// test
var _ gameobject.Component = Player{}
