package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	serializetarget "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize/SerializeTarget"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	scenecanvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/SceneCanvas"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"reflect"

	input "github.com/quasilyte/ebitengine-input"
)

func SelectStageSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.SelectStageScene] {
	return func() (*scene.SelectStageScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		stageMap, err := MakeStageMap(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeStageMap(assetsManager) でエラーが発生しました。\n%w", err)
		}

		// CanvasはSceneを参照するため、Sceneの生成時に生成関数から作る
		newCanvas := func(s *scene.SelectStageScene) (scene.SelectStageCanvasAPI, error) {
			selectStageCanvas, err := scenecanvas.NewSelectStageCanvas(s, gdata)
			if err != nil {
				return nil, fmt.Errorf("scenecanvas.NewSelectStageCanvas() でエラーが発生しました。\n%w", err)
			}
			return selectStageCanvas, nil
		}

		selectStagescene, err := scene.NewSelectStageScene(DrawerFactory, stageMap, newCanvas, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewSelectStageScene() でエラーが発生しました。\n%w", err)
		}

		return selectStagescene, nil
	}
}

// MakeStageMap は行き先選択画面の路線図 (StageMap.toml) を読み込む。
func MakeStageMap(assetsManager *assets.AssetsManager) (stageinfo.StageMap, error) {
	seriData, err := assetsManager.Serialize.NewSerializeData("StageMap", reflect.TypeFor[serializetarget.StageMapConfig](), "StageMap.toml")
	if err != nil {
		return stageinfo.StageMap{}, fmt.Errorf("assetsManager.Serialize.NewSerializeData(\"StageMap\", ..., \"StageMap.toml\") でエラーが発生しました。\n%w", err)
	}
	loaded, err := assetsManager.Serialize.Load(*seriData)
	if err != nil {
		return stageinfo.StageMap{}, fmt.Errorf("assetsManager.Serialize.Load(*seriData) でエラーが発生しました。\n%w", err)
	}
	config, ok := loaded.(serializetarget.StageMapConfig)
	if !ok {
		return stageinfo.StageMap{}, fmt.Errorf("StageMap.toml の型が StageMapConfig ではありません: %T", loaded)
	}
	constructed, err := config.Construct()
	if err != nil {
		return stageinfo.StageMap{}, fmt.Errorf("config.Construct() でエラーが発生しました。\n%w", err)
	}
	stageMap, ok := constructed.(stageinfo.StageMap)
	if !ok {
		return stageinfo.StageMap{}, fmt.Errorf("Construct() の結果が StageMap ではありません: %T", constructed)
	}
	return stageMap, nil
}
