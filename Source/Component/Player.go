package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"
	"reflect"

	input "github.com/quasilyte/ebitengine-input"
)

type Player struct {
	ipt   *input.Handler
	tw    *TileWalker
	dir   Utils.Vec2Int
	tf    gameobject.Transform
	scene *scene.PlayScene

	// アニメーション関連
	anim *AnimRender

	timerInteractInterval *tick.Timer
}

const (
	interactInterval = 10 //unit: 0.1 sec
)

func (p *Player) Update(active bool) (err error) {
	// なにもしない
	if !active {
		return nil
	}

	p.timerInteractInterval.Set(interactInterval)

	if p.ipt.ActionIsPressed(inputconfig.ActionMoveLeft) {
		p.TryMove(Utils.NewVecIntLeft())
	}
	if p.ipt.ActionIsPressed(inputconfig.ActionMoveRight) {
		p.TryMove(Utils.NewVecIntRight())
	}
	if p.ipt.ActionIsPressed(inputconfig.ActionMoveUp) {
		p.TryMove(Utils.NewVecIntUp())
	}
	if p.ipt.ActionIsPressed(inputconfig.ActionMoveDown) {
		p.TryMove(Utils.NewVecIntDown())
	}
	if p.ipt.ActionIsJustPressed(inputconfig.ActionInteract) {
		if p.timerInteractInterval.Expired() {
			err := p.interact()
			if err != nil {
				return fmt.Errorf("p.interact() でエラーが発生しました。\n%w", err)
			}
		}
	}

	// カメラをプレイヤーに追従させる
	p.scene.SetCameraCenter(p.tf.Position())

	return nil
}

func NewPlayer(scene *scene.PlayScene, obj *gameobject.GameObject, anim *AnimRender, tf gameobject.Transform) Utils.Factory[*Player] {
	return func() (*Player, error) {
		// TileWalker を取得
		tw, err := obj.GetComponent(reflect.TypeFor[*TileWalker]())
		if err != nil {
			return nil, fmt.Errorf("obj.GetComponent(reflect.TypeFor[*TileWalker]()) で失敗しました。\n%w", err)
		}

		timer := tick.NewTimer()

		pl := &Player{
			tw:                    tw.(*TileWalker),
			ipt:                   scene.Input,
			scene:                 scene,
			anim:                  anim,
			tf:                    tf,
			timerInteractInterval: timer,
		}

		return pl, nil
	}
}

func (p *Player) interact() error {
	mypos := p.tw.Position.Position
	targetpos := mypos.Add(p.dir)

	if tgt, ok := p.tw.tm.Interact[targetpos]; ok {
		err := tgt.Interact(p)
		if err != nil {
			return fmt.Errorf("target.Interact() でエラーが発生しました。\ntergetpos = %d\n%w", targetpos, err)
		}
	}

	return nil
}

func (p *Player) TryMove(Direction Utils.Vec2Int) bool {
	if p.tw.IsMoving() {
		return false
	}

	moved := p.tw.TryMove(Direction)
	p.dir = Direction

	return moved
}

// test
var _ gameobject.Component = &Player{}
