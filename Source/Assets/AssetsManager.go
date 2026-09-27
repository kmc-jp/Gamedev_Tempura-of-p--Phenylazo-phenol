package assets

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	"embed"
	"fmt"
	"io/fs"
)

type AssetsManager struct {
	Image     *imageassets.ImageManager
	Serialize *serialize.SerializeManager
	assetsFS  *embed.FS
}

func NewAssetsManager(assetsFS *embed.FS) (*AssetsManager, error) {
	imageFS, err := fs.Sub(assetsFS, "Assets/Images")
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの切り出しに失敗しました。\n%w", err)
	}
	im, err := imageassets.NewImageManager(imageFS)
	if err != nil {
		return nil, fmt.Errorf("NewImageManager() でエラーが発生しました。\n%w", err)
	}

	serializeFS, err := fs.Sub(assetsFS, "Assets/Serialize")
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの切り出しに失敗しました。\n%w", err)
	}
	seri, err := serialize.NewSerializeManager(serializeFS)
	if err != nil {
		return nil, fmt.Errorf("NewSerializeManager() でエラーが発生しました。\n%w", err)
	}

	return &AssetsManager{
		Image:     im,
		Serialize: seri,
		assetsFS:  assetsFS,
	}, nil
}
