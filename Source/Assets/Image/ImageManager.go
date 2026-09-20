package imageassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
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

func (im ImageManager) Load(data ImageData) (*ebiten.Image, error) {
	i, err := im.asepriteLoader.Load(data, im.imageFS)
	if err != nil {
		return nil, fmt.Errorf("asepriteLoader.Load() でエラーが発生しました。\n%w", err)
	}
	return i, nil
}
