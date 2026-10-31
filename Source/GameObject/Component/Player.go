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
	anim  *AnimRender
}

func (p *Player) Update(active bool) (err error) {
	// なにもしない
	if !active {
		return nil
	}

	if p.ipt.ActionIsPressed(inputconfig.ActionMoveLeft) {
		p.tw.MoveStart(Utils.NewVecIntLeft())
		p.dir = Utils.NewVecIntLeft()
	}
	if p.ipt.ActionIsPressed(inputconfig.ActionMoveRight) {
		p.tw.MoveStart(Utils.NewVecIntRight())
		p.dir = Utils.NewVecIntRight()
	}
	if p.ipt.ActionIsPressed(inputconfig.ActionMoveUp) {
		p.tw.MoveStart(Utils.NewVecIntUp())
		p.dir = Utils.NewVecIntUp()
	}
	if p.ipt.ActionIsPressed(inputconfig.ActionMoveDown) {
		p.tw.MoveStart(Utils.NewVecIntDown())
		p.dir = Utils.NewVecIntDown()
	}
	if p.ipt.ActionIsJustPressed(inputconfig.ActionInteract) {
		err := p.interact()
		if err != nil {
			return fmt.Errorf("p.interact() でエラーが発生しました。\n%w", err)
		}
	}

	return nil
}

func NewPlayer(scene *scene.PlayScene, obj *gameobject.GameObject, anim *AnimRender) Utils.Factory[*Player] {
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
			anim:  anim,
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
		msg := p.scene.Gamedata.AddPartToFoundList(parts.RandomPartsFromAll())
		fmt.Print(msg, "\n")
		var item []string
		for _, pt := range p.scene.Gamedata.Found {
			item = append(item, pt.Name+" ")
		}
		fmt.Print(item, "\n")
	}

	return nil
}

// test
var _ gameobject.Component = &Player{}
