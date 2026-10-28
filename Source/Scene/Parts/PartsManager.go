package parts

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
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

	//重みづけ抽選をしない場合はここで抽選
	if !Utils.Percent(60) {
		return nil
	}

	p := registry[rand.Intn(len(registry))]()
	if p.Kind == game.Trash {
		return nil
	}
	return p
}
