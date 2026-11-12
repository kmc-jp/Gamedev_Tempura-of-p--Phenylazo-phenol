package stageinfo

type StageLocation struct {
	X uint
	Y uint
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
