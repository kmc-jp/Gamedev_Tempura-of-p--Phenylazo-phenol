package Utils

import "math"

type Rotation struct {
	Angle float64
}

// Rotate in radian
func (rot Rotation) Rotate(theta float64) {
	rot.Angle += theta
}

// Rotate using degrees
func (rot Rotation) RotateDegree(theta float64) {
	rot.Angle += theta * math.Pi / 180
}
