package component

import gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"

type HogeComponent struct {
}

func (h HogeComponent) Update(obj gameobject.GameObject, active bool) (err error) {
	// なにもしない
	return nil
}

func NewHogeComponent() (*HogeComponent, error) {
	return &HogeComponent{}, nil
}
