package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var (
	empty = register(func() *game.Part {
		return &game.Part{}
	})

	CPUExample = register(func() *game.Part {
		return &game.Part{
			Name:        "CPU_Example",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         10000,
			Performance: 1000,
			Detail:      "CPU の仮パーツデータ",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU1 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i5 3500",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         10000,
			Performance: 5000,
			Detail:      "Nitel社のCPU。いろいろなパソコンに使われている。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU2 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i5 3500T",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         12000,
			Performance: 5100,
			Detail:      "Nitel社のCPU。いろいろなパソコンに使われている。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU3 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i7 3700",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         15000,
			Performance: 7000,
			Detail:      "Nitel社のCPU。いろいろな高性能パソコンに使われている。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU4 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i7 3700T",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         17000,
			Performance: 7100,
			Detail:      "Nitel社のCPU。いろいろな高性能パソコンに使われている。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU5 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i9 3900",
			Kind:        game.CPU,
			Rarity:      game.S3,
			Weight:      100,
			Val:         28000,
			Performance: 9000,
			Detail:      "Nitel社のCPU。いろいろなゲーミングパソコンに使われている。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU6 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i9 3990TS",
			Kind:        game.CPU,
			Rarity:      game.S3,
			Weight:      100,
			Val:         35000,
			Performance: 9500,
			Detail:      "Nitel社のCPU。高性能ゲーミングパソコンならこれ。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU7 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i3 2100",
			Kind:        game.CPU,
			Rarity:      game.S1,
			Weight:      100,
			Val:         5000,
			Performance: 2500,
			Detail:      "Nitel社のCPU。性能は微妙かも。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU8 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i5 2500",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         7000,
			Performance: 3900,
			Detail:      "Nitel社のCPU。まあこれで十分。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU9 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i7 2700",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         10000,
			Performance: 6000,
			Detail:      "Nitel社のCPU。いろいろなパソコンに使われている。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU10 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i9 2900",
			Kind:        game.CPU,
			Rarity:      game.S3,
			Weight:      100,
			Val:         20000,
			Performance: 8500,
			Detail:      "Nitel社のCPU。ゲーミングパソコンに使われている。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CPU11 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i5 550",
			Kind:        game.CPU,
			Rarity:      game.S1,
			Weight:      100,
			Val:         5000,
			Performance: 2600,
			Detail:      "Nitel社のCPU。古いけどそこそこ使える。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}

	})

	CPU12 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i7 750",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         7000,
			Performance: 3600,
			Detail:      "Nitel社のCPU。古いけどそこそこ使える。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}

	})

	CPU13 = register(func() *game.Part {
		return &game.Part{
			Name:        "Nitel Pure i9 950",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         10000,
			Performance: 3500,
			Detail:      "Nitel社のCPU。古いけどまあ使える。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}

	})

	CPU31 = register(func() *game.Part {
		return &game.Part{
			Name:        "Ryzin5 3400U",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         10000,
			Performance: 5000,
			Detail:      "ADD社のCPU。いろいろなパソコンに入ってる。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}

	})

	CPU32 = register(func() *game.Part {
		return &game.Part{
			Name:        "Ryzin5 3400U",
			Kind:        game.CPU,
			Rarity:      game.S2,
			Weight:      100,
			Val:         10000,
			Performance: 5000,
			Detail:      "ADD社のCPU。いろいろなパソコンに入ってる。",
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}

	})
)
