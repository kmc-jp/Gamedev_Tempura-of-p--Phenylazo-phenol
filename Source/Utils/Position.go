package Utils

import "math"

type Position struct {
	X float64
	Y float64
}

func (p Position) Add(p2 Position) Position {
	return Position{
		X: p.X + p2.X,
		Y: p.Y + p2.Y,
	}
}

func (p Position) Sub(p2 Position) Position {
	return Position{
		X: p.X - p2.X,
		Y: p.Y - p2.Y,
	}
}

func (p Position) LengthSquared() float64 {
	return p.X*p.X + p.Y*p.Y
}

func (p Position) Length() float64 {
	return math.Pow(p.LengthSquared(), 0.5)
}

func (p Position) DistanceSquared(p2 Position) float64 {
	return p.Sub(p2).LengthSquared()
}

func (p Position) Distance(p2 Position) float64 {
	return p.Sub(p2).Length()
}
