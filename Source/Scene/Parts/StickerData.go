package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var Sticker = register(func() *game.Part {
	return &game.Part{
		Name: "ステッカー", Kind: game.Sticker,
		Weight: 100, Val: 100, Rank: game.A, Performance: 0,
		Detail: "キャラクターのステッカー", ImgName: "Parts/みんげー_キャラステッカー",
		ImgScale: StickerImageScale, ImgOffset: &NoOffset,
	}
})
