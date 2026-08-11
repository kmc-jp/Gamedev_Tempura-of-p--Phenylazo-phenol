package gameobject

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"fmt"
)

type GameObjectUtils struct {
	gameObject scene.GameObject
}

func NewGameObjectUtils(obj scene.GameObject) (*GameObjectUtils, error) {
	if obj == nil {
		return nil, fmt.Errorf("obj が空です。 実装先の GameObject を代入してください。")
	}
	return &GameObjectUtils{gameObject: obj}, nil
}

// func (u GameObjectUtils) GetComponent(target reflect.Type) (component Component, err error) {
// 	if u.gameObject == nil {
// 		return nil, fmt.Errorf("u.obj が未登録です")
// 	}
// 	for _, c := range u.gameObject.components() {
// 		if target == reflect.TypeOf(c) {
// 			return c, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("%s に %s が存在しません。", u.gameObject.Name(), target.Name())
// }

// func (u GameObjectUtils) Update(active bool) (err error) {
// 	obj := u.gameObject
// 	for _, c := range obj.components() {
// 		if err := c.Update(obj, active); err != nil {
// 			return fmt.Errorf("コンポーネント %s でエラーが発生しました。\n%w", reflect.TypeOf(c), err)
// 		}
// 	}
// 	return nil
// }

// func (u GameObjectUtils) Draw(screen *ebiten.Image) {
// 	obj := u.gameObject
// 	obj.render().Draw(obj, screen)
// }
