package stageinfo

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Selection"
	"testing"
)

// testMap はテスト用の路線図を返す（本番の Assets/Serialize/StageMap.toml とは別の配置）。
func testMap() StageMap {
	st := func(label string, x, y float64) Station {
		return Station{Label: label, Pos: Utils.Position{X: x, Y: y}}
	}
	return StageMap{
		Width:  320,
		Height: 240,
		Stations: map[Location]Station{
			Departure: st("自宅", 160, 150),
			Red:       st("Red", 80, 150),
			Green:     st("Green", 240, 150),
			Blue:      st("Blue", 80, 95),
			Orange:    st("Orange", 240, 95),
			Purple:    st("Purple", 160, 70),
		},
		Lines: []Line{
			{From: Departure, To: Red},
			{From: Departure, To: Green},
			{From: Departure, To: Purple},
			{From: Red, To: Blue},
			{From: Green, To: Orange},
			{From: Blue, To: Purple},
			{From: Orange, To: Purple},
		},
	}
}

func TestNewSelectorDirections(t *testing.T) {
	g, err := testMap().NewSelector(Departure)
	if err != nil {
		t.Fatalf("NewSelector でエラー: %v", err)
	}

	cases := []struct {
		from Location
		d    selection.Direction
		want Location
	}{
		{Departure, selection.Left, Red},
		{Departure, selection.Right, Green},
		{Departure, selection.Up, Purple},
		{Red, selection.Right, Departure},
		{Purple, selection.Left, Blue},
		{Blue, selection.Down, Red},
	}
	for _, c := range cases {
		g.SetCurrent(c.from)
		got, ok := g.Next(c.d)
		if !ok || got != c.want {
			t.Errorf("%s から %v: got %s (%v), want %s", name(c.from), c.d, name(got), ok, name(c.want))
		}
	}
}

func TestNewSelectorNoTransition(t *testing.T) {
	g, err := testMap().NewSelector(Departure)
	if err != nil {
		t.Fatalf("NewSelector でエラー: %v", err)
	}
	if g.Fire(selection.Down) {
		t.Fatal("自宅で Down を押して遷移した")
	}
}

func TestValidateIsolatedStation(t *testing.T) {
	m := testMap()
	// Orange につながる路線（Green - Orange、Orange - Purple）をすべて外して孤立させる
	var lines []Line
	for _, line := range m.Lines {
		if line.From != Orange && line.To != Orange {
			lines = append(lines, line)
		}
	}
	m.Lines = lines
	if err := m.Validate(); err == nil {
		t.Fatal("孤立駅があるのに Validate が成功した")
	}
}

func TestValidateUnknownStation(t *testing.T) {
	m := testMap()
	delete(m.Stations, Orange)
	if err := m.Validate(); err == nil {
		t.Fatal("未定義の駅への路線があるのに Validate が成功した")
	}
}

func TestValidateDuplicateLine(t *testing.T) {
	// 同じ向き
	m := testMap()
	m.Lines = append(m.Lines, Line{From: Departure, To: Red})
	if err := m.Validate(); err == nil {
		t.Fatal("同じ路線が2本あるのに Validate が成功した")
	}
	// 逆向き
	m = testMap()
	m.Lines = append(m.Lines, Line{From: Red, To: Departure})
	if err := m.Validate(); err == nil {
		t.Fatal("A-B と B-A があるのに Validate が成功した")
	}
}

func TestNewSelectorUnreachableLine(t *testing.T) {
	st := func(x, y float64) Station { return Station{Label: "x", Pos: Utils.Position{X: x, Y: y}} }
	// Departure の真上に Blue と Red が並ぶ。Departure から見て Up は Blue が取るので、
	// Red への路線は Departure からどのキーでも選べない。
	m := StageMap{
		Width:  100,
		Height: 100,
		Stations: map[Location]Station{
			Departure: st(50, 90),
			Blue:      st(50, 50),
			Red:       st(52, 20),
		},
		Lines: []Line{
			{From: Departure, To: Blue},
			{From: Departure, To: Red},
		},
	}
	if _, err := m.NewSelector(Departure); err == nil {
		t.Fatal("選べない路線があるのに NewSelector が成功した")
	}
}

func TestNewSelectorExplicitKey(t *testing.T) {
	m := testMap()
	// Departure - Red を、自動なら Left のところを、どちらの駅でも Down で行き来するように明示する。
	down := selection.Down
	m.Lines[0].FromKey = &down
	m.Lines[0].ToKey = &down
	g, err := m.NewSelector(Departure)
	if err != nil {
		t.Fatalf("NewSelector でエラー: %v", err)
	}
	if got, ok := g.Next(selection.Down); !ok || got != Red {
		t.Fatalf("自宅で Down を押すと Red へ行くはずが %v, %v", got, ok)
	}
	if _, ok := g.Next(selection.Left); ok {
		t.Fatal("明示した路線は、自動のキー (Left) では選べないはず")
	}
	g.SetCurrent(Red)
	if got, ok := g.Next(selection.Down); !ok || got != Departure {
		t.Fatalf("Red で Down を押すと自宅へ戻るはずが %v, %v", got, ok)
	}
	// 明示していない路線は自動のまま
	g.SetCurrent(Departure)
	if got, ok := g.Next(selection.Right); !ok || got != Green {
		t.Fatalf("自宅で Right は Green のはずが %v, %v", got, ok)
	}
}

func TestNewSelectorExplicitOneWay(t *testing.T) {
	m := testMap()
	// FromKey だけ明示（Departure → Red を Down）。Red → Departure は自動（Right）のまま。
	down := selection.Down
	m.Lines[0].FromKey = &down
	g, err := m.NewSelector(Red)
	if err != nil {
		t.Fatalf("NewSelector でエラー: %v", err)
	}
	if got, ok := g.Next(selection.Right); !ok || got != Departure {
		t.Fatalf("Red で Right は自宅のはずが %v, %v", got, ok)
	}
}

func TestNewSelectorExplicitConflict(t *testing.T) {
	m := testMap()
	// Departure の Down に Red と Green を割り当てる。
	down := selection.Down
	m.Lines[0].FromKey = &down
	m.Lines[1].FromKey = &down
	if _, err := m.NewSelector(Departure); err == nil {
		t.Fatal("同じキーに別の行き先を明示したのに NewSelector が成功した")
	}
}
