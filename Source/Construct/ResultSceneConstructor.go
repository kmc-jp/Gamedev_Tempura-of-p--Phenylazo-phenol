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

func ResultSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.ResultScene] {
	return func() (*scene.ResultScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		// CanvasはSceneを参照するため、Sceneの生成時に生成関数から作る
		newCanvas := func(s *scene.ResultScene) (scene.ResultCanvasAPI, error) {
			resultCanvas, err := scenecanvas.NewResultCanvas(s, gdata)
			if err != nil {
				return nil, fmt.Errorf("scenecanvas.NewResultCanvas() でエラーが発生しました。\n%w", err)
			}
			return resultCanvas, nil
		}

		resultscene, err := scene.NewResultScene(DrawerFactory, newCanvas, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewResultScene() でエラーが発生しました。\n%w", err)
		}

		return resultscene, nil
	}
}
