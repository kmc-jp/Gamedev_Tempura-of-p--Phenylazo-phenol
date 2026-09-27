package tilemap

import "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

type TileMapData struct {
	MapData  map[TilePosition]*TileData
	TileSize Utils.Vec2
}
