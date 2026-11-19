package serializetarget

import (
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Component"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
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

	// エフェクト名のキーとデータの組
	EffectData map[string]EffectData

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
	Name       string
	JsonPath   string
	PngPath    string
	Size       Utils.Vec2Int
	Offset     Utils.Vec2
	Type       tileType
	Reversible bool
	Effect     string
}

type EffectData struct {
	Name     string
	JsonPath string
	PngPath  string
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
	Bottom string
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

	// エフェクトアセットを取り込む
	Fx := map[string]*imageassets.ImageData{}
	for Fxkey, FxData := range tmd.EffectData {
		imgdata, err := tmd.imageManager.NewImageData(FxData.Name, FxData.PngPath, FxData.JsonPath)
		if err != nil {
			return nil, fmt.Errorf("ImageDataの生成でエラーが発生しました。\nKey : %s\n%w", Fxkey, err)
		}
		Fx[Fxkey] = imgdata
	}

	// タイル種別の一覧
	// レイヤーを入れればタイルがデータが返ってくるようになってる
	Tiles := map[string](func(string, Utils.Vec2) (*tilemap.TileData, error)){}
	ReverseTiles := map[string](func(string, Utils.Vec2) (*tilemap.TileData, error)){}
	// タイル毎のOffsetの一覧
	Offsets := map[string]Utils.Vec2{}
	ReverseOffsets := map[string]Utils.Vec2{}
	for tilekey, tiledata := range tmd.TileData {
		imagedata, err := tmd.imageManager.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath)
		if err != nil {
			return nil, fmt.Errorf("tmd.imageManager.NewImageData(tiledata.Name, tiledata.PngPath, tiledata.JsonPath) でエラーが発生しました。\nfilePath: %s, tiledata.Name: %s, tiledata.PngPath: %s, tiledata.JsonPath: %s\n%w", tmd.filePath, tiledata.Name, tiledata.PngPath, tiledata.JsonPath, err)
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

		// エフェクト
		if fx, ok := Fx[tiledata.Effect]; ok {
			const scale = 2. //　マジックナンバー失礼

			tile := makeShelfData(tiledata.Name, tmd.imageManager, imagedata, fx, tileSize, interactivity, scale, false)
			Tiles[tilekey] = tile
			Offsets[tilekey] = tiledata.Offset
			if tiledata.Reversible {
				// タイル反転用DIO作成
				tsv2 := tmd.TileSize.ToVec2()
				DIO := &ebiten.DrawImageOptions{}
				DIO.GeoM.Scale(-1, 1)
				DIO.GeoM.Translate(tsv2.X, 0)
				revtile := makeShelfData(tiledata.Name, tmd.imageManager, imagedata, fx, tileSize, interactivity, scale, true, tilemap.TileData_DIOOption(DIO))
				ReverseTiles[tilekey] = revtile
				ReverseOffsets[tilekey] = tiledata.Offset.HadamardProd(Utils.Vec2{X: -1, Y: 1})
				return nil, fmt.Errorf("反転処理にバグがあるので止めました。\nこのエラーを消して実行し、反転させたエフェクト付きタイルの位置が適切に処理されているか確認してください。")
			} else {
				ReverseTiles[tilekey] = tile
				ReverseOffsets[tilekey] = tiledata.Offset
			}
		} else if tiledata.Effect != "" {
			return nil, fmt.Errorf("エフェクトのキーが間違っています。\nkey : %s", tiledata.Effect)
		} else {
			tile := makeTileData(tiledata.Name, imagedata, tileSize, interactivity)
			Tiles[tilekey] = tile
			Offsets[tilekey] = tiledata.Offset

			if tiledata.Reversible {
				tsv2 := tmd.TileSize.ToVec2()
				DIO := &ebiten.DrawImageOptions{}
				DIO.GeoM.Scale(-1, 1)
				DIO.GeoM.Translate(tsv2.X, 0)
				revtile := makeTileDataWithDIO(tiledata.Name, imagedata, tileSize, interactivity, DIO)
				ReverseTiles[tilekey] = revtile
				ReverseOffsets[tilekey] = tiledata.Offset.HadamardProd(Utils.Vec2{X: -1, Y: 1})
			} else {
				ReverseTiles[tilekey] = tile
				ReverseOffsets[tilekey] = tiledata.Offset
			}
		}
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
			pos := tilemap.TilePosition{Position: Utils.Vec2Int{X: i, Y: j}, Layer: tmd.FloorLayer}

			if err := tmd.addTileData("Floor", pos, floorKey, tmd.FloorLayer, false, Tiles, Offsets, &mapdata, &enterabledata); err != nil {
				return nil, fmt.Errorf("tmd.addTileData(\"Floor\", pos, floorKey, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, floorKey, err)
			}
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
			WallKeyD = tmd.Wall.Bottom
		}

		// posを指定
		posU := tilemap.TilePosition{Position: Utils.Vec2Int{X: i, Y: up - 1}, Layer: tmd.FloorLayer}
		posD := tilemap.TilePosition{Position: Utils.Vec2Int{X: i, Y: down + 1}, Layer: tmd.FloorLayer}

		if err := tmd.addTileData("WallU", posU, WallKeyU, tmd.WallLayer, true, Tiles, Offsets, &mapdata, &enterabledata); err != nil {
			return nil, fmt.Errorf("tmd.addTileData(\"WallU\", posU, WallKeyU, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, WallKeyU, err)
		}
		if err := tmd.addTileData("WallD", posD, WallKeyD, tmd.WallLayer, true, Tiles, Offsets, &mapdata, &enterabledata); err != nil {
			return nil, fmt.Errorf("tmd.addTileData(\"WallD\", posD, WallKeyD, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, WallKeyD, err)
		}
	}
	// 縦壁 (ふたつのばす)
	for j := up - 2; j <= down; j++ {
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
		posL := tilemap.TilePosition{Position: Utils.Vec2Int{X: left - 1, Y: j}, Layer: tmd.FloorLayer}
		posR := tilemap.TilePosition{Position: Utils.Vec2Int{X: right + 1, Y: j}, Layer: tmd.FloorLayer}

		if err := tmd.addTileData("WallL", posL, WallKeyL, tmd.WallLayer, true, ReverseTiles, ReverseOffsets, &mapdata, &enterabledata); err != nil {
			return nil, fmt.Errorf("tmd.addTileData(\"WallL\", posL, WallKeyL, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, WallKeyL, err)
		}
		if err := tmd.addTileData("WallR", posR, WallKeyR, tmd.WallLayer, true, Tiles, Offsets, &mapdata, &enterabledata); err != nil {
			return nil, fmt.Errorf("tmd.addTileData(\"WallR\", posR, WallKeyR, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, WallKeyR, err)
		}
	}
	// // 角
	// // 上角を張る
	// {
	// 	// 各位置に置く角壁タイルを指定
	// 	CornerKeyUL, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", left-1, up-2)]
	// 	if !ok {
	// 		CornerKeyUL = tmd.Wall.Corner
	// 	}
	// 	CornerKeyUR, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", right+1, up-2)]
	// 	if !ok {
	// 		CornerKeyUR = tmd.Wall.Corner
	// 	}

	// 	// posを指定
	// 	posUL := tilemap.TilePosition{Position: Utils.Vec2Int{X: left - 1, Y: up - 2}, Layer: tmd.FloorLayer}
	// 	posUR := tilemap.TilePosition{Position: Utils.Vec2Int{X: right + 1, Y: up - 2}, Layer: tmd.FloorLayer}

	// 	if err := tmd.addTileData("CornerUL", posUL, CornerKeyUL, tmd.WallLayer, true, ReverseTiles, ReverseOffsets, &mapdata, &enterabledata); err != nil {
	// 		return nil, fmt.Errorf("tmd.addTileData(\"CornerUL\", posUL, CornerKeyUL, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, CornerKeyUL, err)
	// 	}
	// 	if err := tmd.addTileData("CornerUR", posUR, CornerKeyUR, tmd.WallLayer, true, Tiles, Offsets, &mapdata, &enterabledata); err != nil {
	// 		return nil, fmt.Errorf("tmd.addTileData(\"CornerUR\", posUR, CornerKeyUR, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, CornerKeyUR, err)
	// 	}
	// }

	// 下角
	{
		// (left,down+1), (right,down+1) に置く壁タイルを指定
		CornerKeyL, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", left-1, down+1)]
		if !ok {
			CornerKeyL = tmd.Wall.Bottom
		}
		CornerKeyR, ok := tmd.AdditionalFloor[fmt.Sprintf("%d,%d", right+1, down+1)]
		if !ok {
			CornerKeyR = tmd.Wall.Bottom
		}

		// posを指定
		posL := tilemap.TilePosition{Position: Utils.Vec2Int{X: left - 1, Y: down + 1}, Layer: tmd.FloorLayer}
		posR := tilemap.TilePosition{Position: Utils.Vec2Int{X: right + 1, Y: down + 1}, Layer: tmd.FloorLayer}

		if err := tmd.addTileData("CornerDL", posL, CornerKeyL, tmd.WallLayer, true, ReverseTiles, ReverseOffsets, &mapdata, &enterabledata); err != nil {
			return nil, fmt.Errorf("tmd.addTileData(\"CornerDL\", posL, CornerKeyL, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, CornerKeyL, err)
		}
		if err := tmd.addTileData("CornerDR", posR, CornerKeyR, tmd.WallLayer, true, Tiles, Offsets, &mapdata, &enterabledata); err != nil {
			return nil, fmt.Errorf("tmd.addTileData(\"CornerDR\", posR, CornerKeyR, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, CornerKeyR, err)
		}
	}

	// 床以外のオブジェクト
	for posstr, ObjectData := range tmd.ObjectMap {
		//　pos の指定
		posvec := Utils.Vec2Int{}
		_, err := fmt.Sscanf(posstr, "%d,%d", &posvec.X, &posvec.Y)
		if err != nil {
			return nil, fmt.Errorf("MapData のキーが不正です。 \"x,y\" の形式で入力してください。\nfilePath: %s, key: %s\n%w", tmd.filePath, posstr, err)
		}
		pos := tilemap.TilePosition{Position: posvec, Layer: ObjectData.Layer}

		if err := tmd.addTileData("Object", pos, ObjectData.Tilekey, ObjectData.Layer, true, Tiles, Offsets, &mapdata, &enterabledata); err != nil {
			return nil, fmt.Errorf("tmd.addTileData(\"Object\", pos, ObjectData.Tilekey, ...) でエラーが発生しました。\nfilePath: %s, key: %s\n%w", tmd.filePath, ObjectData.Tilekey, err)
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
func makeTileData(name string, imgData *imageassets.ImageData, tileSize Utils.Vec2Int, Interact tilemap.InteractType) func(layer string, offset Utils.Vec2) (*tilemap.TileData, error) {
	return func(layer string, offset Utils.Vec2) (*tilemap.TileData, error) {
		return tilemap.NewTileData(name, imgData, layer, tileSize, offset, Interact)
	}
}

func makeTileDataWithDIO(name string, imgData *imageassets.ImageData, tileSize Utils.Vec2Int, Interact tilemap.InteractType, DIO *ebiten.DrawImageOptions) func(layer string, offset Utils.Vec2) (*tilemap.TileData, error) {
	return func(layer string, offset Utils.Vec2) (*tilemap.TileData, error) {
		return tilemap.NewTileData(name, imgData, layer, tileSize, offset, Interact, tilemap.TileData_DIOOption(DIO))
	}
}

func makeShelfData(name string, imageManager *imageassets.ImageManager, imgData *imageassets.ImageData, FxData *imageassets.ImageData, tileSize Utils.Vec2Int, Interact tilemap.InteractType, scale float64, reverse bool, opts ...tilemap.TileDataOpt) func(layer string, offset Utils.Vec2) (*tilemap.TileData, error) {
	return func(layer string, offset Utils.Vec2) (*tilemap.TileData, error) {
		mainAsp, _, err := imageManager.Load(*imgData)
		if err != nil {
			return nil, fmt.Errorf("棚画像のロードに失敗しました。\n%w", err)
		}
		fxAsp, fxImg, err := imageManager.Load(*FxData)
		if err != nil {
			return nil, fmt.Errorf("エフェクト画像のロードに失敗しました。\n%w", err)
		}

		mw, mh := float64(mainAsp.FrameWidth), float64(mainAsp.FrameHeight)
		fw, fh := float64(fxAsp.FrameWidth)*scale, float64(fxAsp.FrameHeight)*scale

		dio := &ebiten.DrawImageOptions{}
		if reverse {
			dio.GeoM.Scale(-scale, scale)
			dio.GeoM.Translate(fw, 0) // 反転で負側に出た分を戻す
		} else {
			dio.GeoM.Scale(scale, scale)
		}
		dio.GeoM.Translate(math.Floor((mw-fw)/2), math.Floor((mh-fh)/2)) // 中心合わせ

		FxRender := Utils.CastFactory[*component.StdRender, gameobject.Render](component.NewLayerableRender(fxAsp, fxImg, layer, component.SR_WithDrawImageOptions(dio)))
		MultiRender := component.MultiRenderDecorator(FxRender)

		allOpts := append(slices.Clone(opts), tilemap.TileData_RenderOption(MultiRender))

		return tilemap.NewTileData(name, imgData, layer, tileSize, offset, Interact, allOpts...)
	}
}

// タイルデータ追加
func (tmd TileMapData) addTileData(
	Target string,
	pos tilemap.TilePosition,
	Key string,
	Layer string,
	Unenterable bool,
	Tiles map[string](func(string, Utils.Vec2) (*tilemap.TileData, error)),
	Offsets map[string]Utils.Vec2,
	mapdata *map[tilemap.TilePosition]*tilemap.TileData,
	enterabledata *map[Utils.Vec2Int]bool,
) error {
	// 実際のタイルのデータを生成
	tileFactory, ok := Tiles[Key]
	if !ok {
		return fmt.Errorf("%s のタイル指定キーが不正です。 存在するキーを入力してください。\nfilePath: %s, key: %s", Target, tmd.filePath, Key)
	}

	realtiledata, err := tileFactory(Layer, Offsets[Key])
	if err != nil {
		return fmt.Errorf("%s %s の生成に失敗しました。\n%w", Target, Key, err)
	}

	(*mapdata)[pos] = realtiledata

	if Unenterable {
		Size := realtiledata.TileSize
		for offset := range Size.Range() {
			(*enterabledata)[pos.Position.Add(offset)] = false
		}
	}

	return nil
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
