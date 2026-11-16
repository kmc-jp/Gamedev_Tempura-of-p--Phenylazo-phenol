package audioassets

import (
	assetsutil "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/AssetsUtils"
	"fmt"
	"path/filepath"
	"strings"
)

type AudioData struct {
	// /Assets/Audio からの相対パス
	Path       string
	ext        string
	SampleRate int
}

func (ad AudioData) Extention() string {
	return ad.ext
}

func NewImageData(path string, SampleRate int, audioManager *AudioManager) (*AudioData, error) {
	cleaned, err := assetsutil.CheckAndCleanPath(path, audioManager.AudioFS)
	if err != nil {
		return nil, fmt.Errorf("path が不正です。\n%w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	id := AudioData{
		Path:       cleaned,
		ext:        ext,
		SampleRate: SampleRate,
	}
	return &id, nil
}
