package tilemap

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileData struct {
	Name          string
	ImageData     imageassets.ImageData
	Layer         string
	TileSize      Utils.Vec2Int
	TileOffset    Utils.Vec2
	Interact      InteractType
	DIO           *ebiten.DrawImageOptions
	RenderFactory func(Utils.Factory[gameobject.Render]) Utils.Factory[gameobject.Render]
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

func NewTileData(
	name string,
	imgData *imageassets.ImageData,
	layer string,
	TileSize Utils.Vec2Int,
	offset Utils.Vec2,
	Interact InteractType,
	opts ...TileDataOpt,
) (*TileData, error) {
	td := &TileData{
		Name:          name,
		ImageData:     *imgData,
		Layer:         layer,
		TileSize:      TileSize,
		TileOffset:    offset,
		Interact:      Interact,
		DIO:           &ebiten.DrawImageOptions{},
		RenderFactory: nil,
	}

	for _, opt := range opts {
		opt(td)
	}

	return td, nil
}

func TileData_DIOOption(DIO *ebiten.DrawImageOptions) TileDataOpt {
	return func(td *TileData) {
		td.DIO = DIO
	}
}

func TileData_RenderOption(RenderFactory func(Utils.Factory[gameobject.Render]) Utils.Factory[gameobject.Render]) TileDataOpt {
	return func(td *TileData) {
		td.RenderFactory = RenderFactory
	}
}
