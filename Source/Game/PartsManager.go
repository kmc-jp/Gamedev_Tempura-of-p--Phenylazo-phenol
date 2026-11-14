package game

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// パーツ種別 kind
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

// 小分類 subcategory
const (
	GPU_Black subcategory = iota
	GPU_White

	MB_Black
	MB_White

	Case_Black
	Case_White
)

const (
	AvailablePartsKind = 5 // 有効なパーツ種の数
	ArbitraryPartsKind = 1 // 任意のパーツ種の数
)

var AvailableParts = [AvailablePartsKind]kinds{CPU, GPU, RAM, MB, Case}

type kinds uint
type subcategory uint

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
	Subcategory subcategory
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
	g.AddFoundParts(p)
	return p.Name + "を見つけた"
}

// PartsRevision は所持品 / 発見した品の版を返す。どちらかが変化するたびに増える。
// UIはこの値の変化を観測し、変化したときだけ一覧を作り直す。
func (g *Gamedata) PartsRevision() int {
	return g.partsRevision
}

// 購入処理結果
const (
	Buy PurchaseResult = iota

	Lack_money
	Lack_weight

	NoItem
)

// 購入処理の結果
type PurchaseResult uint

// BuyFoundPart は発見済みパーツ p を購入し、発見した品から所持品へ移す。
// p が nil、または発見した品に含まれない場合は NoItem を返す。
func (g *Gamedata) BuyFoundPart(p *Part) PurchaseResult {
	if p == nil || !slices.Contains(g.found, p) {
		return NoItem
	}
	if g.Money < p.Val {
		return Lack_money
	}
	if g.Weight < p.Weight {
		return Lack_weight
	}

	g.RemoveFoundParts(p)
	g.AddParts(p)
	g.Money -= p.Val
	g.Weight -= p.Weight
	return Buy
}

func (g *Gamedata) TotalWeight() int {
	var wg int
	for _, p := range g.parts {
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
	for _, p := range g.parts {
		v += p.Value()
	}
	return v
}

func (g *Gamedata) ExtractFoundPartsName() []string {
	var names []string
	for _, p := range g.found {
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

// AvailableParts の順にKindでソート
func (pts Parts) SortWithKind() Parts {
	var sorted Parts
	var eliminated Parts

	//ptsから各Kindを抽出, nil化
	for i := 0; i < AvailablePartsKind; i++ {
		for j, v := range pts {
			if v.Kind == AvailableParts[i] {
				sorted = append(sorted, v)
				pts[j] = nil
			}
		}
	}

	//残りの要素を追加
	sorted = append(sorted, pts...)
	for _, v := range sorted {
		if v != nil {
			eliminated = append(eliminated, v)
		}
	}
	return eliminated
}

func (pts Parts) Delete(idx int) Parts {
	pts[idx] = nil
	var deleted Parts
	for _, v := range pts {
		if v == nil {
			continue
		}
		deleted = append(deleted, v)
	}
	return deleted
}
