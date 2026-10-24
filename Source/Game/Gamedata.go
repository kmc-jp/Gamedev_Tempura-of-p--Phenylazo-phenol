package game

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// シーン間共有データ
type Gamedata struct {
	//文字描画用の基本フォント格納
	Font *text.GoTextFaceSource

	AssetsManager *assets.AssetsManager

	//シーン読み込み時の初回実行制御用
	InitialProcess bool

	Timer *tick.Timer

	Flag []bool

	State state

	//仮UIでのテスト用テキストデータ
	Text string

	//以下ゲーム用データ

	//荷物重量
	weight int

	//所持金
	money int

	//所持品
	Parts []*Part
}

// Flag 識別名
const ()

// state識別名
const (
	Title_Default state = iota

	Play_Explore
	Play_Timeup

	Train_Default

	Result_Default
)

type state uint

// パーツ種別
const (
	CPU kinds = iota
	GPU
	RAM

	Others
)

type kinds uint

// パーツの状態ランク
const (
	A rank = iota
	B
	C
	D
)

type rank uint

// PCパーツ用の構造体
type Part struct {
	Name        string
	Kind        kinds
	Weight      int
	Value       int
	Rank        rank
	Performance int
}

func NewPart() *Part {
	return &Part{}
}

func NewGamedata() *Gamedata {
	return &Gamedata{}
}

func (g *Gamedata) ForNewScene(state state) {
	g.InitialProcess = true
	g.Timer = &tick.Timer{}
	g.Flag = []bool{}
	g.State = state
	g.Text = ""
}
