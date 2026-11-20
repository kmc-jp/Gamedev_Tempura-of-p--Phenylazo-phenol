package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"math/rand"
	"reflect"
)

// ステージを歩き回るおじさん (ランダムウォーク)
type Ojisan struct {
	tw *TileWalker

	// 次に歩き出すまでの残りフレーム数
	wait int
	// 待ちフレーム数の範囲 [min, max]
	min int
	max int
}

var ojisanDirections = []Utils.Vec2Int{
	Utils.NewVecIntUp(),
	Utils.NewVecIntDown(),
	Utils.NewVecIntLeft(),
	Utils.NewVecIntRight(),
}

// test
var _ gameobject.Component = &Ojisan{}

func (o *Ojisan) Update(active bool) (err error) {
	if !active {
		return nil
	}

	// 歩いている間は何もしない
	if o.tw.IsMoving() {
		return nil
	}

	o.wait--
	if o.wait > 0 {
		return nil
	}

	// 動けなくても向きが変わるだけなので結果は無視してよい
	o.tw.TryMove(ojisanDirections[rand.Intn(len(ojisanDirections))])
	o.wait = o.nextWait()

	return nil
}

func (o *Ojisan) nextWait() int {
	return o.min + rand.Intn(o.max-o.min+1)
}

// min, max : 次に歩き出すまでの待ちフレーム数の範囲
func NewOjisan(obj *gameobject.GameObject, min int, max int) Utils.Factory[*Ojisan] {
	return func() (*Ojisan, error) {
		if min < 1 || max < min {
			return nil, fmt.Errorf("待ちフレーム数の範囲が不正です。 1 <= min <= max を満たしてください。\nmin = %d, max = %d", min, max)
		}

		tw, err := obj.GetComponent(reflect.TypeFor[*TileWalker]())
		if err != nil {
			return nil, fmt.Errorf("obj.GetComponent(reflect.TypeFor[*TileWalker]()) で失敗しました。\n%w", err)
		}

		o := &Ojisan{
			tw:  tw.(*TileWalker),
			min: min,
			max: max,
		}
		o.wait = o.nextWait() // 全員が同時に歩き出さないようにする

		return o, nil
	}
}
