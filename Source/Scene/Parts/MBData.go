package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var (
	MBBlack = register(func() *game.Part {
		return &game.Part{
			Name: "マザーボード黒", Kind: game.MB, Subcategory: game.MB_Black,
			Weight: 500, Val: 2000, Rank: game.C, Performance: 1000,
			Detail: "MB黒 の仮パーツデータ", ImgName: "Parts/みんげー_マザボ",
			ImgScale: MBImageScale, ImgOffset: &NoOffset,
		}
	})
	MBWhite = register(func() *game.Part {
		return &game.Part{
			Name: "マザーボード白", Kind: game.MB, Subcategory: game.MB_White,
			Weight: 500, Val: 2000, Rank: game.C, Performance: 1000,
			Detail: "MB白 の仮パーツデータ", ImgName: "Parts/みんげー_マザボ_白",
			ImgScale: MBImageScale, ImgOffset: &NoOffset,
		}
	})
	MBBlue = register(func() *game.Part {
		return &game.Part{
			Name: "マザーボード青", Kind: game.MB, Subcategory: game.MB_Blue,
			Weight: 500, Val: 2000, Rank: game.C, Performance: 1000,
			Detail: "MB青 の仮パーツデータ", ImgName: "Parts/みんげー_マザボ_青",
			ImgScale: MBImageScale, ImgOffset: &NoOffset,
		}
	})
	MBGreen = register(func() *game.Part {
		return &game.Part{
			Name: "マザーボード緑", Kind: game.MB, Subcategory: game.MB_Green,
			Weight: 500, Val: 2000, Rank: game.C, Performance: 1000,
			Detail: "MB緑 の仮パーツデータ", ImgName: "Parts/みんげー_マザボ_緑",
			ImgScale: MBImageScale, ImgOffset: &NoOffset,
		}
	})
)
