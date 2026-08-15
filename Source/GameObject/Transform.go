package gameobject

type Transform interface {
}

type TransformFactory func() (*Transform, error)
