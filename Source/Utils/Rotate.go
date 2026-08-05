package Utils

import "math"

type Rotation struct {
	Angle float64
}

// Rotate in radian
func (rot *Rotation) Rotate(rad float64) {
	rot.Angle += rad
	rot.normalize()
}

// Rotate using degrees
func (rot *Rotation) RotateDegree(deg float64) {
	rot.Angle += deg * math.Pi / 180
	rot.normalize()
}

func (rot *Rotation) normalize() {
	rad := rot.Angle
	for rad > math.Pi {
		rad -= 2 * math.Pi
	}
	for rad < -math.Pi {
		rad += 2 * math.Pi
	}
	rot.Angle = rad
}

// Return unit direction vector
func (rot *Rotation) Forward() Vec2 {
	return Vec2{
		X: math.Cos(rot.Angle),
		Y: math.Sin(rot.Angle),
	}
}

func (rot *Rotation) LookAt(p, target Position) {
	dx := target.X - p.X
	dy := target.Y - p.Y
	rot.Angle = math.Atan2(dy, dx)
}
