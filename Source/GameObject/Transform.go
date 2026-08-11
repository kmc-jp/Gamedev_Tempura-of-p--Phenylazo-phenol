package gameobject

type Transform struct {
}

func NewTransform() (*Transform, error) {
	return &Transform{}, nil
}
