package scene

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

func PartExample() *game.Part {
	return &game.Part{
		Name:        "example",
		Kind:        game.Others,
		Weight:      500,
		Val:         1000,
		Rank:        game.B,
		Performance: 1000,
	}
}
