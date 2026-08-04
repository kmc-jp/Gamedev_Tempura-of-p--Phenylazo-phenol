package Utils

type Vec2 struct {
	X float64
	Y float64
}

func (v Vec2) VtoP() Position {
	return Position{
		X: v.X,
		Y: v.Y,
	}
}

func (p Position) PtoV() Vec2 {
	return Vec2{
		X: p.X,
		Y: p.Y,
	}
}
