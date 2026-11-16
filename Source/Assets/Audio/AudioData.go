package audioassets

import (
	assetsutil "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/AssetsUtils"
	"fmt"
)

type AudioData struct {
	// /Assets/Audio からの相対パス
	Path string
}

func NewImageData(path string, audioManager *AudioManager) (*AudioData, error) {
	cleaned, err := assetsutil.CheckAndCleanPath(path, audioManager.AudioFS)
	if err != nil {
		return nil, fmt.Errorf("path が不正です。\n%w", err)
	}

	id := AudioData{
		Path: cleaned,
	}
	return &id, nil
}
