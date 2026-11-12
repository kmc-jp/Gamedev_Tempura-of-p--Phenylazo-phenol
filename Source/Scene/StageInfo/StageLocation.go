package stageinfo

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type StageLocation struct {
	X float64
	Y float64

	accessibility float64 //周辺の交通手段の豊富さ(= 経路の重みづけ; 使うかわからん) デフォルトは 1
}

const (
	Stage1 Location = iota
	Red
	Blue
	Green
	Orange
	Purple

	Departure
)

type Location uint

// 各ステージの相対位置 (数値は適当)
var Locations []StageLocation = []StageLocation{
	{0, 0, 1},    //Stage1
	{50, 30, 1},  //Red
	{40, 130, 1}, //Blue
	{200, 50, 1}, //Green
	{120, 70, 1}, //Orange
	{80, 170, 1}, //Purple

	{100, 50, 1}, //開始時のプレイヤー出発位置 (仮)
}

type locations []StageLocation

// とりあえず二乗距離 引数にscene.Stageを入力
func (sl locations) Distance(from, to int) float64 {
	f := sl[from].VectorValue()
	t := sl[to].VectorValue()
	return f.DistanceSquared(t) / (sl[from].accessibility * sl[to].accessibility)
}

func (sl StageLocation) VectorValue() Utils.Vec2 {
	return Utils.Vec2{X: sl.X, Y: sl.Y}
}
