package imageassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/goaseprite"
)

type AsepriteLoader struct {
	cacheAseprite map[string]*goaseprite.File
	cacheImage    map[string]*ebiten.Image
}

func NewAsepriteLoader() (*AsepriteLoader, error) {
	al := AsepriteLoader{
		cacheAseprite: map[string]*goaseprite.File{},
		cacheImage:    map[string]*ebiten.Image{},
	}
	return &al, nil
}

func (al AsepriteLoader) Load(data ImageData, fs fs.FS) (*goaseprite.File, *ebiten.Image, error) {
	asp, aspOk := al.cacheAseprite[data.Name]
	img, imgOk := al.cacheImage[data.Name]
	if !aspOk || !imgOk {
		asp, img, err := al.fetch(data, fs)
		if err != nil {
			return nil, nil, fmt.Errorf("AsepriteLoader.fetch() でエラーが発生しました。\n%w", err)
		}
		al.cacheAseprite[data.Name] = asp
		al.cacheImage[data.Name] = img
	}
	asp, _ = al.cacheAseprite[data.Name]
	img, _ = al.cacheImage[data.Name]
	return asp, img, nil
}

func (al AsepriteLoader) fetch(data ImageData, fs fs.FS) (*goaseprite.File, *ebiten.Image, error) {
	aspFile, err := goaseprite.Open(data.JsonPath, fs)
	if err != nil {
		return nil, nil, fmt.Errorf("goaseprite.Open(data.JsonPath) でエラーが発生しました。\n%w", err)
	}
	img, _, err := ebitenutil.NewImageFromFileSystem(fs, aspFile.ImagePath)
	if err != nil {
		return nil, nil, fmt.Errorf("ebitenutil.NewImageFromFile(aspFile.ImagePath) でエラーが発生しました。\n%w", err)
	}
	return aspFile, img, nil
}
