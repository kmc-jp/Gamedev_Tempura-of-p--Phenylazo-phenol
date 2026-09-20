package imageassets

import (
	"fmt"
	"io/fs"
)

type ImageManager struct {
	asepriteLoader *AsepriteLoader
	imageFS        fs.FS
}

func NewImageManager(imageFS fs.FS) (*ImageManager, error) {
	al, err := NewAsepriteLoader()
	if err != nil {
		return &ImageManager{}, fmt.Errorf("NewAsepriteLoader() でエラーが発生しました。\n%w", err)
	}
	return &ImageManager{asepriteLoader: al, imageFS: imageFS}, nil
}
