package tilemap

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TileData struct {
	Name      string
	ImageData imageassets.ImageData
	Layer     string
	Size      Utils.Vec2Int
	Interact  InteractType
}

type InteractType int

const (
	None InteractType = iota
	Shelf
)

func NewTileData(name string, imgData *imageassets.ImageData, layer string, size Utils.Vec2Int) (*TileData, error) {
	return &TileData{Name: name, ImageData: *imgData, Layer: layer, Size: size}, nil
}
