package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var (
	GPUBlack = register(func() *game.Part {
		return &game.Part{
			Name:        "黒GPU",
			Kind:        game.GPU,
			Subcategory: game.GPU_Black,
			Rarity:      game.S2,
			Weight:      1000,
			Val:         20000,
			Rank:        game.C,
			Performance: 1000,
			Detail:      "GPU の仮パーツデータ",
			ImgName:     "Parts/みんげー_グラボ",
			ImgScale:    GPUBlackImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	GPUWhite = register(func() *game.Part {
		return &game.Part{
			Name:        "白GPU",
			Kind:        game.GPU,
			Subcategory: game.GPU_White,
			Rarity:      game.S2,
			Weight:      1000,
			Val:         20000,
			Rank:        game.C,
			Performance: 1000,
			Detail:      "GPU 白 の仮パーツデータ",
			ImgName:     "Parts/みんげー_グラボ_白",
			ImgScale:    GPUWhiteImageScale,
			ImgOffset:   &GPUWhiteOffset,
		}
	})
)
