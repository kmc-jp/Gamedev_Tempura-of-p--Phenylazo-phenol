package tilemap

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileData struct {
	Name      string
	ImageData imageassets.ImageData
	Layer     string
	TileSize  Utils.Vec2Int
	Interact  InteractType
	DIO       ebiten.DrawImageOptions
}

type InteractType int

const (
	None InteractType = iota
	Shelf
)

func IsIntaractive(t InteractType) bool {
	return t != None
}

type TileDataOpt func(*TileData)

func NewTileData(name string, imgData *imageassets.ImageData, layer string, TileSize Utils.Vec2Int, Interact InteractType, opts ...TileDataOpt) (*TileData, error) {
	td := &TileData{
		Name:      name,
		ImageData: *imgData,
		Layer:     layer,
		TileSize:  TileSize,
		Interact:  Interact,
		DIO:       ebiten.DrawImageOptions{},
	}

	for _, opt := range opts {
		opt(td)
	}

	return td, nil
}

func TileData_DIOOption(DIO ebiten.DrawImageOptions) TileDataOpt {
	return func(td *TileData) {
		td.DIO = DIO
	}
}
