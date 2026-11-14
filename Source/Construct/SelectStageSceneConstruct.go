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

func SelectStageSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.SelectStageScene] {
	return func() (*scene.SelectStageScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		selectStageCanvas, err := scenecanvas.NewSelectStageCanvas(gdata.Font)
		if err != nil {
			return nil, fmt.Errorf("scenecanvas.NewSelectStageCanvas() でエラーが発生しました。\n%w", err)
		}

		selectStagescene, err := scene.NewSelectStageScene(DrawerFactory, selectStageCanvas, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewSelectStageScene() でエラーが発生しました。\n%w", err)
		}

		return selectStagescene, nil
	}
}
