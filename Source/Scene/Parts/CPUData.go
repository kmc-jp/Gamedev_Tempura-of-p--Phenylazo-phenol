package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var CPUExample = register(func() *game.Part {
	return &game.Part{
		Name:        "CPU_Example",
		Kind:        game.CPU,
		Rarity:      game.S2,
		Weight:      100,
		Val:         10000,
		Rank:        game.A,
		Performance: 1000,
		Detail:      "CPU の仮パーツデータ",
		ImgName:     "Parts/みんげー_CPU",
		ImgScale:    CPUImageScale,
		ImgOffset:   &NoOffset,
	}
})
