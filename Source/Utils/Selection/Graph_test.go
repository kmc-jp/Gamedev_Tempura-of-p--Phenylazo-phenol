package selection

import "testing"

func TestFireUnconditional(t *testing.T) {
	g := NewGraph[string, Direction]("A")
	g.AddTransition("A", Right, "B", nil)

	if !g.Fire(Right) {
		t.Fatal("Fire(Right) が false")
	}
	if g.Current() != "B" {
		t.Fatalf("Current() = %s, want B", g.Current())
	}
	if g.Revision() != 1 {
		t.Fatalf("Revision() = %d, want 1", g.Revision())
	}
}

func TestFireNoTransition(t *testing.T) {
	g := NewGraph[string, Direction]("A")
	g.AddTransition("A", Right, "B", nil)

	if g.Fire(Left) {
		t.Fatal("遷移がないのに Fire(Left) が true")
	}
	if g.Current() != "A" || g.Revision() != 0 {
		t.Fatalf("状態が変化した: Current=%s Revision=%d", g.Current(), g.Revision())
	}
}

func TestFireCondition(t *testing.T) {
	open := false
	g := NewGraph[string, Direction]("A")
	g.AddTransition("A", Right, "B", func() bool { return open })

	if g.Fire(Right) {
		t.Fatal("条件が false なのに遷移した")
	}
	open = true
	if !g.Fire(Right) || g.Current() != "B" {
		t.Fatal("条件が true なのに遷移しなかった")
	}
}

func TestFireRegistrationOrder(t *testing.T) {
	g := NewGraph[string, Direction]("A")
	g.AddTransition("A", Right, "C", func() bool { return false })
	g.AddTransition("A", Right, "B", nil)

	if !g.Fire(Right) || g.Current() != "B" {
		t.Fatalf("Current() = %s, want B", g.Current())
	}
}

func TestSetCurrent(t *testing.T) {
	g := NewGraph[string, Direction]("A")

	g.SetCurrent("A")
	if g.Revision() != 0 {
		t.Fatalf("同じ値で Revision が増えた: %d", g.Revision())
	}
	g.SetCurrent("B")
	if g.Current() != "B" || g.Revision() != 1 {
		t.Fatalf("Current=%s Revision=%d, want B / 1", g.Current(), g.Revision())
	}
}

func TestNextDoesNotChangeState(t *testing.T) {
	g := NewGraph[string, Direction]("A")
	g.AddTransition("A", Right, "B", nil)

	to, ok := g.Next(Right)
	if !ok || to != "B" {
		t.Fatalf("Next(Right) = %s, %v", to, ok)
	}
	if g.Current() != "A" || g.Revision() != 0 {
		t.Fatalf("Next が状態を変えた: Current=%s Revision=%d", g.Current(), g.Revision())
	}
}

func TestSetGuard(t *testing.T) {
	g := NewGraph[string, Direction]("A")
	g.AddTransition("A", Right, "B", nil)
	g.AddTransition("A", Right, "C", nil)

	// B へは行けないルール → 次の候補 C が採用される
	g.SetGuard(func(from string, trigger Direction, to string) bool { return to != "B" })
	if to, ok := g.Next(Right); !ok || to != "C" {
		t.Fatalf("guard で B を弾いたのに Next = %s, %v", to, ok)
	}

	// 全部弾くルールに差し替え
	g.SetGuard(func(from string, trigger Direction, to string) bool { return false })
	if g.Fire(Right) {
		t.Fatal("全て弾く guard なのに遷移した")
	}

	// 解除すると元に戻る
	g.SetGuard(nil)
	if to, ok := g.Next(Right); !ok || to != "B" {
		t.Fatalf("guard 解除後の Next = %s, %v", to, ok)
	}
}

func TestSetCondition(t *testing.T) {
	g := NewGraph[string, Direction]("A")
	g.AddTransition("A", Right, "B", nil)

	if !g.SetCondition("A", Right, "B", func() bool { return false }) {
		t.Fatal("登録済みの遷移の差し替えが false")
	}
	if g.Fire(Right) {
		t.Fatal("差し替えた条件が false なのに遷移した")
	}

	if !g.SetCondition("A", Right, "B", nil) {
		t.Fatal("nil への差し替えが false")
	}
	if !g.Fire(Right) {
		t.Fatal("無条件に戻したのに遷移しなかった")
	}

	if g.SetCondition("A", Left, "B", nil) {
		t.Fatal("未登録の遷移の差し替えが true")
	}
}
