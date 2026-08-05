package Utils

type RectCollider struct {
	width  float64
	height float64
}

func Collision(rc RectCollider, pos Position, targetrc RectCollider, targetpos Position) bool {
	return targetpos.X <= pos.X+rc.width && pos.X <= targetpos.X+targetrc.width && targetpos.Y <= pos.Y+rc.height && pos.Y <= targetpos.Y+targetrc.height
}
