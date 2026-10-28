package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	serializetarget "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize/SerializeTarget"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject/Component"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/goaseprite"
)

// とりあえず PlayScene を組み立てる
func PlaySceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.PlayScene] {
	return func() (*scene.PlayScene, error) {
		// カメラ設定
		WorldDrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		// UIカメラ
		UIDrawerFactory, err := MakeUIDrawerFactory()
		if err != nil {
			return nil, fmt.Errorf("MakeUIDrawerFactory() でエラーが発生しました。\n%w", err)
		}

		playscene, err := scene.NewPlayScene(WorldDrawerFactory, UIDrawerFactory, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewPlayScene() でエラーが発生しました。\n%w", err)
		}

		// UI
		text, err := ui.NewGameDataText(
			Utils.Position{X: 0, Y: 0},
			color.Red(),
			gdata,
			func(g game.Gamedata) string {
				return fmt.Sprintf("探索フェーズ 残り時間: %s", g.Text)
			},
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		playscene.AddUIElement(text)

		// タイルマップ
		tilemapobj, err := MakeTileMapObj(assetsManager, playscene)
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

// タイルマップ設定
func MakeTileMapData(assetsManager *assets.AssetsManager) (*tilemap.TileMapData, error) {
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

func MakeTileMapObj(assetsManager *assets.AssetsManager, PlScene *scene.PlayScene) (*gameobject.GameObject, error) {
	TileMapData, err := MakeTileMapData(assetsManager)
	if err != nil {
		return nil, fmt.Errorf("MakeTileMap(assetsManager) でエラーが発生しました。\n%w", err)
	}

	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return Utils.CastFactory[*component.StdTransform, gameobject.Transform](component.NewStdTransform(position))
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[gameobject.Render] {
		return Utils.CastFactory[*component.StdRender, gameobject.Render](component.NewLayerableRender(asp, img, layer))
	}
	tilemapobj, err := NewTileMapObject("TileMapObj", *TileMapData, PlScene, NewTransform, NewRenderer, assetsManager)()
	if err != nil {
		return nil, fmt.Errorf("NewTileMapObject() でエラーが発生しました。\n%w", err)
	}

	return tilemapobj, nil
}

// プレイヤー設定
func MakePlayer(assetsManager *assets.AssetsManager, PlScene *scene.PlayScene, TileMap *component.TileMap) (*gameobject.GameObject, error) {
	SerializeData, err := serialize.NewSerializeData("Player", reflect.TypeFor[serializetarget.PlayerConfig](), "Player.toml", assetsManager.Serialize)
	if err != nil {
		return nil, fmt.Errorf("serialize.NewSerializeData(\"Player\", reflect.TypeFor[serializetarget.PlayerConfig](), \"Player.toml\", assetsManager.Serialize) でエラーが発生しました。\n%w", err)
	}

	ConfigAny, err := assetsManager.Serialize.Load(*SerializeData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*SerializeData) でエラーが発生しました。\n%w", err)
	}

	PlayerConfig := ConfigAny.(serializetarget.PlayerConfig)

	err = PlayerConfig.Init(assetsManager)
	if err != nil {
		return nil, fmt.Errorf("PlayerConfig.Init(assetsManager) でエラーが発生しました。\n%w", err)
	}

	constructDataAny, err := PlayerConfig.Construct()
	if err != nil {
		return nil, fmt.Errorf("PlayerConfig.Construct() でエラーが発生しました。\n%w", err)
	}
	constructData := constructDataAny.(serializetarget.PlayerConstructData)

	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return Utils.CastFactory[*component.StdTransform, gameobject.Transform](component.NewStdTransform(position))
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[gameobject.Render] {
		return Utils.CastFactory[*component.StdRender, gameobject.Render](component.NewLayerableRender(asp, img, layer))
	}

	plobj, err := NewPlayer("Player", constructData.Position, *constructData.TileData, TileMap, PlScene, NewTransform, NewRenderer, assetsManager)()
	if err != nil {
		return nil, fmt.Errorf("NewPlayer() でエラーが発生しました。\n%w", err)
	}

	return plobj, nil
}
