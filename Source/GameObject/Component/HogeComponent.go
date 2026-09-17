package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type HogeComponent struct {
	GameObject *gameobject.GameObject
}

func (h HogeComponent) Update(obj gameobject.GameObject, active bool) (err error) {
	// なにもしない
	return nil
}

func NewHogeComponent(gameobject *gameobject.GameObject) Utils.Factory[HogeComponent] {
	return func() (*HogeComponent, error) {
		c := HogeComponent{
			GameObject: gameobject,
		}
		return &c, nil
	}
}

// test
var _ Utils.Factory[HogeComponent] = NewHogeComponent(nil)
