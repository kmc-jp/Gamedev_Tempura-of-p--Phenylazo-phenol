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
