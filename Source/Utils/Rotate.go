package Utils

import "math"

type Rotation struct {
	Angle float64
}

// Rotate in radian
func (rot Rotation) Rotate(rad float64) {
	rot.Angle += rad
}

// Rotate using degrees
func (rot Rotation) RotateDegree(deg float64) {
	rot.Angle += deg * math.Pi / 180
}
