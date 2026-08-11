package gameobject

import scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"

type Component interface {
	Update(obj scene.GameObject, active bool) (err error)
}
