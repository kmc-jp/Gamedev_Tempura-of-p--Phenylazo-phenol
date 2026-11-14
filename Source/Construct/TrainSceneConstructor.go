package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	scenecanvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/SceneCanvas"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	input "github.com/quasilyte/ebitengine-input"
)

func TrainSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.TrainScene] {
	return func() (*scene.TrainScene, error) {
		worldDrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		BGDrawerFactory, err := MakeBGDrawerFactory(assetsManager, "みんげー_電車内.png", "みんげー_電車内.json")

		// CanvasはSceneを参照するため、Sceneの生成時に生成関数から作る
		newCanvas := func(s *scene.TrainScene) (scene.TrainCanvasAPI, error) {
			trainCanvas, err := scenecanvas.NewTrainCanvas(s, gdata)
			if err != nil {
				return nil, fmt.Errorf("scenecanvas.NewTrainCanvas() でエラーが発生しました。\n%w", err)
			}
			return trainCanvas, nil
		}

		trainscene, err := scene.NewTrainScene(worldDrawerFactory, BGDrawerFactory, newCanvas, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewTrainScene() でエラーが発生しました。\n%w", err)
		}

		return trainscene, nil
	}
}
