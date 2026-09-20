package construct

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject/Component"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

// とりあえず PlayScene を組み立てる
func PlaySceneConstruct() Utils.Factory[*scene.PlayScene] {
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
	NewRenderer := func() (gameobject.Render, error) {
		rnd, err := component.NewStdRender()
		return rnd, err
	}
	tilemapobj, err := NewTileMapObject("TileMapObj", mapdata, playscene, NewTransform, NewRenderer)()
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
