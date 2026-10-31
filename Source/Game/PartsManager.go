package game

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// パーツ種別
const (
	CPU kinds = iota
	GPU
	RAM
	MB
	Case

	//ここより上に必須パーツの定義、const hogePartsKind と AvailableParts も要更新 ↓ .
	LED

	//ここより上に任意パーツの定義
	Trash
	Others
)

const (
	AvailablePartsKind = 5 // 有効なパーツ種の数
	ArbitraryPartsKind = 1 // 任意のパーツ種の数
)

var AvailableParts = [AvailablePartsKind]kinds{CPU, GPU, RAM, MB, Case}

type kinds uint

// パーツの状態ランク, UI上でのアルファベット表記はしない
// Value倍率を定義
const (
	A rank = 4 //新品
	B rank = 3 //良
	C rank = 2 //普通
	D rank = 1 //ジャンク
)

type rank uint

// パーツのレアリティ(s : star)
const (
	S0 rarity = iota //原則使用しない
	S1
	S2
	S3
)

type rarity uint

// PCパーツ用の構造体
type Part struct {
	Name        string
	Kind        kinds
	Rarity      rarity
	Weight      int
	Val         int
	Rank        rank
	Performance int

	Detail string

	ImgName string
	Img     *ebiten.Image
}

type Parts []*Part

func NewPart() *Part {
	return &Part{}
}

func (g *Gamedata) AddPartToFoundList(p *Part) string {
	if p == nil {
		return "ごみを見つけた"
	}
	g.Found = append(g.Found, p)
	return p.Name + "を見つけた"
}

func (g *Gamedata) AddPartToPartsList(p *Part) {
	g.Parts = append(g.Parts, p)
}

// 購入処理結果
const (
	Buy purchase = iota

	Lack_money
	Lack_weight

	NoItem
)

type purchase uint

func (g *Gamedata) BuyPart(index int) purchase {
	if len(g.Found) == 0 {
		return NoItem
	}

	p := g.Found[index]

	if g.Money < p.Val {
		return Lack_money
	}
	if g.Weight < p.Weight {
		return Lack_weight
	}

	g.AddPartToPartsList(p)
	g.Money -= p.Val
	g.Weight -= p.Weight
	g.Found = slices.Delete(g.Found, index, index)
	return Buy
}

func (g *Gamedata) TotalWeight() int {
	var wg int
	for _, p := range g.Parts {
		wg += p.Weight
	}
	return wg
}

// 所持品としての価値
func (p *Part) Value() int {
	return p.Val * int(p.Rank)
}

// 所持品としての価値
func (g *Gamedata) TotalValue() int {
	var v int
	for _, p := range g.Parts {
		v += p.Value()
	}
	return v
}

func (g *Gamedata) ExtractFoundPartsName() []string {
	var names []string
	for _, p := range g.Found {
		names = append(names, p.Name)
	}
	return names
}

func (g *Gamedata) ExtractAssemblePartsName() []string {
	var names []string
	for _, p := range g.Assemble {
		names = append(names, p.Name)
	}
	return names
}

// performance に -0.05 ~ 0.05 の倍率
func (p *Part) Score() float32 {
	s := float32(p.Performance) + float32(p.Performance)*(rand.Float32()-0.5)*0.1
	return s
}

// 種類k kindsのパーツを含むならindexとtrue, 含まないなら-1とfalseを返す
func (pts Parts) Contain(k kinds) (int, bool) {
	for i, p := range pts {
		if p.Kind == k {
			return i, true
		}
	}
	return -1, false
}

// 必須パーツがそろっているかの確認
func (pts Parts) Enough() ([]uint, bool) {
	l := AvailableParts
	var lack []uint
	var enough bool = true
	for _, p := range pts {
		if p.Kind >= AvailablePartsKind {
			continue
		} else {
			l[p.Kind] = AvailablePartsKind + ArbitraryPartsKind + 1
		}
	}
	for i, v := range l {
		if v == AvailablePartsKind+ArbitraryPartsKind+1 {
			continue
		} else {
			lack = append(lack, uint(i))
			enough = false
		}
	}
	return lack, enough
}

// 各状態ランクを任意の確率で設定。例: p.SetRandomRank(5, 20, 60, 15) A: 5%, B: 20%, C: 60%, D: 15%
func (p *Part) SetRandomRank(a, b, c, d float32) error {
	if a+b+c+d != 100 {
		return fmt.Errorf("合計確率が100になっていません")
	}
	r := rand.Float32() * 100
	if 0 <= r && r < a {
		p.Rank = A
	} else if a <= r && r < a+b {
		p.Rank = B
	} else if a+b <= r && r < a+b+c {
		p.Rank = C
	} else if a+b+c <= r && r < 100 {
		p.Rank = D
	}
	return nil
}
