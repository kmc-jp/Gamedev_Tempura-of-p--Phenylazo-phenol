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
	// /Assets/Sounds からの相対パス
	Path string
	ext  string
	Type AudioType
}

func (ad AudioData) Extension() string {
	return ad.ext
}

func (am AudioManager) NewAudioData(path string, Type AudioType) (AudioData, error) {
	cleaned, err := assetsutil.CheckAndCleanPath(path, am.AudioFS)
	if err != nil {
		return AudioData{}, fmt.Errorf("path が不正です。\n%w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	ad := AudioData{
		Path: cleaned,
		ext:  ext,
		Type: Type,
	}
	return ad, nil
}
