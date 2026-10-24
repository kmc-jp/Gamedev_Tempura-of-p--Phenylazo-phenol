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
}

// Flag 識別名
const ()

// state識別名
const (
	Play_Explore state = iota
	Play_Timeup

	Train_Default

	Result_Default
)

type state uint

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
