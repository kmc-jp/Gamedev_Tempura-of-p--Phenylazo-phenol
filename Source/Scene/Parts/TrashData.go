package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var Trash = register(func() *game.Part {
	return &game.Part{
		Name: "Part", Kind: game.Trash,
		Weight: 0, Val: 0, Rank: game.D, Performance: 0,
	}
})
