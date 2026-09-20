package assets

import imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"

type AssetsManager struct {
	ImageManager imageassets.ImageManager
}

func NewAssetsManager() (*AssetsManager, error) {
	return &AssetsManager{}, nil
}
