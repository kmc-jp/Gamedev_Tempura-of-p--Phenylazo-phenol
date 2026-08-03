package Utils

type Rotation struct {
	angle float64
}

func (rot Rotation) Rotate(theta float64) {
	rot.angle += theta
}
