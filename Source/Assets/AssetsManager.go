package assets

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	"fmt"
)

type AssetsManager struct {
	Image *imageassets.ImageManager
}

func NewAssetsManager() (*AssetsManager, error) {
	im, err := imageassets.NewImageManager()
	if err != nil {
		return &AssetsManager{}, fmt.Errorf("ImageManager() でエラーが発生しました。\n%w", err)
	}
	return &AssetsManager{
		Image: im,
	}, nil
}
