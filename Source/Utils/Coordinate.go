package Utils

import (
	"errors"
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

func (v1 Vec2) VtoP() Position {
	return Position{
		X: v1.X,
		Y: v1.Y,
	}
}

func (p Position) PtoV() Vec2 {
	return Vec2{
		X: p.X,
		Y: p.Y,
	}
}

func NewZeroVec2() Vec2 {
	return Vec2{
		X: 0,
		Y: 0,
	}
}

func NewLeftVec2() Vec2 {
	return Vec2{
		X: -1,
		Y: 0,
	}
}

func NewRightVec2() Vec2 {
	return Vec2{
		X: 1,
		Y: 0,
	}
}

func NewUpVec2() Vec2 {
	return Vec2{
		X: 0,
		Y: -1,
	}
}

func NewDownVec2() Vec2 {
	return Vec2{
		X: 0,
		Y: 1,
	}
}

func NewOneVec2() Vec2 {
	return Vec2{
		X: 1,
		Y: 1,
	}
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X + v2.X, Y: v1.Y + v2.Y,
	}
}

func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X - v2.X, Y: v1.Y - v2.Y,
	}
}

func (v1 Vec2) Minus() Vec2 {
	return Vec2{
		X: -v1.X, Y: -v1.Y,
	}
}

func (v1 Vec2) Length() float64 {
	return math.Pow(v1.LengthSquared(), 0.5)
}

func (v1 Vec2) LengthSquared() float64 {
	return v1.X*v1.X + v1.Y*v1.Y
}

func (v1 Vec2) Normalize() (Vec2, error) {
	if v1.LengthSquared() == 0 {
		return NewZeroVec2(), errors.New("零ベクトルは正規化できません")
	} else {
		return v1.ScalarProd(1 / v1.Length()), nil
	}
}

const atol float64 = 0.0000000000000001 // e-(7*2)
const rtol float64 = 0.0000000001       // e-(5*2)

func (v1 Vec2) Equal(v2 Vec2) bool {
	return v1.Sub(v2).LengthSquared() <= atol+rtol*v2.LengthSquared()
}

func (v1 Vec2) ScalarProd(n float64) Vec2 {
	return Vec2{
		X: v1.X * n, Y: v1.Y * n,
	}
}

func ScalarProd(v1 Vec2, n float64) Vec2 {
	return Vec2{
		X: v1.X * n, Y: v1.Y * n,
	}
}

func (v1 Vec2) InnerProd(v2 Vec2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func InnerProd(v1, v2 Vec2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v1 Vec2) CrossProd(v2 Vec2) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

func CrossProd(v1, v2 Vec2) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

func (v1 Vec2) HadamardProd(v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X * v2.X, Y: v1.Y * v2.Y,
	}
}

func HadamardProd(v1, v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X * v2.X, Y: v1.Y * v2.Y,
	}
}
