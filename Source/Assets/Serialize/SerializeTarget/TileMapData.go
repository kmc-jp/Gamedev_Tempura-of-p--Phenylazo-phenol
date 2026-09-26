package serializetarget

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type TileMapData struct {
	// タイル名のキーとデータの組
	TileData map[string]TileData
	// 位置 "x,y" とタイルキーの組
	MapData map[string]string
	// タイルサイズ
	TileSize Utils.Vec2
	// デコードされない : ファイルパスを入れておく
	filePath string
	// デコードされない
	imageManager *imageassets.ImageManager
	// デコードされない : 初期化した？
	isInitialized bool
}

type TileData struct {
	Name     string
	JsonPath string
	PngPath  string
}

func (tmd *TileMapData) Init(filePath string, imageManager *imageassets.ImageManager) {
	tmd.filePath = filePath
	tmd.imageManager = imageManager
	tmd.isInitialized = true
}

// tilemap.TileMapData を生成して返す
func (tmd TileMapData) Construct() (any, error) {
	if !tmd.isInitialized {
		return nil, fmt.Errorf("TileMapData が初期化されていません。 *TileMapData.Init() を実行してください。")
	}

	Tiles := map[string]*tilemap.TileData{}
	for tilekey, tiledata := range tmd.TileData {
		imagedata, err := imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, tmd.imageManager)
		if err != nil {
			return nil, fmt.Errorf("imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, &tmd.imageManager) でエラーが発生しました。\nfilePath: %s, tiledata.Name: %s, tiledata.PngPath: %s, tiledata.JsonPath: %s\n%w", tmd.filePath, tiledata.Name, tiledata.PngPath, tiledata.JsonPath, err)
		}
		tile, err := tilemap.NewTileData(tiledata.Name, imagedata)
		Tiles[tilekey] = tile
	}

	mapdata := map[tilemap.TilePosition]*tilemap.TileData{}
	for posstr, tilekey := range tmd.MapData {
		pos := tilemap.TilePosition{}
		_, err := fmt.Scanf(posstr, "%d,%d", &pos.X, &pos.Y)
		if err != nil {
			return nil, fmt.Errorf("MapData のキーが不正です。 \"x,y\" の形式で入力してください。\nfilePath: %s, key: %s\n%w", tmd.filePath, posstr, err)
		}
		mapdata[pos] = Tiles[tilekey]
	}

	result := tilemap.TileMapData{
		MapData:  mapdata,
		TileSize: tmd.TileSize,
	}

	return result, nil
}
