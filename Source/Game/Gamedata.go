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
	Trash
)

type kinds uint

// パーツの状態ランク
// Value倍率を定義
const (
	A rank = 2
	B rank = 1.5
	C rank = 1
	D rank = 0.5
)

type rank float32

// PCパーツ用の構造体
type Part struct {
	Name        string
	Kind        kinds
	Weight      int
	Val         int
	Rank        rank
	Performance int
}

func NewPart() *Part {
	return &Part{}
}

func (g *Gamedata) AddPart(p *Part) string {
	if p == nil {
		return "ごみを見つけた"
	}
	g.Parts = append(g.Parts, p)
	return p.Name + "を見つけた"
}

func (g *Gamedata) TotalWeight() int {
	var wg int
	for _, p := range g.Parts {
		wg += p.Weight
	}
	return wg
}

func (p *Part) Value() int {
	return p.Val * int(p.Rank)
}

func (g *Gamedata) TotalValue() int {
	var v int
	for _, p := range g.Parts {
		v += p.Value()
	}
	return v
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
