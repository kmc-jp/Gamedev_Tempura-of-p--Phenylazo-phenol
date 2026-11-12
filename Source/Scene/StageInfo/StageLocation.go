package stageinfo

import "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

type StageLocation struct {
	X float64
	Y float64
}

// 各ステージの相対位置 (適当)
var Locations []StageLocation = []StageLocation{
	{0, 0},    //Stage1
	{50, 30},  //Red
	{40, 130}, //Blue
	{200, 50}, //Green
	{120, 70}, //Orange
	{80, 170}, //Purple
}

func (sl StageLocation) VectorValue() Utils.Vec2 {
	return Utils.Vec2{X: sl.X, Y: sl.Y}
}
