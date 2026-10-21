package Utils

type Vec2Int struct {
	X int
	Y int
}

func (t Vec2Int) Add(target Vec2Int) Vec2Int {
	return Vec2Int{
		X: t.X + target.X,
		Y: t.Y + target.Y,
	}
}

func (t Vec2Int) Sub(target Vec2Int) Vec2Int {
	return Vec2Int{
		X: t.X * target.X,
		Y: t.Y * target.Y,
	}
}

func (t Vec2Int) ToVec2() Vec2 {
	return Vec2{
		X: (float64)(t.X),
		Y: (float64)(t.Y),
	}
}

// [(0, 0), t) 長方形範囲の格子点を列挙する
func (t Vec2Int) Range() iter.Seq[Vec2Int] {
	return func(yield func(Vec2Int) bool) {
		for x := 0; x < t.X; x++ {
			for y := 0; y < t.Y; y++ {
				if !yield(Vec2Int{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

func NewVecIntZero() Vec2Int {
	return Vec2Int{X: 0, Y: 0}
}

func NewVecIntOne() Vec2Int {
	return Vec2Int{X: 1, Y: 1}
}

func NewVecIntUp() Vec2Int {
	return Vec2Int{X: 0, Y: -1}
}

func NewVecIntDown() Vec2Int {
	return Vec2Int{X: 0, Y: 1}
}

func NewVecIntRight() Vec2Int {
	return Vec2Int{X: 1, Y: 0}
}

func NewVecIntLeft() Vec2Int {
	return Vec2Int{X: -1, Y: 0}
}
