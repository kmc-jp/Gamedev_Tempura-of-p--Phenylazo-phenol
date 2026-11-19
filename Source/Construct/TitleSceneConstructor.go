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

func TitleSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.TitleScene] {
	return func() (*scene.TitleScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		BGDrawerFactory, err := MakeBGDrawerFactory(assetsManager, "Backgrounds/みんげー_タイトル画面.png", "Backgrounds/みんげー_タイトル画面.json")
		if err != nil {
			return nil, fmt.Errorf("MakeBGDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		newCanvas := func(s *scene.TitleScene) (scene.TitleCanvasAPI, error) {
			return scenecanvas.NewTitleCanvas(s)
		}

		titlescene, err := scene.NewTitleScene(DrawerFactory, BGDrawerFactory, newCanvas, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewTitleScene() でエラーが発生しました。\n%w", err)
		}

		return titlescene, nil
	}
}
