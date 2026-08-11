package component

import scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"

type HogeComponent struct {
}

func (h HogeComponent) Update(obj scene.GameObject, active bool) (err error) {
	// なにもしない
	return nil
}

func NewHogeComponent() (*HogeComponent, error) {
	return &HogeComponent{}, nil
}
