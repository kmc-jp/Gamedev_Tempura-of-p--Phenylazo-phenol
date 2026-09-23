package imageassets

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

type ImageManager struct {
	asepriteLoader *AsepriteLoader
	ImageFS        fs.FS
}

func NewImageManager(imageFS fs.FS) (*ImageManager, error) {
	al, err := NewAsepriteLoader()
	if err != nil {
		return &ImageManager{}, fmt.Errorf("NewAsepriteLoader() でエラーが発生しました。\n%w", err)
	}
	return &ImageManager{asepriteLoader: al, ImageFS: imageFS}, nil
}

func (im ImageManager) Load(data ImageData) (*goaseprite.File, *ebiten.Image, error) {
	asp, img, err := im.asepriteLoader.Load(data, im.ImageFS)
	if err != nil {
		return nil, img, fmt.Errorf("asepriteLoader.Load() でエラーが発生しました。\n%w", err)
	}
	return asp, img, nil
}

func (im ImageManager) CheckAndCleanPath(path string) (string, error) {
	if len(strings.TrimSpace(path)) == 0 {
		return path, fmt.Errorf("path が空です。")
	}

	cleanedpath := filepath.Clean(path)
	if cleanedpath == "." {
		return path, fmt.Errorf("jsonpath が不正です。")
	}
	cleanedpath = strings.TrimPrefix(cleanedpath, "/")

	info, err := fs.Stat(im.ImageFS, cleanedpath)
	if err != nil {
		return cleanedpath, fmt.Errorf("pngpath が使用できません。\n%w", err)
	}
	if info.IsDir() {
		return cleanedpath, fmt.Errorf("pngpath はディレクトリです。")
	}

	return cleanedpath, nil
}
