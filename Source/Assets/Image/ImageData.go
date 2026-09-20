package imageassets

import (
	"fmt"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"
)

type ImageData struct {
	Name string
	// /Assets/Image からの相対パス
	Path string
	// 拡張子タイプ
	mimeType string
}

func NewImageData(name string, path string, ImageDirectory fs.FS) (*ImageData, error) {
	if len(strings.TrimSpace(path)) == 0 {
		return nil, fmt.Errorf("path が空です。")
	}

	cleaned := filepath.Clean(filepath.Join("/", path))

	if cleaned == "." {
		return nil, fmt.Errorf("path が不正です。")
	}

	info, err := fs.Stat(ImageDirectory, cleaned)
	if err != nil {
		return nil, fmt.Errorf("path が使用できません。")
	}
	if info.IsDir() {
		return nil, fmt.Errorf("path はディレクトリです。")
	}

	ext := filepath.Ext(cleaned)
	mimeType := mime.TypeByExtension(ext)
	id := ImageData{
		Name:     name,
		Path:     cleaned,
		mimeType: mimeType,
	}
	return &id, nil
}
