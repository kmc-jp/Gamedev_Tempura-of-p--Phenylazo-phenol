package parts

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"math/rand"
)

var registry []func() *game.Part

func register(builder func() *game.Part) func() *game.Part {
	registry = append(registry, builder)
	return builder
}

func RandomPartsFromAll() *game.Part {
	if len(registry) == 0 {
		return nil
	}
	p := registry[rand.Intn(len(registry))]()
	if p.Kind == game.Trash {
		return nil
	}
	return p
}
