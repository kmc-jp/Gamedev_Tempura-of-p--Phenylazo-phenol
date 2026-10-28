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

	//持ち運び可能な残り荷物重量
	Weight int

	//所持金
	Money int

	//所持品
	Parts []*Part

	//発見した品
	Found []*Part
}

// Gamedata初期値
const (
	InitialWeight = 1000
	InitialMoney  = 10000
)

// Flag 識別名
const ()

// state識別名
const (
	Title_Default state = iota

	Play_Explore
	Play_Timeup
	Play_Select

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

// ゲームサイクル終了時のゲームデータ初期化
func (g *Gamedata) ResetGame() {
	g.Parts = []*Part{}
	g.Found = []*Part{}
	g.Weight = InitialWeight
	g.Money = InitialMoney
}
