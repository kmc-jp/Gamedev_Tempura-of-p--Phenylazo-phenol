package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type TileMap struct {
	Tiles    []*Tile
	TileSize Utils.Vec2
}

// test
var _ gameobject.Component = TileMap{}

func (h TileMap) Update(active bool) (err error) {
	// なにもしない
	return nil
}

func NewTileMap(mapdata tilemap.TileMapData) Utils.Factory[TileMap] {
	c := len(mapdata.MapData)
	tilelist := make(map[tilemap.TilePosition]*Tile, c)

	for pos, tile := range mapdata.MapData {
		t, err := NewTile(tile)()
		if err != nil {
			return func() (*TileMap, error) {
				return nil, fmt.Errorf("NewTile() でエラーが発生しました。\n%w", err)
			}
		}

		// tilelist に登録
		tilelist[pos] = t
	}
	t := TileMap{}
	return func() (*TileMap, error) {
		return &t, nil
	}
}

// test
var _ Utils.Factory[HogeComponent] = NewHogeComponent(nil)
