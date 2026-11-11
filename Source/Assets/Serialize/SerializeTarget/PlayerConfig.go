package serializetarget

import (
	assetsutil "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/AssetsUtils"
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type PlayerConfig struct {
	Name     string
	Position tilemap.TilePosition // 座標とレイヤー
	Walktime int                  // 歩行にかかる時間

	// タグ
	U string
	D string
	L string
	R string

	TileData     TileData                  // Source\Assets\Serialize\SerializeTarget\TileMapData.go 参照
	initialized  bool                      // シリアライズされない
	imageManager *imageassets.ImageManager // シリアライズされない
}

func (pc *PlayerConfig) Init(
	imageManager *imageassets.ImageManager,
) error {
	if imageManager == nil {
		return fmt.Errorf("assetsManager が nil です。")
	}

	// ファイルパスチェック
	_, err := assetsutil.CheckAndCleanPath(pc.TileData.JsonPath, imageManager.GetFS())
	if err != nil {
		return fmt.Errorf("assetsutil.CheckAndCleanPath(pc.TileData.JsonPath) でエラーが発生しました。 JsonPath が不正です。\n%w", err)
	}
	_, err = assetsutil.CheckAndCleanPath(pc.TileData.PngPath, imageManager.GetFS())
	if err != nil {
		return fmt.Errorf("assetsutil.CheckAndCleanPath(pc.TileData.PngPath) でエラーが発生しました。 PngPath が不正です。\n%w", err)
	}

	// タイルサイズのデフォルト値設定
	if pc.TileData.Size == Utils.NewVecIntZero() {
		pc.TileData.Size = Utils.NewVecIntOne()
	}

	pc.initialized = true
	pc.imageManager = imageManager
	return nil
}

type PlayerConstructData struct {
	Name     string
	Position tilemap.TilePosition
	TileData *tilemap.TileData
	Walktime int
	DirTag   map[Utils.Vec2Int]string
}

// プレイヤーを生成するために必要な情報を生成して返す
func (pc PlayerConfig) Construct() (any, error) {
	if !pc.initialized {
		return nil, fmt.Errorf("PlayerConfigがイニシャライズされていません。")
	}

	if pc.TileData.Size == Utils.NewVecIntZero() {
		return nil, fmt.Errorf("TileData.Size がゼロ値になっています。")
	}

	imageData, err := imageassets.NewImageData(pc.Name, pc.TileData.PngPath, pc.TileData.JsonPath, pc.imageManager)
	if err != nil {
		return nil, fmt.Errorf("imageassets.NewImageData でエラーが発生しました。\n%w", err)
	}

	td, err := tilemap.NewTileData(pc.TileData.Name, imageData, pc.Position.Layer, pc.TileData.Size, tilemap.None)
	if err != nil {
		return nil, fmt.Errorf("tilemap.NewTileData でエラーが発生しました。\n%w", err)
	}

	dirTag := map[Utils.Vec2Int]string{
		Utils.NewVecIntUp():    pc.U,
		Utils.NewVecIntDown():  pc.D,
		Utils.NewVecIntLeft():  pc.L,
		Utils.NewVecIntRight(): pc.R,
	}

	return PlayerConstructData{
		Name:     pc.Name,
		Position: pc.Position,
		TileData: td,
		Walktime: pc.Walktime,
		DirTag:   dirTag,
	}, nil
}

// test
var _ serialize.Serializable = PlayerConfig{}
