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
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/goaseprite"
)

// とりあえず PlayScene を組み立てる
func PlaySceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler) Utils.Factory[*scene.PlayScene] {
	// カメラ設定
	DrawerFactory, err := MakeDrawerFactory(assetsManager)
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}
	}

	playscene, err := scene.NewPlayScene(DrawerFactory, inputHandler)()
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return nil, fmt.Errorf("NewPlayScene() でエラーが発生しました。\n%w", err)
		}
	}

	// タイルマップ
	TileMapData, err := MakeTileMap(assetsManager)
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return nil, fmt.Errorf("MakeTileMap(assetsManager) でエラーが発生しました。\n%w", err)
		}
	}

	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return Utils.CastFactory[*component.StdTransform, gameobject.Transform](component.NewStdTransform(position))
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[gameobject.Render] {
		return Utils.CastFactory[*component.StdRender, gameobject.Render](component.NewLayerableRender(asp, img, layer))
	}
	tilemapobj, err := NewTileMapObject("TileMapObj", *TileMapData, playscene, NewTransform, NewRenderer, assetsManager)()
	if err != nil {
		return func() (*scene.PlayScene, error) {
			return nil, fmt.Errorf("NewTileMapObject() でエラーが発生しました。\n%w", err)
		}
	}
	playscene.AddEntity(tilemapobj)
	return func() (*scene.PlayScene, error) {
		return playscene, nil
	}

	// Player を設置

}

// カメラ設定
func MakeDrawerFactory(assetsManager *assets.AssetsManager) (Utils.Factory[scene.Drawer], error) {
	cameraConfigData, err := serialize.NewSerializeData("CameraConfig", reflect.TypeFor[serializetarget.CameraConfig](), "CameraConfig.toml", assetsManager.Serialize)
	if err != nil {
		return nil, fmt.Errorf("serialize.NewSerializeData(\"CameraConfig\", reflect.TypeFor[serializetarget.CameraConfig](), \"CameraConfig.toml\", assetsManager.Serialize) でエラーが発生しました。\n%w", err)
	}
	cameraConfigAny, err := assetsManager.Serialize.Load(*cameraConfigData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*cameraConfigData) でエラーが発生しました。\n%w", err)
	}

	cameraConfig := cameraConfigAny.(serializetarget.CameraConfig)
	cameraFactory, err := cameraConfig.Construct()
	if err != nil {
		return nil, fmt.Errorf("cameraConfig.Construct() でエラーが発生しました。\n%w", err)
	}

	DrawerFactory := Utils.CastFactory[*scene.Camera, scene.Drawer](cameraFactory.(Utils.Factory[*scene.Camera]))
	return DrawerFactory, nil
}

// タイルマップ設定
func MakeTileMap(assetsManager *assets.AssetsManager) (*tilemap.TileMapData, error) {
	TileMapSeriData, err := serialize.NewSerializeData("shelf", reflect.TypeFor[serializetarget.TileMapData](), "TileMap.toml", assetsManager.Serialize)
	if err != nil {
		return nil, fmt.Errorf("serialize.NewSerializeData(\"shelf\", reflect.TypeFor[serializetarget.TileMapData](), \"TileMap.toml\", assetsManager.Serialize) でエラーが発生しました。\n%w", err)
	}
	mapdataAny, err := assetsManager.Serialize.Load(*TileMapSeriData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*shelfdata) でエラーが発生しました。\n%w", err)
	}
	mapdataDeserialized := mapdataAny.(serializetarget.TileMapData)
	mapdataDeserialized.Init("TileMap.toml", assetsManager.Image)
	mapdataConstructed, err := mapdataDeserialized.Construct()
	if err != nil {
		return nil, fmt.Errorf("mapdataDeserialized.Construct() でエラーが発生しました。\n%w", err)
	}

	return mapdataConstructed.(*tilemap.TileMapData), nil
}
