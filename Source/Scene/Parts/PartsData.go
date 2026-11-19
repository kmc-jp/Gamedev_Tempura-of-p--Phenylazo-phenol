package parts

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

// パーツリスト表示の画像サイズ指定
const (
	CPUImageScale      game.Imgscale = 3
	GPUBlackImageScale game.Imgscale = 2.2
	GPUWhiteImageScale game.Imgscale = 2
	RAMImageScale      game.Imgscale = 2.2
	MBImageScale       game.Imgscale = 2.2
	CaseImageScale     game.Imgscale = 1.2

	LEDImageScale     game.Imgscale = 3
	StickerImageScale game.Imgscale = 2
)

// パーツリスト表示の画像位置オフセットを指定
var (
	GPUWhiteOffset = Utils.Vec2{X: 0, Y: 20}
	NoOffset       = Utils.NewZeroVec2()
)

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

	// CPU
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
			ImgName:     "Parts/みんげー_CPU",
			ImgScale:    CPUImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	// GPU
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

	// RAM
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
			ImgName:     "Parts/みんげー_メモリ",
			ImgScale:    RAMImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	// MB
	MBBlack = register(func() *game.Part {
		return &game.Part{
			Name:        "マザーボード黒",
			Kind:        game.MB,
			Subcategory: game.MB_Black,
			Weight:      500,
			Val:         2000,
			Rank:        game.C,
			Performance: 1000,
			Detail:      "MB黒 の仮パーツデータ",
			ImgName:     "Parts/みんげー_マザボ",
			ImgScale:    MBImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	MBWhite = register(func() *game.Part {
		return &game.Part{
			Name:        "マザーボード白",
			Kind:        game.MB,
			Subcategory: game.MB_White,
			Weight:      500,
			Val:         2000,
			Rank:        game.C,
			Performance: 1000,
			Detail:      "MB白 の仮パーツデータ",
			ImgName:     "Parts/みんげー_マザボ_白",
			ImgScale:    MBImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	MBBlue = register(func() *game.Part {
		return &game.Part{
			Name:        "マザーボード青",
			Kind:        game.MB,
			Subcategory: game.MB_Blue,
			Weight:      500,
			Val:         2000,
			Rank:        game.C,
			Performance: 1000,
			Detail:      "MB青 の仮パーツデータ",
			ImgName:     "Parts/みんげー_マザボ_青",
			ImgScale:    MBImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	MBGreen = register(func() *game.Part {
		return &game.Part{
			Name:        "マザーボード緑",
			Kind:        game.MB,
			Subcategory: game.MB_Green,
			Weight:      500,
			Val:         2000,
			Rank:        game.C,
			Performance: 1000,
			Detail:      "MB緑 の仮パーツデータ",
			ImgName:     "Parts/みんげー_マザボ_緑",
			ImgScale:    MBImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	// Case
	CaseBlackSideGlass = register(func() *game.Part {
		return &game.Part{
			Name:        "PCケース黒",
			Kind:        game.Case,
			Subcategory: game.Case_Black_SideGlass,
			Weight:      2000,
			Val:         2000,
			Rank:        game.B,
			Performance: 1000,
			Detail:      "側面がガラスの黒PCケース",
			ImgName:     "Parts/みんげー_PCケース_黒_蓋閉じてる",
			ImgScale:    CaseImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CaseBlackSideFrontGlass = register(func() *game.Part {
		return &game.Part{
			Name:        "PCケース黒",
			Kind:        game.Case,
			Subcategory: game.Case_Black_SideFrontGlass,
			Weight:      2000,
			Val:         2000,
			Rank:        game.B,
			Performance: 1000,
			Detail:      "側面と前面がガラスの黒PCケース",
			ImgName:     "Parts/みんげー_PCケース_黒_蓋開いてる",
			ImgScale:    CaseImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CaseWhiteSideGlass = register(func() *game.Part {
		return &game.Part{
			Name:        "PCケース白",
			Kind:        game.Case,
			Subcategory: game.Case_White_SideGlass,
			Weight:      2000,
			Val:         2000,
			Rank:        game.B,
			Performance: 1000,
			Detail:      "側面がガラスの白PCケース",
			ImgName:     "Parts/みんげー_PCケース_白_蓋閉じてる",
			ImgScale:    CaseImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	CaseWhiteSideFrontGlass = register(func() *game.Part {
		return &game.Part{
			Name:        "PCケース白",
			Kind:        game.Case,
			Subcategory: game.Case_White_SideFrontGlass,
			Weight:      2000,
			Val:         2000,
			Rank:        game.B,
			Performance: 1000,
			Detail:      "側面と前面がガラスの白PCケース",
			ImgName:     "Parts/みんげー_PCケース_白_蓋開いてる",
			ImgScale:    CaseImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	// 以下、任意パーツ
	// LED
	LED = register(func() *game.Part {
		return &game.Part{
			Name:        "LED_Tape",
			Kind:        game.LED,
			Weight:      100,
			Val:         100,
			Rank:        game.A,
			Performance: 0,
			Detail:      "PC装飾用LEDテープ",
			ImgName:     "Parts/みんげー_LEDテープ",
			ImgScale:    LEDImageScale,
			ImgOffset:   &NoOffset,
		}
	})

	// Sticker
	Sticker = register(func() *game.Part {
		return &game.Part{
			Name:        "ステッカー",
			Kind:        game.Sticker,
			Weight:      100,
			Val:         100,
			Rank:        game.A,
			Performance: 0,
			Detail:      "キャラクターのステッカー",
			ImgName:     "Parts/みんげー_キャラステッカー",
			ImgScale:    StickerImageScale,
			ImgOffset:   &NoOffset,
		}
	})
)
