package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	serializetarget "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize/SerializeTarget"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Component"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	scenecanvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/SceneCanvas"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/goaseprite"
)

// とりあえず PlayScene を組み立てる
func PlaySceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes, Stage2URL stageinfo.Stage2URL, stage stageinfo.Location) Utils.Factory[*scene.PlayScene] {
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

		// おじさんを設置 (Player の位置を避けるため Player の後に作る)
		ojisans, err := MakeOjisans(assetsManager, playscene, TM)
		if err != nil {
			return nil, fmt.Errorf("MakeOjisans(assetsManager, playscene, TM) でエラーが発生しました。\n%w", err)
		}
		for _, o := range ojisans {
			playscene.AddEntity(o)
		}

		return playscene, nil
	}
}

// ステージ類設定
func MakeStage2URL(assetsManager *assets.AssetsManager) (stageinfo.Stage2URL, error) {
	StagesSeriData, err := assetsManager.Serialize.NewSerializeData("Stages", reflect.TypeFor[serializetarget.StagesConfig](), "StageConfig.toml")
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.NewSerializeData(\"Stages\", reflect.TypeFor[serializetarget.StagesConfig](), \"StageConfig.toml\") でエラーが発生しました。\n%w", err)
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

	return StagesDataConstructed.(stageinfo.Stage2URL), nil
}

// タイルマップ設定
func MakeTileMapData(assetsManager *assets.AssetsManager, URL string) (*tilemap.TileMapData, error) {
	TileMapSeriData, err := assetsManager.Serialize.NewSerializeData("shelf", reflect.TypeFor[serializetarget.TileMapData](), URL)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.NewSerializeData(\"shelf\", reflect.TypeFor[serializetarget.TileMapData](), \"%s\") でエラーが発生しました。\n%w", URL, err)
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

func MakeTileMapObj(assetsManager *assets.AssetsManager, Stage2URL stageinfo.Stage2URL, stage stageinfo.Location, PlScene *scene.PlayScene) (*gameobject.GameObject, error) {
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
	SerializeData, err := assetsManager.Serialize.NewSerializeData("Player", reflect.TypeFor[serializetarget.PlayerConfig](), "Player.toml")
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.NewSerializeData(\"Player\", reflect.TypeFor[serializetarget.PlayerConfig](), \"Player.toml\") でエラーが発生しました。\n%w", err)
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

// おじさん設定
// 絵が届くまでは仮に暗くした主人公の絵を使う
func MakeOjisans(assetsManager *assets.AssetsManager, PlScene *scene.PlayScene, TileMap *component.TileMap) ([]*gameobject.GameObject, error) {
	SerializeData, err := assetsManager.Serialize.NewSerializeData("Ojisan", reflect.TypeFor[serializetarget.OjisanConfig](), "Ojisan.toml")
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.NewSerializeData(\"Ojisan\", reflect.TypeFor[serializetarget.OjisanConfig](), \"Ojisan.toml\") でエラーが発生しました。\n%w", err)
	}

	ConfigAny, err := assetsManager.Serialize.Load(*SerializeData)
	if err != nil {
		return nil, fmt.Errorf("assetsManager.Serialize.Load(*SerializeData) でエラーが発生しました。\n%w", err)
	}

	OjisanConfig := ConfigAny.(serializetarget.OjisanConfig)

	err = OjisanConfig.Init(assetsManager.Image)
	if err != nil {
		return nil, fmt.Errorf("OjisanConfig.Init(assetsManager.Image) でエラーが発生しました。\n%w", err)
	}

	constructDataAny, err := OjisanConfig.Construct()
	if err != nil {
		return nil, fmt.Errorf("OjisanConfig.Construct() でエラーが発生しました。\n%w", err)
	}
	constructData := constructDataAny.(serializetarget.OjisanConstructData)

	// 出現位置を決める
	spawns := ojisanSpawnPositions(TileMap, constructData)

	NewTransform := func(position Utils.Position) Utils.Factory[gameobject.Transform] {
		return Utils.CastFactory[*component.StdTransform, gameobject.Transform](component.NewStdTransform(position))
	}
	NewRenderer := func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[*component.AnimRender] {
		// 仮の見た目 : 主人公と区別できるよう暗い茶色にする (専用の絵が届いたら外す)
		dio := &ebiten.DrawImageOptions{}
		dio.ColorScale.Scale(0.75, 0.55, 0.4, 1)
		return component.NewAnimRender(asp, img, layer, "", component.AR_WithDrawImageOptions(dio))
	}

	ojisans := make([]*gameobject.GameObject, 0, len(spawns))
	for i, pos := range spawns {
		tpos := constructData.Position
		tpos.Position = pos
		obj, err := NewOjisan(
			fmt.Sprintf("Ojisan%d", i+1), tpos, *constructData.TileData, TileMap, PlScene,
			NewTransform, NewRenderer, assetsManager,
			constructData.Walktime, constructData.DirTag,
			constructData.MoveIntervalMin, constructData.MoveIntervalMax,
		)()
		if err != nil {
			return nil, fmt.Errorf("NewOjisan() でエラーが発生しました。\n%w", err)
		}
		ojisans = append(ojisans, obj)
	}

	return ojisans, nil
}

// 進入可能で誰もいないマスか
func isFreeTile(tm *component.TileMap, pos Utils.Vec2Int) bool {
	if enterable, ok := tm.Enterable[pos]; ok && !enterable {
		return false
	}
	return !tm.Occupied[pos]
}

// おじさんの出現位置を決める
// Spawns 指定があればそれを (進入不可・占有済みはスキップ)、なければ床からランダムに Count 箇所選ぶ
func ojisanSpawnPositions(tm *component.TileMap, data serializetarget.OjisanConstructData) []Utils.Vec2Int {
	result := []Utils.Vec2Int{}

	if len(data.Spawns) > 0 {
		used := map[Utils.Vec2Int]bool{}
		for _, pos := range data.Spawns {
			if !isFreeTile(tm, pos) || used[pos] {
				fmt.Printf("おじさんの出現位置 %v は進入不可、または他と重なっているためスキップしました。\n", pos)
				continue
			}
			used[pos] = true
			result = append(result, pos)
		}
		return result
	}

	// ランダム : 床タイルの座標を重複なく集める
	seen := map[Utils.Vec2Int]bool{}
	candidates := []Utils.Vec2Int{}
	for tpos := range tm.Tiles {
		pos := tpos.Position
		if seen[pos] {
			continue
		}
		seen[pos] = true
		if isFreeTile(tm, pos) {
			candidates = append(candidates, pos)
		}
	}
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	n := min(data.Count, len(candidates))
	return append(result, candidates[:n]...)
}
