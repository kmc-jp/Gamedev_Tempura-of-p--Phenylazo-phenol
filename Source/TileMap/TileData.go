package tilemap

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
)

type TileData struct {
	Name      string
	ImageData imageassets.ImageData
}

func NewTileData(name string, imgData imageassets.ImageData) (*TileData, error) {
	return &TileData{Name: name, ImageData: imgData}, nil
}
