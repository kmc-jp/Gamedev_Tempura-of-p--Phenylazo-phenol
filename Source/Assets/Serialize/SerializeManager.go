package serialize

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type SerializeManager struct {
	tomlLoader  *TomlLoader
	SerializeFS fs.FS
}

func NewSerializeManager(serializeFS fs.FS) (*SerializeManager, error) {
	tl, err := NewTomlLoader()
	if err != nil {
		return &SerializeManager{}, fmt.Errorf("NewTomlLoader() でエラーが発生しました。\n%w", err)
	}
	return &SerializeManager{tomlLoader: tl, SerializeFS: serializeFS}, nil
}

func (im SerializeManager) Load(data SerializeData) (any, error) {
	item, err := im.tomlLoader.Load(data, im.SerializeFS)
	if err != nil {
		return nil, fmt.Errorf("tomlLoader.Load() でエラーが発生しました。\n%w", err)
	}
	return item, nil
}

func (im SerializeManager) CheckAndCleanPath(path string) (string, error) {
	if len(strings.TrimSpace(path)) == 0 {
		return path, fmt.Errorf("path が空です。")
	}

	cleanedpath := filepath.Clean(path)
	if cleanedpath == "." {
		return path, fmt.Errorf("path \"%s\" が不正です。", path)
	}
	cleanedpath = strings.TrimPrefix(cleanedpath, "/")

	info, err := fs.Stat(im.SerializeFS, cleanedpath)
	if err != nil {
		return cleanedpath, fmt.Errorf("path \"%s\" が使用できません。\n%w", path, err)
	}
	if info.IsDir() {
		return cleanedpath, fmt.Errorf("path \"%s\" はディレクトリです。", path)
	}

	return cleanedpath, nil
}
