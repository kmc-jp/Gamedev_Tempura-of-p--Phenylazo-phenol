package imageassets

import "fmt"

type ImageManager struct {
	asepriteLoader *AsepriteLoader
}

func NewImageManager() (*ImageManager, error) {
	al, err := NewAsepriteLoader()
	if err != nil {
		return &ImageManager{}, fmt.Errorf("NewAsepriteLoader() でエラーが発生しました。\n%w", err)
	}
	return &ImageManager{asepriteLoader: al}, nil
}
