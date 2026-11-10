package game

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"

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

	Message    string //ゲーム内メッセージ用
	Message_sn uint   //メッセージの変更検知用serial number

	//以下ゲーム用データ

	//持ち運び可能な残り荷物重量
	Weight int

	//所持金
	Money int

	//所持品
	Parts Parts

	//発見した品
	Found Parts

	// 組み立てに使うパーツ
	Assemble Parts
}

// Gamedata初期値
const (
	InitialWeight = 5000
	InitialMoney  = 50000
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

	Result_SelectParts
	Result_Assembled
	Result_Benchmark
	Result_QuitGame
)

const states = 9 //stateの数

type state uint

func NewGamedata() *Gamedata {
	return &Gamedata{
		Timer:  tick.NewTimer(),
		Weight: InitialWeight,
		Money:  InitialMoney,
	}
}

func (g *Gamedata) ForNewScene(state state) {
	g.InitialProcess = true
	g.Flag = []bool{}
	g.State = state
	g.Text = ""
}

// ゲームサイクル終了時のゲームデータ初期化
func (g *Gamedata) ResetGame() {
	g.Parts = []*Part{}
	g.Found = []*Part{}
	g.Assemble = []*Part{}
	g.Weight = InitialWeight
	g.Money = InitialMoney
}

// 同じScene内での利用を想定
func (g *Gamedata) NextState() error {
	if g.State >= states {
		return fmt.Errorf("未知のstateに遷移しようとしています")
	}
	g.State++
	return nil
}

func (g *Gamedata) SetMessage(msg string) {
	g.Message = msg
	g.Message_sn++
}
