package Utils

import "math"

type RectCollider struct {
	Width  float64
	Height float64
}

func (rc RectCollider) Collision(pos Position, targetrc RectCollider, targetpos Position) bool {
	return targetpos.X < pos.X+rc.Width &&
		pos.X < targetpos.X+targetrc.Width &&
		targetpos.Y < pos.Y+rc.Height &&
		pos.Y < targetpos.Y+targetrc.Height
}

func (rc RectCollider) Repulsion(pos *Position, targetrc RectCollider, targetpos Position) {
	if !rc.Collision(*pos, targetrc, targetpos) {
		return
	}

	overlapLeft := (pos.X + rc.Width) - targetpos.X
	overlapRight := (targetpos.X + targetrc.Width) - pos.X
	overlapTop := (pos.Y + rc.Height) - targetpos.Y
	overlapBottom := (targetpos.Y + targetrc.Height) - pos.Y

	minOverlapX := math.Min(overlapLeft, overlapRight)
	minOverlapY := math.Min(overlapTop, overlapBottom)

	if minOverlapX < minOverlapY {
		if overlapLeft < overlapRight {
			pos.X -= overlapLeft
		} else {
			pos.X += overlapRight
		}
	} else {
		if overlapTop < overlapBottom {
			pos.Y -= overlapTop
		} else {
			pos.Y += overlapBottom
		}
	}
}
