package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
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
	mapdata := tilemap.TileMapData{
		MapData: map[tilemap.TilePosition]tilemap.TileData{},
		TileSize: Utils.Vec2{
			X: 40,
			Y: 40,
		},
	}
	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return func() (gameobject.Transform, error) {
			tf, err := component.NewStdTransform(position)()
			return tf, err
		}
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image) Utils.Factory[gameobject.Render] {
		rnd, err := component.NewStdRender(asp, img)()
		return func() (gameobject.Render, error) { return rnd, err }
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
