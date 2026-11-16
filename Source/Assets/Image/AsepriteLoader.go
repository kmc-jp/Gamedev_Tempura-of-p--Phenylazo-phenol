package imageassets

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/goaseprite"
)

type imagecachekey struct {
	JsonPath string
	PngPath  string
}

type AsepriteLoader struct {
	cacheAseprite map[imagecachekey]*goaseprite.File
	cacheImage    map[imagecachekey]*ebiten.Image
}

func NewAsepriteLoader() (*AsepriteLoader, error) {
	al := AsepriteLoader{
		cacheAseprite: map[imagecachekey]*goaseprite.File{},
		cacheImage:    map[imagecachekey]*ebiten.Image{},
	}
	return &al, nil
}

func (al AsepriteLoader) Load(data ImageData, fs fs.FS) (*goaseprite.File, *ebiten.Image, error) {
	cachekey := imagecachekey{
		JsonPath: data.JsonPath,
		PngPath:  data.PngPath,
	}

	asp, aspOk := al.cacheAseprite[cachekey]
	img, imgOk := al.cacheImage[cachekey]
	if !aspOk || !imgOk {
		asp, img, err := al.fetch(data, fs)
		if err != nil {
			return nil, nil, fmt.Errorf("AsepriteLoader.fetch() でエラーが発生しました。\n%w", err)
		}
		al.cacheAseprite[cachekey] = asp
		al.cacheImage[cachekey] = img
	}
	asp, _ = al.cacheAseprite[cachekey]
	img, _ = al.cacheImage[cachekey]
	return asp, img, nil
}

func (al AsepriteLoader) fetch(data ImageData, fs fs.FS) (*goaseprite.File, *ebiten.Image, error) {
	aspFile, err := goaseprite.Open(data.JsonPath, fs)
	if err != nil {
		return nil, nil, fmt.Errorf("goaseprite.Open(data.JsonPath) でエラーが発生しました。\ndata.JsonPath : %s\n%w", data.JsonPath, err)
	}
	img, _, err := ebitenutil.NewImageFromFileSystem(fs, data.PngPath)
	if err != nil {
		return nil, nil, fmt.Errorf("ebitenutil.NewImageFromFile(aspFile.ImagePath) でエラーが発生しました。\naspFile.ImagePath : %s\n%w", aspFile.ImagePath, err)
	}
	return aspFile, img, nil
}
