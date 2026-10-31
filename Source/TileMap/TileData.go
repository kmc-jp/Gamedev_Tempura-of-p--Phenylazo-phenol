package tilemap

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TileData struct {
	Name      string
	ImageData imageassets.ImageData
	Layer     string
	TileSize  Utils.Vec2Int
	Interact  InteractType
}

type InteractType int

const (
	None InteractType = iota
	Shelf
)

func IsIntaractive(t InteractType) bool {
	return t != None
}

func NewTileData(name string, imgData *imageassets.ImageData, layer string, TileSize Utils.Vec2Int, Interact InteractType) (*TileData, error) {
	return &TileData{Name: name, ImageData: *imgData, Layer: layer, TileSize: TileSize, Interact: Interact}, nil
}
