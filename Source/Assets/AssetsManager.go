package assets

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	"embed"
	"fmt"
)

type AssetsManager struct {
	Image    *imageassets.ImageManager
	assetsFS *embed.FS
}

func NewAssetsManager(assetsFS *embed.FS) (*AssetsManager, error) {
	im, err := imageassets.NewImageManager()
	if err != nil {
		return &AssetsManager{}, fmt.Errorf("ImageManager() でエラーが発生しました。\n%w", err)
	}
	return &AssetsManager{
		Image:    im,
		assetsFS: assetsFS,
	}, nil
}
