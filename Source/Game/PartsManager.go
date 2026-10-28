package game

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
