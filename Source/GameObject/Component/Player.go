package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	parts "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/Parts"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"reflect"

	input "github.com/quasilyte/ebitengine-input"
)

type Player struct {
	ipt   *input.Handler
	tw    *TileWalker
	dir   Utils.Vec2Int
	scene *scene.PlayScene
}

func (p *Player) Update(active bool) (err error) {
	// なにもしない
	if !active {
		return nil
	}

	for _, action := range inputconfig.ActionConstList() {
		if p.ipt.ActionIsJustPressed(action) {
			switch action {
			case inputconfig.ActionMoveLeft:
				p.tw.Move(Utils.NewVecIntLeft())
				p.dir = Utils.NewVecIntLeft()
			case inputconfig.ActionMoveRight:
				p.tw.Move(Utils.NewVecIntRight())
				p.dir = Utils.NewVecIntRight()
			case inputconfig.ActionMoveUp:
				p.tw.Move(Utils.NewVecIntUp())
				p.dir = Utils.NewVecIntUp()
			case inputconfig.ActionMoveDown:
				p.tw.Move(Utils.NewVecIntDown())
				p.dir = Utils.NewVecIntDown()
			case inputconfig.ActionInteract:
				err := p.interact()
				if err != nil {
					return fmt.Errorf("p.interact() でエラーが発生しました。\n%w", err)
				}
			default: // 何もしない
			}
		}
	}

	return nil
}

func NewPlayer(scene *scene.PlayScene, obj *gameobject.GameObject) Utils.Factory[*Player] {
	return func() (*Player, error) {
		// TileWalker を取得
		tw, err := obj.GetComponent(reflect.TypeFor[*TileWalker]())
		if err != nil {
			return nil, fmt.Errorf("obj.GetComponent(reflect.TypeFor[*TileWalker]()) で失敗しました。\n%w", err)
		}

		pl := &Player{
			tw:    tw.(*TileWalker),
			ipt:   scene.Input,
			scene: scene,
		}

		return pl, nil
	}
}

func (p *Player) interact() error {
	mypos := p.tw.Position.Position
	targetpos := mypos.Add(p.dir)

	if tgt, ok := p.tw.tm.Interact[targetpos]; ok {
		err := tgt.Interact()
		if err != nil {
			return fmt.Errorf("target.Interact() でエラーが発生しました。\ntergetpos = %d\n%w", targetpos, err)
		}
		p.scene.Gamedata.AddPart(parts.RandomPartsFromAll())
		var item []string
		for _, p := range p.scene.Gamedata.Parts {
			item = append(item, p.Name+"\n")
		}
		fmt.Print(item)
	}

	return nil
}

// test
var _ gameobject.Component = &Player{}
