package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	serializetarget "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize/SerializeTarget"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Component"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	scenecanvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/SceneCanvas"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/goaseprite"
)

// とりあえず PlayScene を組み立てる
func PlaySceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes, Stage2URL scene.Stage2URL, stage scene.Stage) Utils.Factory[*scene.PlayScene] {
	return func() (*scene.PlayScene, error) {
		// カメラ設定
		WorldDrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		// UI（CanvasはSceneを参照するため、Sceneの生成時に生成関数から作る）
		newCanvas := func(s *scene.PlayScene) (scene.PlayCanvasAPI, error) {
			playCanvas, err := scenecanvas.NewPlayCanvas(s, gdata)
			if err != nil {
				return nil, fmt.Errorf("scenecanvas.NewPlayCanvas() でエラーが発生しました。\n%w", err)
			}
			return playCanvas, nil
		}

		playscene, err := scene.NewPlayScene(WorldDrawerFactory, newCanvas, inputHandler, gdata, initScenes, stage)()
		if err != nil {
			return nil, fmt.Errorf("NewPlayScene() でエラーが発生しました。\n%w", err)
		}

		// タイルマップ
		tilemapobj, err := MakeTileMapObj(assetsManager, Stage2URL, stage, playscene)
		if err != nil {
			return nil, fmt.Errorf("MakeTileMapObj(assetsManager, playscene) でエラーが発生しました。\n%w", err)
		}
		playscene.AddEntity(tilemapobj)

		// Player を設置 (タイルマップ情報が必要)
		TMC, err := tilemapobj.GetComponent(reflect.TypeFor[*component.TileMap]())
		if err != nil {
			return nil, fmt.Errorf("tilemapobj.GetComponent(reflect.TypeFor[*component.TileMap]()) でエラーが発生しました。\n%w", err)
		}
		TM := TMC.(*component.TileMap)
		plobj, err := MakePlayer(assetsManager, playscene, TM)
		if err != nil {
			return nil, fmt.Errorf("MakePlayer(assetsManager, playscene) でエラーが発生しました。\n%w", err)
		}
		playscene.AddEntity(plobj)

		return playscene, nil
	}
}

// ステージ類設定
func MakeStage2URL(assetsManager *assets.AssetsManager) (scene.Stage2URL, error) {
	StagesSeriData, err := serialize.NewSerializeData("Stages", reflect.TypeFor[serializetarget.StagesConfig](), "StageConfig.toml", assetsManager.Serialize)
	if err != nil {
		return nil, fmt.Errorf("serialize.NewSerializeData(\"Stages\", reflect.TypeFor[serializetarget.TileMapData](), \"StageConfig.toml\", assetsManager.Serialize) でエラーが発生しました。\n%w", err)
	}
	StagesDataAny, err := assetsManager.Serialize.Load(*StagesSeriData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*shelfdata) でエラーが発生しました。\n%w", err)
	}
	StagesDataDeserialized := StagesDataAny.(serializetarget.StagesConfig)
	StagesDataConstructed, err := StagesDataDeserialized.Construct()
	if err != nil {
		return nil, fmt.Errorf("mapdataDeserialized.Construct() でエラーが発生しました。\n%w", err)
	}

	return StagesDataConstructed.(scene.Stage2URL), nil
}

// タイルマップ設定
func MakeTileMapData(assetsManager *assets.AssetsManager, URL string) (*tilemap.TileMapData, error) {
	TileMapSeriData, err := serialize.NewSerializeData("shelf", reflect.TypeFor[serializetarget.TileMapData](), URL, assetsManager.Serialize)
	if err != nil {
		return nil, fmt.Errorf("serialize.NewSerializeData(\"shelf\", reflect.TypeFor[serializetarget.TileMapData](), \"%s\", assetsManager.Serialize) でエラーが発生しました。\n%w", URL, err)
	}
	mapdataAny, err := assetsManager.Serialize.Load(*TileMapSeriData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*shelfdata) でエラーが発生しました。\n%w", err)
	}
	mapdataDeserialized := mapdataAny.(serializetarget.TileMapData)
	mapdataDeserialized.Init(URL, assetsManager.Image)
	mapdataConstructed, err := mapdataDeserialized.Construct()
	if err != nil {
		return nil, fmt.Errorf("mapdataDeserialized.Construct() でエラーが発生しました。\n%w", err)
	}

	return mapdataConstructed.(*tilemap.TileMapData), nil
}

func MakeTileMapObj(assetsManager *assets.AssetsManager, Stage2URL scene.Stage2URL, stage scene.Stage, PlScene *scene.PlayScene) (*gameobject.GameObject, error) {
	TileMapData, err := MakeTileMapData(assetsManager, Stage2URL[stage])
	if err != nil {
		return nil, fmt.Errorf("MakeTileMap(assetsManager) でエラーが発生しました。\n%w", err)
	}

	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return Utils.CastFactory[*component.StdTransform, gameobject.Transform](component.NewStdTransform(position))
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image, DIO *ebiten.DrawImageOptions, layer string) Utils.Factory[gameobject.Render] {
		return Utils.CastFactory[*component.StdRender, gameobject.Render](component.NewLayerableRender(asp, img, layer, component.SR_WithDrawImageOptions(DIO)))
	}
	tilemapobj, err := NewTileMapObject("TileMapObj", *TileMapData, PlScene, NewTransform, NewRenderer, assetsManager)()
	if err != nil {
		return nil, fmt.Errorf("NewTileMapObject() でエラーが発生しました。\n%w", err)
	}

	return tilemapobj, nil
}

// プレイヤー設定
func MakePlayer(assetsManager *assets.AssetsManager, PlScene *scene.PlayScene, TileMap *component.TileMap) (*gameobject.GameObject, error) {
	// PlayerConfig
	SerializeData, err := serialize.NewSerializeData("Player", reflect.TypeFor[serializetarget.PlayerConfig](), "Player.toml", assetsManager.Serialize)
	if err != nil {
		return nil, fmt.Errorf("serialize.NewSerializeData(\"Player\", reflect.TypeFor[serializetarget.PlayerConfig](), \"Player.toml\", assetsManager.Serialize) でエラーが発生しました。\n%w", err)
	}

	ConfigAny, err := assetsManager.Serialize.Load(*SerializeData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*SerializeData) でエラーが発生しました。\n%w", err)
	}

	PlayerConfig := ConfigAny.(serializetarget.PlayerConfig)

	err = PlayerConfig.Init(assetsManager.Image)
	if err != nil {
		return nil, fmt.Errorf("PlayerConfig.Init(assetsManager) でエラーが発生しました。\n%w", err)
	}

	constructDataAny, err := PlayerConfig.Construct()
	if err != nil {
		return nil, fmt.Errorf("PlayerConfig.Construct() でエラーが発生しました。\n%w", err)
	}
	constructData := constructDataAny.(serializetarget.PlayerConstructData)

	// Transform / AnimRender
	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return Utils.CastFactory[*component.StdTransform, gameobject.Transform](component.NewStdTransform(position))
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[*component.AnimRender] {
		return component.NewAnimRender(asp, img, layer, "")
	}

	// オブジェクト生成
	plobj, err := NewPlayer("Player", constructData.Position, *constructData.TileData, TileMap, PlScene, NewTransform, NewRenderer, assetsManager, constructData.Walktime, constructData.DirTag)()
	if err != nil {
		return nil, fmt.Errorf("NewPlayer() でエラーが発生しました。\n%w", err)
	}

	return plobj, nil
}
