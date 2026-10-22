package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
)

type InteractableTile interface {
	gameobject.Component
	Interact() error
}
