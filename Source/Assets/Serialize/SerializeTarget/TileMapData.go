package serializetarget

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type TileMapData struct {
	// タイルサイズ
	TileSize Utils.Vec2
	// 床のデータ
	Floor FloorData

	// タイル名のキーとデータの組
	TileData map[string]TileData

	// Floor 以外のタイルで敷き詰めたい床
	// 位置 "x,y" とタイルキーの組
	AdditionalFloor map[string]string

	// 床の上に置いてあるもの
	MapData map[string]string

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

type FloorData struct {
	// 敷き詰めるデフォルトのタイル
	Tile string
	// 敷き詰め領域の左上角
	UpLeft tilemap.TilePosition
	// 敷き詰め領域の右下角
	DownRight tilemap.TilePosition
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

	// タイル種別の一覧
	Tiles := map[string]*tilemap.TileData{}
	for tilekey, tiledata := range tmd.TileData {
		imagedata, err := imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, tmd.imageManager)
		if err != nil {
			return nil, fmt.Errorf("imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, &tmd.imageManager) でエラーが発生しました。\nfilePath: %s, tiledata.Name: %s, tiledata.PngPath: %s, tiledata.JsonPath: %s\n%w", tmd.filePath, tiledata.Name, tiledata.PngPath, tiledata.JsonPath, err)
		}
		tile, err := tilemap.NewTileData(tiledata.Name, imagedata, "")
		Tiles[tilekey] = tile
	}

	mapdata := map[tilemap.TilePosition]*tilemap.TileData{}
	for posstr, tilekey := range tmd.MapData {
		pos := tilemap.TilePosition{Layer: Tiles[tilekey].Layer}
		_, err := fmt.Sscanf(posstr, "%d,%d", &pos.X, &pos.Y)
		if err != nil {
			return nil, fmt.Errorf("MapData のキーが不正です。 \"x,y\" の形式で入力してください。\nfilePath: %s, key: %s\n%w", tmd.filePath, posstr, err)
		}
		realtiledata, ok := Tiles[tilekey]
		if !ok {
			return nil, fmt.Errorf("MapData のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, tilekey)
		}
		mapdata[pos] = realtiledata
	}

	result := tilemap.TileMapData{
		MapData:  mapdata,
		TileSize: tmd.TileSize,
	}

	return result, nil
}
