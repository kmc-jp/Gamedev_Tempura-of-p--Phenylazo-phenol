package audioassets

import (
	assetsutil "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/AssetsUtils"
	"fmt"
	"path/filepath"
	"strings"
)

type AudioType int

const (
	BGM AudioType = iota
	SE
)

type AudioData struct {
	// /Assets/Audio からの相対パス
	Path       string
	ext        string
	SampleRate int
	Type       AudioType
}

func (ad AudioData) Extention() string {
	return ad.ext
}

func NewImageData(path string, SampleRate int, Type AudioType, audioManager *AudioManager) (*AudioData, error) {
	cleaned, err := assetsutil.CheckAndCleanPath(path, audioManager.AudioFS)
	if err != nil {
		return nil, fmt.Errorf("path が不正です。\n%w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	id := AudioData{
		Path:       cleaned,
		ext:        ext,
		Type:       Type,
		SampleRate: SampleRate,
	}
	return &id, nil
}
