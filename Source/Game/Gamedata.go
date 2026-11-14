package game

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"
	"slices"

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

	//世界時間用のタイマー
	WorldTimer *tick.WorldTimer

	Flag []bool

	nextStage Utils.Factory[Scene]

	currentLocation stageinfo.Location

	State state

	//世界時間進行の有無
	timeProgressing bool

	//仮UIでのテスト用テキストデータ
	Text string

	Message    string //ゲーム内メッセージ用
	Message_sn uint   //メッセージの変更検知用serial number

	//以下ゲーム用データ

	//持ち運び可能な残り荷物重量
	Weight int

	//所持金
	Money int

	//訪れた店の数
	Visited uint

	// 所持品。書き換えは AddParts / RemoveParts を通す
	parts Parts

	// 発見した品。書き換えは AddFoundParts / RemoveFoundParts を通す
	found Parts

	// partsRevision は parts / found が変化するたびに増える。UIが変化を観測するために使う。
	partsRevision int

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

	SelectStage_Default

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
		Timer:           tick.NewTimer(),
		WorldTimer:      tick.NewWorldTimer(),
		Weight:          InitialWeight,
		Money:           InitialMoney,
		currentLocation: stageinfo.Departure,
	}
}

func (g *Gamedata) ForNewScene(state state) {
	g.InitialProcess = true
	g.Flag = []bool{}
	g.State = state
	g.Text = ""
}

func (g *Gamedata) ForNewStage() {
	g.Visited++
}

// ゲームサイクル終了時のゲームデータ初期化
func (g *Gamedata) ResetGame() {
	g.parts = Parts{}
	g.found = Parts{}
	g.partsRevision++
	g.Assemble = []*Part{}
	g.Weight = InitialWeight
	g.Money = InitialMoney
	g.Visited = 0
	g.currentLocation = stageinfo.Departure
	g.timeProgressing = false
}

// 同じScene内での利用を想定
func (g *Gamedata) NextState() error {
	if g.State >= states {
		return fmt.Errorf("未知のstateに遷移しようとしています")
	}
	g.State++
	return nil
}

func (g *Gamedata) SetNextStage(nextStage Utils.Factory[Scene]) {
	g.nextStage = nextStage
}

func (g *Gamedata) NextStage() Utils.Factory[Scene] {
	return g.nextStage
}

func (g *Gamedata) EnableWorldTimeProgressing() {
	g.timeProgressing = true
}

func (g *Gamedata) DisableWorldTimeProgressing() {
	g.timeProgressing = false
}

func (g *Gamedata) WorldTimeProgress() bool {
	return g.timeProgressing
}

// Found / Parts 関連
// 一覧の書き換えはここのメソッドに限定し、変化するたびに partsRevision を増やす。
// 一覧内のポインタは一意である前提で、削除はポインタで要素を特定する。

// FoundParts は発見した品の一覧を返す。返り値はコピーなので、書き換えても Gamedata には反映されない。
func (g *Gamedata) FoundParts() Parts {
	return slices.Clone(g.found)
}

// OwnedParts は所持品の一覧を返す。返り値はコピーなので、書き換えても Gamedata には反映されない。
func (g *Gamedata) OwnedParts() Parts {
	return slices.Clone(g.parts)
}

// AddFoundParts は p を発見した品に追加する。p が nil の場合は何もしない。
func (g *Gamedata) AddFoundParts(p *Part) {
	if p == nil {
		return
	}
	g.found = append(g.found, p)
	g.partsRevision++
}

// RemoveFoundParts は p を発見した品から取り除く。p が含まれていなければ false を返す。
func (g *Gamedata) RemoveFoundParts(p *Part) bool {
	return g.removePart(&g.found, p)
}

// AddParts は p を所持品に追加する。p が nil の場合は何もしない。
func (g *Gamedata) AddParts(p *Part) {
	if p == nil {
		return
	}
	g.parts = append(g.parts, p)
	g.partsRevision++
}

// RemoveParts は p を所持品から取り除く。p が含まれていなければ false を返す。
func (g *Gamedata) RemoveParts(p *Part) bool {
	return g.removePart(&g.parts, p)
}

func (g *Gamedata) removePart(list *Parts, p *Part) bool {
	if p == nil {
		return false
	}
	index := slices.Index(*list, p)
	if index < 0 {
		return false
	}
	*list = slices.Delete(*list, index, index+1)
	g.partsRevision++
	return true
}
