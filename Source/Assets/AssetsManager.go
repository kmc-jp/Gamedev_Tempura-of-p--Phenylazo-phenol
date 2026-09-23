package assets

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	"embed"
	"fmt"
	"io/fs"
)

type AssetsManager struct {
	Image    *imageassets.ImageManager
	assetsFS *embed.FS
}

func NewAssetsManager(assetsFS *embed.FS) (*AssetsManager, error) {
	imageFS, err := fs.Sub(assetsFS, "Assets/Images")
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの切り出しに失敗しました。\n%w", err)
	}
	im, err := imageassets.NewImageManager(imageFS)
	if err != nil {
		return nil, fmt.Errorf("ImageManager() でエラーが発生しました。\n%w", err)
	}
	return &AssetsManager{
		Image:    im,
		assetsFS: assetsFS,
	}, nil
}
