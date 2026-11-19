package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var RAMExample = register(func() *game.Part {
	return &game.Part{
		Name:        "RAM_Example",
		Kind:        game.RAM,
		Rarity:      game.S2,
		Weight:      100,
		Val:         5000,
		Rank:        game.B,
		Performance: 1000,
		Detail:      "RAM の仮パーツデータ",
		ImgName:     "Parts/みんげー_メモリ",
		ImgScale:    RAMImageScale,
		ImgOffset:   &NoOffset,
	}
})
