package assetsutil

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func CheckAndCleanPath(path string, FS fs.FS) (string, error) {
	if len(strings.TrimSpace(path)) == 0 {
		return path, fmt.Errorf("path が空です。")
	}

	cleanedpath := filepath.ToSlash(filepath.Clean(path))
	if cleanedpath == "." || !fs.ValidPath(cleanedpath) {
		return path, fmt.Errorf("path \"%s\" が不正です。", path)
	}
	cleanedpath = strings.TrimPrefix(cleanedpath, "/")

	info, err := fs.Stat(FS, cleanedpath)
	if err != nil {
		return cleanedpath, fmt.Errorf("path \"%s\" が使用できません。\n%w", path, err)
	}
	if info.IsDir() {
		return cleanedpath, fmt.Errorf("path \"%s\" はディレクトリです。", path)
	}

	return cleanedpath, nil
}
