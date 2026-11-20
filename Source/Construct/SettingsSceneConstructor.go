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

func SettingsSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.SettingsScene] {
	return func() (*scene.SettingsScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeWorldDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		newCanvas := func(s *scene.SettingsScene) (scene.SettingsCanvasAPI, error) {
			return scenecanvas.NewSettingsCanvas(s, gdata)
		}

		settingsScene, err := scene.NewSettingsScene(DrawerFactory, newCanvas, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewSettingsScene() でエラーが発生しました。\n%w", err)
		}

		return settingsScene, nil
	}
}
