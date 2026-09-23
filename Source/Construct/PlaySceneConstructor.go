package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	imageassets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Image"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject/Component"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

// とりあえず PlayScene を組み立てる
func PlaySceneConstruct(assetsManager *assets.AssetsManager) Utils.Factory[*scene.PlayScene] {
	playscene, err := scene.NewPlayScene()
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return playscene, fmt.Errorf("NewPlayScene() でエラーが発生しました。\n%w", err)
		}
	}

	ShelfImgData, err := imageassets.NewImageData("Shelf", "みんげー_棚.png", "みんげー_棚.json", assetsManager.Image)
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return playscene, fmt.Errorf("imageassets.NewImageData(\"Shelf\", \"みんげー_棚.png\", \"みんげー_棚.json\", assetsManager.Image.ImageFS) でエラーが発生しました。\n%w", err)
		}
	}
	ShelfTileData, _ := tilemap.NewTileData("Shelf", ShelfImgData)

	mapdata := tilemap.TileMapData{
		MapData: map[tilemap.TilePosition]*tilemap.TileData{
			{X: 0, Y: 0}: ShelfTileData,
			{X: 1, Y: 0}: ShelfTileData,
			{X: 2, Y: 0}: ShelfTileData,
			{X: 3, Y: 0}: ShelfTileData,
			{X: 4, Y: 0}: ShelfTileData,
			{X: 5, Y: 0}: ShelfTileData,
			{X: 6, Y: 0}: ShelfTileData,
			{X: 7, Y: 0}: ShelfTileData,
			{X: 8, Y: 0}: ShelfTileData,
			{X: 9, Y: 0}: ShelfTileData,
			{X: 0, Y: 1}: ShelfTileData,
			{X: 1, Y: 1}: ShelfTileData,
			{X: 2, Y: 1}: ShelfTileData,
			{X: 3, Y: 1}: ShelfTileData,
			{X: 4, Y: 1}: ShelfTileData,
			{X: 5, Y: 1}: ShelfTileData,
			{X: 6, Y: 1}: ShelfTileData,
			{X: 7, Y: 1}: ShelfTileData,
			{X: 8, Y: 1}: ShelfTileData,
			{X: 9, Y: 1}: ShelfTileData,
		},
		TileSize: Utils.Vec2{
			X: 54,
			Y: 54,
		},
	}
	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return func() (gameobject.Transform, error) {
			tf, err := component.NewStdTransform(position)()
			return tf, err
		}
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image) Utils.Factory[gameobject.Render] {
		return Utils.CastFactory[*component.StdRender, gameobject.Render](component.NewStdRender(asp, img))
	}
	tilemapobj, err := NewTileMapObject("TileMapObj", mapdata, playscene, NewTransform, NewRenderer, assetsManager)()
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return playscene, fmt.Errorf("NewTileMapObject() でエラーが発生しました。\n%w", err)
		}
	}
	playscene.AddEntity(tilemapobj)
	return func() (*scene.PlayScene, error) {
		return playscene, nil
	}
}
