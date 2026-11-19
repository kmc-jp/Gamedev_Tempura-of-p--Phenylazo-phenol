package parts

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"strconv"
)

type ramModelSpec struct {
	capacityGB  int
	performance int
	value       int
	rarityLevel int
}

func registerRAM(spec ramModelSpec) {
	rarity := game.S1
	if spec.rarityLevel >= 2 {
		rarity = game.S2
	}
	if spec.rarityLevel >= 3 {
		rarity = game.S3
	}
	register(func() *game.Part {
		capacity := strconv.Itoa(spec.capacityGB) + "GB"
		return &game.Part{
			Name:        "メモリ " + capacity,
			Kind:        game.RAM,
			Rarity:      rarity,
			Weight:      100,
			Val:         spec.value,
			Performance: spec.performance,
			Detail:      capacity + " のメモリモジュール。",
			ImgName:     "Parts/みんげー_メモリ",
			ImgScale:    RAMImageScale,
			ImgOffset:   &NoOffset,
		}
	})
}

func init() {
	capacities := []struct {
		gb          int
		performance int
		value       int
		rarityLevel int
	}{
		{gb: 4, performance: 200, value: 1000, rarityLevel: 1},
		{gb: 8, performance: 400, value: 1800, rarityLevel: 1},
		{gb: 16, performance: 650, value: 3000, rarityLevel: 1},
		{gb: 24, performance: 850, value: 4000, rarityLevel: 2},
		{gb: 32, performance: 1000, value: 5200, rarityLevel: 2},
		{gb: 48, performance: 1000, value: 6900, rarityLevel: 2},
		{gb: 64, performance: 1000, value: 8700, rarityLevel: 2},
		{gb: 96, performance: 1000, value: 11200, rarityLevel: 3},
		{gb: 128, performance: 1000, value: 14000, rarityLevel: 3},
	}

	for range 5 {
		for _, spec := range capacities {
			registerRAM(ramModelSpec{
				capacityGB:  spec.gb,
				performance: spec.performance,
				value:       spec.value,
				rarityLevel: spec.rarityLevel,
			})
		}
	}
}
