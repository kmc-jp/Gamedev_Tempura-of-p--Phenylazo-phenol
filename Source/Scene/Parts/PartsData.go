package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var (
	Trash = register(func() *game.Part {
		return &game.Part{
			Name:        "Part",
			Kind:        game.Trash,
			Weight:      0,
			Val:         0,
			Rank:        game.D,
			Performance: 0,
		}
	})

	CPUExample = register(func() *game.Part {
		return &game.Part{
			Name:        "CPU_Example",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         10000,
			Rank:        game.A,
			Performance: 1000,
			Detail:      "CPU の仮パーツデータ",
			ImgName:     "みんげー_CPU",
		}
	})

	GPUExample = register(func() *game.Part {
		return &game.Part{
			Name:        "GPU_Example",
			Kind:        game.GPU,
			Rarity:      game.S2,
			Weight:      1000,
			Val:         20000,
			Rank:        game.C,
			Performance: 1000,
			Detail:      "GPU の仮パーツデータ",
			ImgName:     "みんげー_グラボ",
		}
	})

	RAMExample = register(func() *game.Part {
		return &game.Part{
			Name:        "RAM_Example",
			Kind:        game.RAM,
			Rarity:      game.S2,
			Weight:      100,
			Val:         5000,
			Rank:        game.B,
			Performance: 1000,
			Detail:      "RAM の仮パーツデータ",
			ImgName:     "みんげー_メモリ",
		}
	})

	MBExample = register(func() *game.Part {
		return &game.Part{
			Name:        "MB_Example",
			Kind:        game.MB,
			Weight:      500,
			Val:         2000,
			Rank:        game.C,
			Performance: 1000,
		}
	})

	CaseExample = register(func() *game.Part {
		return &game.Part{
			Name:        "Case_Example",
			Kind:        game.Case,
			Weight:      2000,
			Val:         2000,
			Rank:        game.B,
			Performance: 1000,
		}
	})

	// 以下、任意パーツ
	LED = register(func() *game.Part {
		return &game.Part{
			Name:        "LED_Tape",
			Kind:        game.LED,
			Weight:      100,
			Val:         100,
			Rank:        game.A,
			Performance: 0,
		}
	})
)
