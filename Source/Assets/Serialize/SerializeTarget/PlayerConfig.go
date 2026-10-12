package serializetarget

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
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
	newPlayer func(
		name string,
		position tilemap.TilePosition,
		TileData tilemap.TileData,
		PlScene *scene.PlayScene,
		NewTransform func(Utils.Position) Utils.Factory[gameobject.Transform],
		NewRender Utils.Factory[gameobject.Render],
	) Utils.Factory[*gameobject.GameObject],
) error {
	if assetsManager == nil {
		return fmt.Errorf("assetsManager が nil です。")
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

	imageData, err := imageassets.NewImageData(pc.Name, pc.TileData.PngPath, pc.TileData.JsonPath, pc.assetsManager.Image)
	if err != nil {
		return nil, fmt.Errorf("imageassets.NewImageData でエラーが発生しました。\n%w", err)
	}

	td, err := tilemap.NewTileData(pc.TileData.Name, imageData, pc.Position.Layer)
	if err != nil {
		return nil, fmt.Errorf("tilemap.NewTileData でエラーが発生しました。\n%w", err)
	}

	return PlayerConstructData{Name: pc.Name, Position: pc.Position, TileData: td}, nil
}

// test
var _ serialize.Serializable = PlayerConfig{}
