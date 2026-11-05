package serializetarget

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type TileMapData struct {
	// タイルサイズ
	TileSize Utils.Vec2Int
	// 床のデータ
	Floor FloorData
	// 床のレイヤー
	FloorLayer string

	// 壁のデータ
	Wall WallData
	// 壁のレイヤー
	WallLayer string

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
	Size     Utils.Vec2Int
	Offset   Utils.Vec2
	Type     tileType
}

type tileType int

const (
	object tileType = iota
	shelf
)

// tileTypeのデシリアライズに使えるようにする
func (tt *tileType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "object":
		*tt = object
	case "shelf":
		*tt = shelf
	case "":
		*tt = object
	default:
		return fmt.Errorf("invalid Hoge value: %s", text)
	}
	return nil
}

type ObjectData struct {
	Tilekey string `toml:"Tile"`
	Layer   string
}

type FloorData struct {
	// 敷き詰めるデフォルトのタイル
	Tile string
	// 敷き詰め領域の左上角
	UpLeft Utils.Vec2Int
	// 敷き詰め領域の右下角
	DownRight Utils.Vec2Int
}

type WallData struct {
	Top    string
	Side   string
	Corner string
}

func (tmd *TileMapData) Init(filePath string, imageManager *imageassets.ImageManager) {
	tmd.filePath = filePath
	tmd.imageManager = imageManager
	tmd.isInitialized = true
}

// *tilemap.TileMapData を生成して返す
func (tmd TileMapData) Construct() (any, error) {
	if !tmd.isInitialized {
		return nil, fmt.Errorf("TileMapData が初期化されていません。 *TileMapData.Init() を実行してください。")
	}

	// タイル種別の一覧
	// レイヤーを入れればタイルがデータが返ってくるようになってる
	Tiles := map[string](func(layer string) (*tilemap.TileData, error)){}
	// タイル毎のOffsetの一覧
	Offsets := map[string]Utils.Vec2{}
	for tilekey, tiledata := range tmd.TileData {
		imagedata, err := imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, tmd.imageManager)
		if err != nil {
			return nil, fmt.Errorf("imageassets.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath, &tmd.imageManager) でエラーが発生しました。\nfilePath: %s, tiledata.Name: %s, tiledata.PngPath: %s, tiledata.JsonPath: %s\n%w", tmd.filePath, tiledata.Name, tiledata.PngPath, tiledata.JsonPath, err)
		}

		// タイルサイズのデフォルト値設定
		tileSize := tiledata.Size
		if tileSize == Utils.NewVecIntZero() {
			tileSize = Utils.NewVecIntOne()
		}

		// 干渉可能設定
		interactivity, err := Interactivity(tiledata.Type)
		if err != nil {
			return nil, fmt.Errorf("Interactivity(tiledata.Type) で失敗しました。\n%w", err)
		}

		tile := makeTileData(tiledata.Name, imagedata, tileSize, interactivity)
		Tiles[tilekey] = tile
		Offsets[tilekey] = tiledata.Offset
	}

	// マップを組み立てる
	mapdata := map[tilemap.TilePosition]*tilemap.TileData{}
	enterabledata := map[Utils.Vec2Int]bool{}

	// 床
	up := tmd.Floor.UpLeft.Y
	left := tmd.Floor.UpLeft.X
	down := tmd.Floor.DownRight.Y
	right := tmd.Floor.DownRight.X

	for j := up; j <= down; j++ {
		for i := left; i <= right; i++ {
			// (i,j) に置く床タイルを指定
			floorKey, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", i, j)]
			if !ok {
				floorKey = tmd.Floor.Tile
			}

			// posを指定
			pos := tilemap.TilePosition{Position: Utils.Vec2Int{X: i, Y: j}, Layer: tmd.FloorLayer, Offset: Offsets[floorKey]}

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

	// 壁
	// 横壁
	for i := left; i <= right; i++ {
		// (i,up-1), (i,down-1) に置く壁タイルを指定
		WallKeyU, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", i, up-1)]
		if !ok {
			WallKeyU = tmd.Wall.Top
		}
		WallKeyD, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", i, down+1)]
		if !ok {
			WallKeyD = tmd.Wall.Top
		}

		// posを指定
		posU := tilemap.TilePosition{Position: Utils.Vec2Int{X: i, Y: up - 1}, Layer: tmd.FloorLayer, Offset: Offsets[WallKeyU]}
		posD := tilemap.TilePosition{Position: Utils.Vec2Int{X: i, Y: down + 1}, Layer: tmd.FloorLayer, Offset: Offsets[WallKeyD]}

		// 実際のタイルのデータを生成
		tileFactoryU, ok := Tiles[WallKeyU]
		if !ok {
			return nil, fmt.Errorf("WallKeyUp のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, WallKeyU)
		}
		tileFactoryD, ok := Tiles[WallKeyD]
		if !ok {
			return nil, fmt.Errorf("WallKeyDown のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, WallKeyD)
		}

		realtiledataU, err := tileFactoryU(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("WallUp %s の生成に失敗しました。\n%w\n", WallKeyU, err)
		}
		realtiledataD, err := tileFactoryD(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("WallDown %s の生成に失敗しました。\n%w\n", WallKeyD, err)
		}

		mapdata[posU] = realtiledataU
		enterabledata[posU.Position] = false
		mapdata[posD] = realtiledataD
		enterabledata[posD.Position] = false
	}
	// 縦壁
	for j := up; j <= down; j++ {
		// (i,up-1), (i,down-1) に置く壁タイルを指定
		WallKeyL, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", left-1, j)]
		if !ok {
			WallKeyL = tmd.Wall.Side
		}
		WallKeyR, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", right+1, j)]
		if !ok {
			WallKeyR = tmd.Wall.Side
		}

		// posを指定
		posL := tilemap.TilePosition{Position: Utils.Vec2Int{X: left - 1, Y: j}, Layer: tmd.FloorLayer, Offset: Offsets[WallKeyL]}
		posR := tilemap.TilePosition{Position: Utils.Vec2Int{X: right + 1, Y: j}, Layer: tmd.FloorLayer, Offset: Offsets[WallKeyR]}

		// 実際のタイルのデータを生成
		tileFactoryL, ok := Tiles[WallKeyL]
		if !ok {
			return nil, fmt.Errorf("WallKeyUp のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, WallKeyL)
		}
		tileFactoryR, ok := Tiles[WallKeyR]
		if !ok {
			return nil, fmt.Errorf("WallKeyDown のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, WallKeyR)
		}

		realtiledataL, err := tileFactoryL(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("WallUp %s の生成に失敗しました。\n%w\n", WallKeyL, err)
		}
		realtiledataR, err := tileFactoryR(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("WallDown %s の生成に失敗しました。\n%w\n", WallKeyR, err)
		}

		mapdata[posL] = realtiledataL
		enterabledata[posL.Position] = false
		mapdata[posR] = realtiledataR
		enterabledata[posR.Position] = false
	}
	// 角
	// 上角を張る
	{
		// 各位置に置く角壁タイルを指定
		CornerKeyUL, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", left-1, up-1)]
		CornerKeyDL := ""
		if !ok {
			CornerKeyUL = tmd.Wall.Corner
			CornerKeyDL = tmd.Wall.Side
		}
		CornerKeyUR, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", right+1, up-1)]
		CornerKeyDR := ""
		if !ok {
			CornerKeyUR = tmd.Wall.Corner
			CornerKeyDR = tmd.Wall.Side
		}

		// posを指定
		posUL := tilemap.TilePosition{Position: Utils.Vec2Int{X: left - 1, Y: up - 1}, Layer: tmd.FloorLayer, Offset: Offsets[CornerKeyUL]}
		posDL := tilemap.TilePosition{Position: Utils.Vec2Int{X: left - 1, Y: up - 1}, Layer: tmd.FloorLayer, Offset: Offsets[CornerKeyDL]}
		posUR := tilemap.TilePosition{Position: Utils.Vec2Int{X: right + 1, Y: up - 1}, Layer: tmd.FloorLayer, Offset: Offsets[CornerKeyUR]}
		posDR := tilemap.TilePosition{Position: Utils.Vec2Int{X: right + 1, Y: up - 1}, Layer: tmd.FloorLayer, Offset: Offsets[CornerKeyDR]}

		// 実際のタイルのデータを生成
		tileFactoryUL, ok := Tiles[CornerKeyUL]
		if !ok {
			return nil, fmt.Errorf("CornerUL のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, CornerKeyUL)
		}
		tileFactoryUR, ok := Tiles[CornerKeyUR]
		if !ok {
			return nil, fmt.Errorf("CornerUR のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, CornerKeyUR)
		}
		tileFactoryDL, ok := Tiles[CornerKeyDL]
		if !ok {
			return nil, fmt.Errorf("CornerDL のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, CornerKeyDL)
		}
		tileFactoryDR, ok := Tiles[CornerKeyDR]
		if !ok {
			return nil, fmt.Errorf("CornerDR のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, CornerKeyDR)
		}

		realtiledataUL, err := tileFactoryUL(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("Floor %s の生成に失敗しました。\n%w\n", CornerKeyUL, err)
		}
		realtiledataUR, err := tileFactoryUR(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("Floor %s の生成に失敗しました。\n%w\n", CornerKeyUR, err)
		}
		realtiledataDL, err := tileFactoryDL(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("Floor %s の生成に失敗しました。\n%w\n", CornerKeyDL, err)
		}
		realtiledataDR, err := tileFactoryDR(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("Floor %s の生成に失敗しました。\n%w\n", CornerKeyDR, err)
		}

		mapdata[posUL] = realtiledataUL
		mapdata[posDL] = realtiledataDL
		enterabledata[posUL.Position] = false
		mapdata[posUR] = realtiledataUR
		mapdata[posDR] = realtiledataDR
		enterabledata[posUR.Position] = false
	}
	// 下角
	{
		// (left,down+1), (right,down+1) に置く壁タイルを指定
		CornerKeyL, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", left-1, down+1)]
		if !ok {
			CornerKeyL = tmd.Wall.Top
		}
		CornerKeyR, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", right+1, down+1)]
		if !ok {
			CornerKeyR = tmd.Wall.Top
		}

		// posを指定
		posL := tilemap.TilePosition{Position: Utils.Vec2Int{X: left - 1, Y: down + 1}, Layer: tmd.FloorLayer, Offset: Offsets[CornerKeyL]}
		posR := tilemap.TilePosition{Position: Utils.Vec2Int{X: right + 1, Y: down + 1}, Layer: tmd.FloorLayer, Offset: Offsets[CornerKeyR]}

		// 実際のタイルのデータを生成
		tileFactoryL, ok := Tiles[CornerKeyL]
		if !ok {
			return nil, fmt.Errorf("CornerKeyL のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, CornerKeyL)
		}
		tileFactoryR, ok := Tiles[CornerKeyR]
		if !ok {
			return nil, fmt.Errorf("CornerKeyR のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, CornerKeyR)
		}

		realtiledataL, err := tileFactoryL(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("CornerL %s の生成に失敗しました。\n%w\n", CornerKeyL, err)
		}
		realtiledataR, err := tileFactoryR(tmd.WallLayer)
		if err != nil {
			return nil, fmt.Errorf("CornerR %s の生成に失敗しました。\n%w\n", CornerKeyR, err)
		}

		mapdata[posL] = realtiledataL
		enterabledata[posL.Position] = false
		mapdata[posR] = realtiledataR
		enterabledata[posR.Position] = false
	}

	// 床以外のオブジェクト
	for posstr, ObjectData := range tmd.ObjectMap {
		//　pos の指定
		posvec := Utils.Vec2Int{}
		_, err := fmt.Sscanf(posstr, "%d,%d", &posvec.X, &posvec.Y)
		if err != nil {
			return nil, fmt.Errorf("MapData のキーが不正です。 \"x,y\" の形式で入力してください。\nfilePath: %s, key: %s\n%w", tmd.filePath, posstr, err)
		}
		pos := tilemap.TilePosition{Position: posvec, Layer: ObjectData.Layer, Offset: Offsets[ObjectData.Tilekey]}

		// 実際のタイルのデータを生成
		tileFactory, ok := Tiles[ObjectData.Tilekey]
		if !ok {
			return nil, fmt.Errorf("FloorLayer のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", tmd.filePath, ObjectData.Tilekey)
		}

		realtiledata, err := tileFactory(ObjectData.Layer)
		if err != nil {
			return nil, fmt.Errorf("Object %s の生成に失敗しました。\n%w\n", ObjectData.Tilekey, err)
		}
		mapdata[pos] = realtiledata
		Size := realtiledata.TileSize
		for offset := range Size.Range() {
			enterabledata[pos.Position.Add(offset)] = false
		}
	}

	result := &tilemap.TileMapData{
		MapData:   mapdata,
		TileSize:  tmd.TileSize,
		Enterable: enterabledata,
	}

	return result, nil
}

// 部分適用
func makeTileData(name string, imgData *imageassets.ImageData, tileSize Utils.Vec2Int, Interact tilemap.InteractType) func(layer string) (*tilemap.TileData, error) {
	return func(layer string) (*tilemap.TileData, error) {
		return tilemap.NewTileData(name, imgData, layer, tileSize, Interact)
	}
}

// tileType から干渉可能設定を取り出す関数
func Interactivity(Type tileType) (tilemap.InteractType, error) {
	switch Type {
	case object:
		return tilemap.None, nil
	case shelf:
		return tilemap.Shelf, nil
	default:
		return tilemap.None, fmt.Errorf("未設定の tileType です。\ntileType = %d", Type)
	}
}
