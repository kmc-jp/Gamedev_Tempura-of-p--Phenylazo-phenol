package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	serializetarget "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize/SerializeTarget"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject/Component"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"reflect"

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

	serimanager := assetsManager.Serialize
	shelfdata, err := serialize.NewSerializeData("shelf", reflect.TypeFor[serializetarget.TileMapData](), "TileMap.toml", serimanager)
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return playscene, fmt.Errorf("serialize.NewSerializeData(\"shelf\", reflect.TypeFor[serializetarget.TileMapData](), \"TileMap.toml\", serimanager) でエラーが発生しました。\n%w", err)
		}
	}
	mapdataAny, err := assetsManager.Serialize.Load(*shelfdata)
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return playscene, fmt.Errorf("assetsManager.Serialize.Load(*shelfdata) でエラーが発生しました。\n%w", err)
		}
	}
	mapdataDeserialized := mapdataAny.(serializetarget.TileMapData)
	mapdataDeserialized.Init("TileMap.toml", assetsManager.Image)
	mapdataConstructed, err := mapdataDeserialized.Construct()
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return playscene, fmt.Errorf("mapdataDeserialized.Construct() でエラーが発生しました。\n%w", err)
		}
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
	tilemapobj, err := NewTileMapObject("TileMapObj", mapdataConstructed.(tilemap.TileMapData), playscene, NewTransform, NewRenderer, assetsManager)()
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
