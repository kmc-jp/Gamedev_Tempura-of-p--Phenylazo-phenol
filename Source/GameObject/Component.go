package gameobject

type Component interface {
	Update(active bool) (err error)
}
