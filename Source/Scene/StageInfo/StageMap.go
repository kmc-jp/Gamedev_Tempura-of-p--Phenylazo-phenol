package stageinfo

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Selection"
	"fmt"
	"maps"
	"slices"
)

// Station は路線図上の駅である。座標は StageMap の基準サイズ上の値。
type Station struct {
	Label string
	Pos   Utils.Position
}

// Line は2駅を結ぶ路線（無向）である。
//
// FromKey / ToKey は、その向きに進むために押すキーを明示する。nil なら駅の位置関係から自動で決める。
//   - FromKey: From の駅で To へ進むキー
//   - ToKey:   To の駅で From へ戻るキー
//
// Minutes は2駅間の所要時間（単位: 分。ゲーム内時刻と同じ）。0 でもよい。
type Line struct {
	From    Location
	To      Location
	FromKey *selection.Direction
	ToKey   *selection.Direction
	Minutes uint64
}

// StageMap は行き先選択画面の路線図の定義である。
type StageMap struct {
	Width    float64
	Height   float64
	Stations map[Location]Station
	Lines    []Line
}

// directionVectors は方向キーに対応する画面上のベクトル（Y は下向き）。
var directionVectors = map[selection.Direction]Utils.Vec2{
	selection.Left:  {X: -1, Y: 0},
	selection.Right: {X: 1, Y: 0},
	selection.Up:    {X: 0, Y: -1},
	selection.Down:  {X: 0, Y: 1},
}

// directionThreshold は、路線の向きが方向キーの向きと近いとみなす内積の下限（60°以内）。
const directionThreshold = 0.5

// name は駅の表示用の名前を返す（エラーメッセージ用）。
func name(loc Location) string {
	if s, ok := Location2String[loc]; ok {
		return s
	}
	return fmt.Sprintf("Location(%d)", loc)
}

// Neighbors は loc と路線でつながった駅を、Lines の定義順で返す。
func (m StageMap) Neighbors(loc Location) []Location {
	var result []Location
	for _, line := range m.Lines {
		switch loc {
		case line.From:
			result = append(result, line.To)
		case line.To:
			result = append(result, line.From)
		}
	}
	return result
}

// ShortestMinutes は from から to までの最短所要時間（分）を返す（ダイクストラ法）。
// from == to なら 0。駅が未定義、または到達できない場合はエラーを返す。
func (m StageMap) ShortestMinutes(from, to Location) (uint64, error) {
	for _, loc := range []Location{from, to} {
		if _, ok := m.Stations[loc]; !ok {
			return 0, fmt.Errorf("駅 %s が路線図にありません", name(loc))
		}
	}

	dist := map[Location]uint64{from: 0}
	done := map[Location]bool{}
	for {
		// 未確定の駅のうち、距離が最小のものを選ぶ（駅数が少ないので線形探索）
		var cur Location
		found := false
		for loc, d := range dist {
			if done[loc] {
				continue
			}
			if !found || d < dist[cur] {
				cur, found = loc, true
			}
		}
		if !found {
			return 0, fmt.Errorf("駅 %s から %s へは路線でたどり着けません", name(from), name(to))
		}
		if cur == to {
			return dist[cur], nil
		}
		done[cur] = true
		for _, line := range m.Lines {
			var next Location
			switch cur {
			case line.From:
				next = line.To
			case line.To:
				next = line.From
			default:
				continue
			}
			if nd := dist[cur] + line.Minutes; !done[next] {
				if d, ok := dist[next]; !ok || nd < d {
					dist[next] = nd
				}
			}
		}
	}
}

// Validate は路線図の整合性を確かめる。
//   - Width / Height が正
//   - Departure（自宅）が Stations にある
//   - Lines の両端が Stations にあり、From != To
//   - 同じ2駅を結ぶ路線がない（A-B と B-A も重複とみなす）
//   - 同じ座標の駅がない
//   - 全駅が Departure から Lines を辿って到達できる
func (m StageMap) Validate() error {
	if m.Width <= 0 || m.Height <= 0 {
		return fmt.Errorf("路線図の基準サイズが不正です (Width=%v, Height=%v)", m.Width, m.Height)
	}
	if _, ok := m.Stations[Departure]; !ok {
		return fmt.Errorf("自宅駅 (Departure) が定義されていません")
	}
	seen := map[[2]Location]bool{}
	for _, line := range m.Lines {
		if line.From == line.To {
			return fmt.Errorf("路線 %s - %s の両端が同じ駅です", name(line.From), name(line.To))
		}
		for _, end := range []Location{line.From, line.To} {
			if _, ok := m.Stations[end]; !ok {
				return fmt.Errorf("路線 %s - %s の駅 %s が定義されていません", name(line.From), name(line.To), name(end))
			}
		}
		// 路線は無向なので、A-B と B-A は同じ路線として数える
		pair := [2]Location{min(line.From, line.To), max(line.From, line.To)}
		if seen[pair] {
			return fmt.Errorf("路線 %s - %s が重複しています", name(line.From), name(line.To))
		}
		seen[pair] = true
	}

	locations := slices.Sorted(maps.Keys(m.Stations))
	for i, a := range locations {
		for _, b := range locations[i+1:] {
			if m.Stations[a].Pos == m.Stations[b].Pos {
				return fmt.Errorf("駅 %s と %s の座標が同じです", name(a), name(b))
			}
		}
	}

	// 自宅から幅優先で到達できる駅を数える
	reached := map[Location]bool{Departure: true}
	queue := []Location{Departure}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range m.Neighbors(cur) {
			if !reached[next] {
				reached[next] = true
				queue = append(queue, next)
			}
		}
	}
	for _, loc := range locations {
		if !reached[loc] {
			return fmt.Errorf("駅 %s は自宅から路線でたどり着けません", name(loc))
		}
	}
	return nil
}

// NewSelector は路線図に沿ってカーソルを動かす選択グラフを作る。
// 上下左右の移動先は、Line の FromKey / ToKey で明示されたものを優先し、
// 明示のない向きは駅の位置関係から自動で決める。
// どのキーでも辿れない路線がある場合はエラーを返す。
func (m StageMap) NewSelector(initial Location) (*selection.Graph[Location, selection.Direction], error) {
	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("路線図が不正です。\n%w", err)
	}
	if _, ok := m.Stations[initial]; !ok {
		return nil, fmt.Errorf("初期位置 %s が路線図にありません", name(initial))
	}

	g := selection.NewGraph[Location, selection.Direction](initial)
	// 方向キーで辿れる (from, to) の組
	reachable := map[[2]Location]bool{}

	// 明示されたキー。explicitTo は (from, キー) → 行き先、explicitPair は明示のある (from, to) の組。
	type fromKey struct {
		from Location
		key  selection.Direction
	}
	explicitTo := map[fromKey]Location{}
	explicitPair := map[[2]Location]bool{}
	addExplicit := func(from, to Location, key *selection.Direction) error {
		if key == nil {
			return nil
		}
		if _, ok := directionVectors[*key]; !ok {
			return fmt.Errorf("路線 %s → %s のキーが不正です (%d)", name(from), name(to), *key)
		}
		k := fromKey{from, *key}
		if prev, dup := explicitTo[k]; dup && prev != to {
			return fmt.Errorf("駅 %s の同じキーに %s と %s が割り当てられています", name(from), name(prev), name(to))
		}
		explicitTo[k] = to
		explicitPair[[2]Location{from, to}] = true
		g.AddTransition(from, *key, to, nil)
		reachable[[2]Location{from, to}] = true
		return nil
	}
	for _, line := range m.Lines {
		if err := addExplicit(line.From, line.To, line.FromKey); err != nil {
			return nil, err
		}
		if err := addExplicit(line.To, line.From, line.ToKey); err != nil {
			return nil, err
		}
	}

	for _, from := range slices.Sorted(maps.Keys(m.Stations)) {
		fromPos := m.Stations[from].Pos.PtoV()
		for _, d := range selection.Directions() {
			if _, taken := explicitTo[fromKey{from, d}]; taken {
				continue // このキーは明示で使われている
			}
			dirVec := directionVectors[d]
			best := -1.0
			var bestTo Location
			found := false
			for _, to := range m.Neighbors(from) {
				if explicitPair[[2]Location{from, to}] {
					continue // この向きはキーが明示されているので自動では決めない
				}
				unit, err := m.Stations[to].Pos.PtoV().Sub(fromPos).Normalize()
				if err != nil {
					return nil, fmt.Errorf("駅 %s と %s の位置が同じです\n%w", name(from), name(to), err)
				}
				dot := unit.InnerProd(dirVec)
				if dot >= directionThreshold && dot > best {
					best = dot
					bestTo = to
					found = true
				}
			}
			if found {
				g.AddTransition(from, d, bestTo, nil)
				reachable[[2]Location{from, bestTo}] = true
			}
		}
	}

	for _, line := range m.Lines {
		if !reachable[[2]Location{line.From, line.To}] || !reachable[[2]Location{line.To, line.From}] {
			return nil, fmt.Errorf("路線 %s - %s はどの方向キーでも選べません。駅の位置を調整してください。", name(line.From), name(line.To))
		}
	}
	return g, nil
}
