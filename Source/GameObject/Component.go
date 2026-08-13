package gameobject

type Component interface {
	Update(obj GameObject, active bool) (err error)
}
