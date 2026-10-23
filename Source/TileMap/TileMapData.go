package tilemap

import "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

type TileMapData struct {
	MapData   map[TilePosition]*TileData
	TileSize  Utils.Vec2Int
	Enterable map[Utils.Vec2Int]bool
	Interact  map[Utils.Vec2Int]InteractType
}

type InteractType int

const (
	None InteractType = iota
	Shelf
)
