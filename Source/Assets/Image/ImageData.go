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

func (im *ImageManager) NewImageData(name string, pngpath string, jsonpath string) (*ImageData, error) {
	pngcleaned, err := assetsutil.CheckAndCleanPath(pngpath, im.imageFS)
	if err != nil {
		return nil, fmt.Errorf("pngpath が不正です。\n%w", err)
	}
	jsoncleaned, err := assetsutil.CheckAndCleanPath(jsonpath, im.imageFS)
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
