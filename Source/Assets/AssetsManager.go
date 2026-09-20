package assets

type AssetsManager struct {
	ImageManager ImageManager
}

func NewAssetsManager() (*AssetsManager, error) {
	return &AssetsManager{}, nil
}
