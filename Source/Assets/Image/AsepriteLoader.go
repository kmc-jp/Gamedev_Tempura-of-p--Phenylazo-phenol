package imageassets

import (
	"fmt"
	"io/fs"
	"mime"

	"github.com/hajimehoshi/ebiten/v2"
)

type AsepriteLoader struct {
	cache map[string]*ebiten.Image
}

func NewAsepriteLoader() (*AsepriteLoader, error) {
	// mime のシステムに .aseprite を追加
	mime.AddExtensionType(".aseprite", "image/aseprite")
	return &AsepriteLoader{}, nil
}

func (al AsepriteLoader) Load(data ImageData, fs fs.FS) (*ebiten.Image, error) {
	image, ok := al.cache[data.Name]
	if !ok {
		image, err := al.fetch(data.Path, fs)
		if err != nil {
			return nil, fmt.Errorf("AsepriteLoader.fetch() でエラーが発生しました。\n%w", err)
		}
		al.cache[data.Name] = image
	}
	return image, nil
}

func (al AsepriteLoader) fetch(path string, fs fs.FS) (*ebiten.Image, error) {
	panic("Not Implimented!")
}
