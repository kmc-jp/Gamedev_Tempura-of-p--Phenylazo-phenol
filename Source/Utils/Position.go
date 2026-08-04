package Utils

type Position struct {
	X float64
	Y float64
}

func (p Position) Add(v Vec2) Position {
	return p.PtoV().Add(v).VtoP()
}

func (p Position) Offset(p2 Position) Vec2 {
	return p.PtoV().Sub(p2.PtoV())
}

func (p Position) DistanceSquared(p2 Position) float64 {
	return p.Offset(p2).LengthSquared()
}

func (p Position) Distance(p2 Position) float64 {
	return p.Offset(p2).Length()
}
