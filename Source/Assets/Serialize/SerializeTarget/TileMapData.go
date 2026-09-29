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
	// 床のレイヤー
	FloorLayer string

	// タイル名のキーとデータの組
	TileData map[string]TileData

	// Floor 以外のタイルで敷き詰めたい床
	// 位置 "x,y" とタイルキーの組
	AdditionalFloor map[string]string

	// 床の上に置いてあるもの
	ObjectMap map[string]ObjectData

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

type ObjectData struct {
	Tile  string
	Layer string
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
	// レイヤーを入れればタイルがデータが返ってくるようになってる
	Tiles := map[string](func(layer string) (*tilemap.TileData, error)){}
	for tilekey, tiledata := range tmd.TileData {
		imagedata, err := imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, tmd.imageManager)
		if err != nil {
			return nil, fmt.Errorf("imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, &tmd.imageManager) でエラーが発生しました。\nfilePath: %s, tiledata.Name: %s, tiledata.PngPath: %s, tiledata.JsonPath: %s\n%w", tmd.filePath, tiledata.Name, tiledata.PngPath, tiledata.JsonPath, err)
		}
		tile := makeTileData(tiledata.Name, imagedata)
		Tiles[tilekey] = tile
	}

	// マップを組み立てる
	mapdata := map[tilemap.TilePosition]*tilemap.TileData{}

	// 床
	up := tmd.Floor.UpLeft.Y
	left := tmd.Floor.UpLeft.X
	down := tmd.Floor.DownRight.Y
	right := tmd.Floor.DownRight.X
	fmt.Printf("%d,%d,%d,%d,", up, left, down, right)

	for j := up; j <= down; j++ {
		for i := left; i <= right; i++ {
			// (i,j) に置く床タイルを指定
			floorKey, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", i, j)]
			if !ok {
				floorKey = tmd.Floor.Tile
			}

			// posを指定
			pos := tilemap.TilePosition{X: i, Y: j, Layer: tmd.FloorLayer}

			// 実際のタイルのデータを生成
			tileFactory, ok := Tiles[floorKey]
			if !ok {
				return nil, fmt.Errorf("FloorLayer のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, floorKey)
			}

			realtiledata, err := tileFactory(tmd.FloorLayer)
			if err != nil {
				return nil, fmt.Errorf("Floor %s の生成に失敗しました。\n%w\n", floorKey, err)
			}

			mapdata[pos] = realtiledata
		}
	}

	for posstr, ObjectData := range tmd.ObjectMap {
		//　pos の指定
		pos := tilemap.TilePosition{Layer: ObjectData.Layer}
		_, err := fmt.Sscanf(posstr, "%d,%d", &pos.X, &pos.Y)
		if err != nil {
			return nil, fmt.Errorf("MapData のキーが不正です。 \"x,y\" の形式で入力してください。\nfilePath: %s, key: %s\n%w", tmd.filePath, posstr, err)
		}

		// 実際のタイルのデータを生成
		tileFactory, ok := Tiles[ObjectData.Tile]
		if !ok {
			return nil, fmt.Errorf("FloorLayer のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, ObjectData.Tile)
		}

		realtiledata, err := tileFactory(ObjectData.Layer)
		if err != nil {
			return nil, fmt.Errorf("Object %s の生成に失敗しました。\n%w\n", ObjectData.Tile, err)
		}
		mapdata[pos] = realtiledata
	}

	result := tilemap.TileMapData{
		MapData:  mapdata,
		TileSize: tmd.TileSize,
	}

	return result, nil
}

// 部分適用
func makeTileData(name string, imgData *imageassets.ImageData) func(layer string) (*tilemap.TileData, error) {
	return func(layer string) (*tilemap.TileData, error) {
		return tilemap.NewTileData(name, imgData, "hogehoge")
	}
}
