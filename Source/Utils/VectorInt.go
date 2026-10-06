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

func NewVecIntZero() Vec2Int {
	return Vec2Int{X: 0, Y: 0}
}

func NewVecIntUp() Vec2Int {
	return Vec2Int{X: 0, Y: 1}
}

func NewVecIntDown() Vec2Int {
	return Vec2Int{X: 0, Y: -1}
}

func NewVecIntRight() Vec2Int {
	return Vec2Int{X: 1, Y: 0}
}

func NewVecIntLeft() Vec2Int {
	return Vec2Int{X: -1, Y: 0}
}
