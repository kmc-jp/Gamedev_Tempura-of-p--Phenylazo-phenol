package imageassets

import "github.com/hajimehoshi/ebiten/v2"

type AsepriteLoader struct {
	cache map[string]ebiten.Image
}

func NewAsepriteLoader() (*AsepriteLoader, error) {
	return &AsepriteLoader{}, nil
}
