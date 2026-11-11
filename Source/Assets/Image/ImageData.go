package imageassets

import (
	assetsutil "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/AssetsUtils"
	"fmt"
)

type ImageData struct {
	Name string
	// /Assets/Image からの相対パス
	PngPath  string
	JsonPath string
}

func NewImageData(name string, pngpath string, jsonpath string, imageManager *ImageManager) (*ImageData, error) {
	pngcleaned, err := assetsutil.CheckAndCleanPath(pngpath, imageManager.imageFS)
	if err != nil {
		return nil, fmt.Errorf("pngpath が不正です。\n%w", err)
	}
	jsoncleaned, err := assetsutil.CheckAndCleanPath(jsonpath, imageManager.imageFS)
	if err != nil {
		return nil, fmt.Errorf("jsonpath が不正です。\n%w", err)
	}

	id := ImageData{
		Name:     name,
		PngPath:  pngcleaned,
		JsonPath: jsoncleaned,
	}
	return &id, nil
}
