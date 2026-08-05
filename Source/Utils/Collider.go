package Utils

type RectCollider struct {
	Width  float64
	Height float64
}

// Detect the collision of two objects.
func (rc RectCollider) Collision(pos Position, targetrc RectCollider, targetpos Position) bool {
	left := targetpos.X <= pos.X+rc.Width
	right := pos.X <= targetpos.X+targetrc.Width
	top := targetpos.Y <= pos.Y+rc.Height
	bottom := pos.Y <= targetpos.Y+targetrc.Height
	return left && right && top && bottom
}

func Repulsion(rc RectCollider, pos Position, targetrc RectCollider, targetpos Position) {
	for rc.Collision(pos, targetrc, targetpos) {
		center := pos.PtoV().Add(Vec2{rc.Width / 2, rc.Height / 2})
		targetCenter := targetpos.PtoV().Add(Vec2{targetrc.Width / 2, targetrc.Height / 2})
		direction := center.Sub(targetCenter)
		direction, _ = direction.Normalize()
		pos.Move(direction)
	}
}
