package component

type StdTransform struct {
}

func NewTransform() (*StdTransform, error) {
	return &StdTransform{}, nil
}
