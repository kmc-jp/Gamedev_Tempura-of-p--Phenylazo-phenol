package imageassets

import (
	"mime"

	"github.com/hajimehoshi/ebiten/v2"
)

type AsepriteLoader struct {
	cache map[string]ebiten.Image
}

func NewAsepriteLoader() (*AsepriteLoader, error) {
	// mime のシステムに .aseprite を追加
	mime.AddExtensionType(".aseprite", "image/aseprite")
	return &AsepriteLoader{}, nil
}
