package tilemap

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
)

type TileData struct {
	Name      string
	ImageData imageassets.ImageData
	Layer     string
}

func NewTileData(name string, imgData *imageassets.ImageData, layer string) (*TileData, error) {
	return &TileData{Name: name, ImageData: *imgData, Layer: layer}, nil
}
