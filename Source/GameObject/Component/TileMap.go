package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type TileMap struct {
	Tiles     map[tilemap.TilePosition]*Tile
	TileSize  Utils.Vec2Int
	Enterable map[Utils.Vec2Int]bool // 何も入ってなければ進入可能
}

// test
var _ gameobject.Component = TileMap{}

func (h TileMap) Update(active bool) (err error) {
	// なにもしない
	return nil
}

func NewTileMap(mapdata tilemap.TileMapData) Utils.Factory[*TileMap] {
	c := len(mapdata.MapData)
	tilelist := make(map[tilemap.TilePosition]*Tile, c)

	for pos, tile := range mapdata.MapData {
		t, err := NewTile(*tile)()
		if err != nil {
			return func() (*TileMap, error) {
				return nil, fmt.Errorf("NewTile() でエラーが発生しました。\n%w", err)
			}
		}

		// tilelist に登録
		tilelist[pos] = t
	}
	t := TileMap{
		Tiles:     tilelist,
		TileSize:  mapdata.TileSize,
		Enterable: mapdata.Enterable,
	}
	return func() (*TileMap, error) {
		return &t, nil
	}
}

func (tm TileMap) SetGameObject(
	scene *scene.PlayScene,
	NewTileObj func(
		name string,
		data tilemap.TileData,
		scene *scene.PlayScene,
		layer string,
		transformFactry Utils.Factory[gameobject.Transform],
	) Utils.Factory[*gameobject.GameObject],
	NewTransform func(Utils.Position) Utils.Factory[gameobject.Transform],
) error {
	println("TileMap.SetGameObject()")

	tilecount := map[tilemap.TileData]int{}
	fmt.Printf("tilecount : %d\n", len(tm.Tiles))

	for pos, t := range tm.Tiles {
		data := t.Data
		// 名前
		i, _ := tilecount[data]
		tilecount[data] = i + 1 // キーが存在しなければ i = 0 なので問題なし
		name := fmt.Sprintf("%s%d", data.Name, i+1)

		// 位置
		position := tm.TilePositionToWorldPosition(pos, data)
		transformFactry := NewTransform(position)

		// レイヤー
		layer := pos.Layer

		// GameObject を作成
		obj, err := NewTileObj(name, data, scene, layer, transformFactry)()
		if err != nil {
			return fmt.Errorf("NewTileObj() でエラーが発生しました。\n%w", err)
		}

		scene.AddEntity(obj)

		// fmt.Printf("Tile %s が追加されました。\n", obj.Name())
	}
	return nil
}

func (tm TileMap) TilePositionToWorldPosition(pos tilemap.TilePosition, data tilemap.TileData) Utils.Position {
	tileposf := pos.Position.ToVec2()
	sizeoffset := data.Size.Sub(Utils.NewVecIntOne()).ToVec2().HadamardDiv(data.Size.ToVec2())
	to := tileposf.Add(sizeoffset).HadamardProd(tm.TileSize.ToVec2()).Add(pos.Offset)
	// if data.Size != Utils.NewVecIntOne() {
	// 	fmt.Printf("TilePositionToWorldPosition\npos %d -> size %d -> %f\ntileposf = %f, pos.Offset = %f, sizeoffset = %f", pos.Position, data.Size, to.HadamardDiv(tm.TileSize.ToVec2()), tileposf, pos.Offset, sizeoffset)
	// }
	return to.VtoP()
}
