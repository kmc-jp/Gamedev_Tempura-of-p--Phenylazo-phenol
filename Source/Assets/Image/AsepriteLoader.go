package imageassets

import (
	"fmt"
	"io/fs"

	"github.com/solarlune/goaseprite"
)

type AsepriteLoader struct {
	cache map[string]*goaseprite.File
}

func NewAsepriteLoader() (*AsepriteLoader, error) {
	return &AsepriteLoader{}, nil
}

func (al AsepriteLoader) Load(data ImageData, fs fs.FS) (*goaseprite.File, error) {
	image, ok := al.cache[data.Name]
	if !ok {
		image, err := al.fetch(data, fs)
		if err != nil {
			return nil, fmt.Errorf("AsepriteLoader.fetch() でエラーが発生しました。\n%w", err)
		}
		al.cache[data.Name] = image
	}
	return image, nil
}

func (al AsepriteLoader) fetch(data ImageData, fs fs.FS) (*goaseprite.File, error) {
	aspFile, err := goaseprite.Open(data.JsonPath, fs)
	if err != nil {
		return nil, fmt.Errorf("goaseprite.Open(data.JsonPath) でエラーが発生しました。\n%w", err)
	}
	al.cache[data.Name] = aspFile
	return aspFile, nil
}
