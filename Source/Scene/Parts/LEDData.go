package parts

import game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

var LED = register(func() *game.Part {
	return &game.Part{
		Name: "LED_Tape", Kind: game.LED,
		Weight: 100, Val: 100, Rank: game.A, Performance: 0,
		Detail: "PC装飾用LEDテープ", ImgName: "Parts/みんげー_LEDテープ",
		ImgScale: LEDImageScale, ImgOffset: &NoOffset,
	}
})

func init() {
	for range 29 {
		register(LED)
	}
}
