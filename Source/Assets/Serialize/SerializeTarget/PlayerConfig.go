package serializetarget

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type PlayerConfig struct {
	Name          string
	Position      tilemap.TilePosition  // 座標とレイヤー
	TileData      TileData              // Source\Assets\Serialize\SerializeTarget\TileMapData.go 参照
	initialized   bool                  // シリアライズされない
	assetsManager *assets.AssetsManager // シリアライズされない
}

func (pc *PlayerConfig) Init(
	assetsManager *assets.AssetsManager,
) error {
	if assetsManager == nil {
		return fmt.Errorf("assetsManager が nil です。")
	}

	// ファイルパスチェック
	_, err := assetsManager.Image.CheckAndCleanPath(pc.TileData.JsonPath)
	if err != nil {
		return fmt.Errorf("assetsManager.Image.CheckAndCleanPath(pc.TileData.JsonPath) でエラーが発生しました。 JsonPath が不正です。\n%w", err)
	}
	_, err = assetsManager.Image.CheckAndCleanPath(pc.TileData.PngPath)
	if err != nil {
		return fmt.Errorf("assetsManager.Image.CheckAndCleanPath(pc.TileData.PngPath) でエラーが発生しました。 PngPath が不正です。\n%w", err)
	}

	// タイルサイズのデフォルト値設定
	if pc.TileData.Size == Utils.NewVecIntZero() {
		pc.TileData.Size = Utils.NewVecIntOne()
	}

	pc.initialized = true
	pc.assetsManager = assetsManager
	return nil
}

type PlayerConstructData struct {
	Name     string
	Position tilemap.TilePosition
	TileData *tilemap.TileData
}

// プレイヤーを生成するために必要な情報を生成して返す
func (pc PlayerConfig) Construct() (any, error) {
	if !pc.initialized {
		return nil, fmt.Errorf("PlayerConfigがイニシャライズされていません。")
	}

	if pc.TileData.Size == Utils.NewVecIntZero() {
		return nil, fmt.Errorf("TileData.Size がゼロ値になっています。")
	}

	imageData, err := imageassets.NewImageData(pc.Name, pc.TileData.PngPath, pc.TileData.JsonPath, pc.assetsManager.Image)
	if err != nil {
		return nil, fmt.Errorf("imageassets.NewImageData でエラーが発生しました。\n%w", err)
	}

	td, err := tilemap.NewTileData(pc.TileData.Name, imageData, pc.Position.Layer, pc.TileData.Size)
	if err != nil {
		return nil, fmt.Errorf("tilemap.NewTileData でエラーが発生しました。\n%w", err)
	}

	return PlayerConstructData{Name: pc.Name, Position: pc.Position, TileData: td}, nil
}

// test
var _ serialize.Serializable = PlayerConfig{}
