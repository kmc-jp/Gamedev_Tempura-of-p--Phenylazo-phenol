package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	serialize "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize"
	serializetarget "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/Serialize/SerializeTarget"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"reflect"
)

func GameConstruct(assetsManager *assets.AssetsManager) Utils.Factory[*game.Game] {
	return func() (*game.Game, error) {
		inputConfigPath := "InputConfig.toml"
		inputConfigData, err := serialize.NewSerializeData("inputConfig", reflect.TypeFor[serializetarget.InputConfig](), inputConfigPath, assetsManager.Serialize)
		if err != nil {
			return nil, fmt.Errorf("serialize.NewSerializeData(\"inputConfig\", reflect.TypeFor[serializetarget.CameraConfig](), \"%s\", assetsManager.Serialize) でエラーが発生しました。\n%w", inputConfigPath, err)
		}

		inputConfigSerializedAny, err := assetsManager.Serialize.Load(*inputConfigData)
		if err != nil {
			return nil, fmt.Errorf("assetsManager.Serialize.Load(*inputConfigData) でエラーが発生しました。\n%w", err)
		}

		inputConfigAny, err := inputConfigSerializedAny.(serializetarget.InputConfig).Construct()
		if err != nil {
			return nil, fmt.Errorf("inputConfigSerializedAny.(serializetarget.InputConfig).Construct() でエラーが発生しました。\n%w", err)
		}

		Handler, err := inputConfigAny.(*inputconfig.InputConfig).Setup()
		if err != nil {
			return nil, fmt.Errorf("inputConfigAny.(inputconfig.InputConfig).Setup() でエラーが発生しました。\n%w", err)
		}

		gdata := game.NewGamedata()
		gdata.Font, err = ui.LoadFont("Assets/Fonts/Noto_Sans_JP/static/NotoSansJP-Medium.ttf")
		if err != nil {
			return nil, fmt.Errorf("ui.LoadFont(\"Assets/Fonts/Noto_Sans_JP/static/NotoSansJP-Medium.ttf\")でエラーが発生しました \n%w", err)
		}
		gdata.AssetsManager = assetsManager

		//各シーンの Constructor
		initializedScenes := &scene.InitializedScenes{}
		initializedScenes.InitTitleScene = TitleSceneConstruct(assetsManager, Handler, gdata, initializedScenes)
		initializedScenes.InitPlayScene = PlaySceneConstruct(assetsManager, Handler, gdata, initializedScenes)
		initializedScenes.InitPauseScene = PauseSceneConstruct(assetsManager, Handler, gdata, initializedScenes)
		initializedScenes.InitTrainScene = TrainSceneConstruct(assetsManager, Handler, gdata, initializedScenes)
		//(ここまで)各シーンの Constructor

		initScene, err := TitleSceneConstruct(assetsManager, Handler, gdata, initializedScenes)() // 起動時に表示されるシーン
		if err != nil {
			return nil, fmt.Errorf("初期シーンの読み込みに失敗しました\n%v", err)
		}

		created, err := game.NewGame(initScene, Handler)()
		if err != nil {
			return nil, fmt.Errorf("ゲームの構築に失敗しました。\n%v", err)
		}

		return created, nil
	}
}
