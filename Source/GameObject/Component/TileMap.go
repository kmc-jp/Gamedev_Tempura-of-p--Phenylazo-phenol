package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type TileMap struct {
	Tiles    map[tilemap.TilePosition]*Tile
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

func (tm TileMap) SetGameObject(scene *scene.PlayScene, NewTileObj func(name string, data tilemap.TileData, scene *scene.PlayScene, transformFactry Utils.Factory[gameobject.Transform]) Utils.Factory[gameobject.GameObject], NewTransform func(Utils.Position) Utils.Factory[gameobject.Transform]) error {
	tilecount := map[tilemap.TileData]int{}

	for pos, t := range tm.Tiles {
		data := t.Data
		// 名前
		i, _ := tilecount[data]
		tilecount[data] = i + 1 // キーが存在しなければ i = 0 なので問題なし
		name := fmt.Sprintf("%s%d", data.Name, i+1)

		// 位置
		position := pos.ToPosition().PtoV().HadamardProd(tm.TileSize).VtoP()
		transformFactry := NewTransform(position)
		// GameObject を作成
		obj, err := NewTileObj(name, data, scene, transformFactry)()
		if err != nil {
			return fmt.Errorf("NewTileObj() でエラーが発生しました。\n%w", err)
		}

		scene.AddEntity(obj)
	}
	return nil
}

// test
var _ Utils.Factory[HogeComponent] = NewHogeComponent(nil)
