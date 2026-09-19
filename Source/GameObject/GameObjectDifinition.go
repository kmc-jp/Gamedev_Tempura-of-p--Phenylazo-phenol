package gameobject

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

// 個々の GameObject の定義を書く

func NewHogeObject(scene *scene.PlayScene, NewTransform Utils.Factory[Transform], NewRender Utils.Factory[Render]) Utils.Factory[GameObject] {
	return func() (*GameObject, error) {
		t, err := NewTransform()
		if err != nil {
			return nil, fmt.Errorf("NewTransform() でエラーが発生しました。\n %w", err)
		}
		r, err := NewRender()
		if err != nil {
			return nil, fmt.Errorf("NewRender() でエラーが発生しました。\n %w", err)
		}
		h := GameObject{
			Scene: scene,
			tf:    *t,
			rd:    *r,
			cmps:  []Component{},
		}
		return &h, nil
	}
}

// test
var _ Utils.Factory[GameObject] = NewHogeObject(nil, nil, nil)
