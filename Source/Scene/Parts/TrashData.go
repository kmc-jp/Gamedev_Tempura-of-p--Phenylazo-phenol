package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var Trash = register(func() *game.Part {
	return &game.Part{
		Name: "Part", Kind: game.Trash,
		Weight: 0, Val: 0, Rank: game.D, Performance: 0,
	}
})

var Spare = register(func() *game.Part {
	return &game.Part{
		Name:        "あまりもの",
		Kind:        game.Trash,
		Weight:      0,
		Val:         0,
		Rank:        game.D,
		Performance: 100,
		Detail:      "家にあった謎のあまりパーツ",
		IsSpare:     true,
	}
})
