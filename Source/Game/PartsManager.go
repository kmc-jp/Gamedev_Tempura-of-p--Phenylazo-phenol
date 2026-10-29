package game

import "math/rand"

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
	A rank = 4
	B rank = 3
	C rank = 2
	D rank = 1
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

//購入処理結果
const (
	Buy purchase = iota

	Lack_money
	Lack_weight

	NoItem
)

type purchase uint

func (g *Gamedata) BuyPart() purchase {
	if len(g.Found) == 0 {
		return NoItem
	}

	p := g.Found[0]

	if g.Money < p.Val {
		return Lack_money
	}
	if g.Weight < p.Weight {
		return Lack_weight
	}

	g.AddPartToPartsList(p)
	g.Money -= p.Val
	g.Weight -= p.Weight
	g.Found = g.Found[1:] //Foundの先頭を削除
	return Buy
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

func (g *Gamedata) ExtractFoundPartsName() []string {
	var names []string
	for _, p := range g.Found {
		names = append(names, p.Name)
	}
	return names
}

func (p *Part) Score() float32 {
	s := float32(p.Performance) + float32(p.Performance)*(rand.Float32()-0.5)*0.1
	return s
}
