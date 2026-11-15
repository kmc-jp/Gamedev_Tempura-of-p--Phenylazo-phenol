// Package selection は、選択可能UIの「いまどれを選んでいるか」と
// 「どの入力でどこへ移るか」を管理する共通基盤を提供する。
//
// オートマトン（状態遷移グラフ）として扱い、遷移先と遷移条件を登録しておけば
// Fire でトリガーを渡すだけで選択が移る。描画や入力の取得は行わない。
package selection

// Direction は上下左右の選択移動を表す汎用トリガーである。
type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

// Directions は全方向を固定順で返す（順序を安定させるため map を使わない）。
func Directions() []Direction { return []Direction{Left, Right, Up, Down} }

type edgeKey[N comparable, T comparable] struct {
	from    N
	trigger T
}

type transition[N comparable] struct {
	to   N
	cond func() bool // nil なら無条件
}

// Graph は N をノード（選択肢）、T をトリガー（入力の意味）とする選択グラフである。
type Graph[N comparable, T comparable] struct {
	current     N
	transitions map[edgeKey[N, T]][]transition[N]
	guard       func(from N, trigger T, to N) bool
	revision    int
}

// NewGraph は initial を現在のノードとする空のグラフを生成する。
func NewGraph[N comparable, T comparable](initial N) *Graph[N, T] {
	return &Graph[N, T]{
		current:     initial,
		transitions: map[edgeKey[N, T]][]transition[N]{},
	}
}

// AddTransition は from で trigger を受けたとき to へ移る遷移を登録する。
// cond が nil なら無条件。同じ (from, trigger) に複数登録した場合は登録順に評価し、
// 最初に条件を満たしたものを採用する。
func (g *Graph[N, T]) AddTransition(from N, trigger T, to N, cond func() bool) {
	key := edgeKey[N, T]{from: from, trigger: trigger}
	g.transitions[key] = append(g.transitions[key], transition[N]{to: to, cond: cond})
}

// SetGuard は全遷移に共通の条件を設定する（差し替え可能。nil で解除）。
// 遷移ごとの cond を満たした候補に対して guard(from, trigger, to) を評価し、
// false ならその候補を不採用にして次の候補へ進む。
// 「隣接駅のみ」「訪問済みは不可」のようなルールを、遷移の登録をやり直さずに入れ替えるために使う。
func (g *Graph[N, T]) SetGuard(guard func(from N, trigger T, to N) bool) {
	g.guard = guard
}

// SetCondition は登録済みの遷移 (from, trigger, to) の条件を差し替える（nil で無条件に戻す）。
// 同じ (from, trigger, to) が複数登録されている場合はすべて差し替える。
// 該当する遷移がなければ false を返す。
func (g *Graph[N, T]) SetCondition(from N, trigger T, to N, cond func() bool) bool {
	key := edgeKey[N, T]{from: from, trigger: trigger}
	found := false
	candidates := g.transitions[key]
	for i := range candidates {
		if candidates[i].to == to {
			candidates[i].cond = cond
			found = true
		}
	}
	return found
}

// Next は Fire した場合の遷移先を、状態を変えずに返す（UIのヒント表示・テスト用）。
// 遷移先がなければ false を返す。
func (g *Graph[N, T]) Next(trigger T) (N, bool) {
	key := edgeKey[N, T]{from: g.current, trigger: trigger}
	for _, t := range g.transitions[key] {
		if t.cond != nil && !t.cond() {
			continue
		}
		if g.guard != nil && !g.guard(g.current, trigger, t.to) {
			continue
		}
		return t.to, true
	}
	var zero N
	return zero, false
}

// Fire は現在のノードで trigger を発火する。遷移したら true を返し revision を増やす。
// 該当する遷移がない・条件を満たす遷移がないときは何もせず false。
func (g *Graph[N, T]) Fire(trigger T) bool {
	to, ok := g.Next(trigger)
	if !ok {
		return false
	}
	g.current = to
	g.revision++
	return true
}

// Current は現在のノードを返す。
func (g *Graph[N, T]) Current() N {
	return g.current
}

// SetCurrent は現在のノードを直接設定する。値が変わったときだけ revision を増やす。
func (g *Graph[N, T]) SetCurrent(n N) {
	if g.current == n {
		return
	}
	g.current = n
	g.revision++
}

// Revision は選択が変化するたびに増える整数を返す。UIが変化を観測するために使う。
func (g *Graph[N, T]) Revision() int {
	return g.revision
}
