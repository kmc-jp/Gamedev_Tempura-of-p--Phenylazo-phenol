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

		resultscene, err := scene.NewResultScene(DrawerFactory, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewResultScene() でエラーが発生しました。\n%w", err)
		}

		resultCanvas, err := scenecanvas.NewResultCanvas(resultscene, gdata)
		if err != nil {
			return nil, fmt.Errorf("scenecanvas.NewResultCanvas() でエラーが発生しました。\n%w", err)
		}
		resultscene.SetCanvas(resultCanvas)

		return resultscene, nil
	}
}
