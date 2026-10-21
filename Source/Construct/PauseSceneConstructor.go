package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	input "github.com/quasilyte/ebitengine-input"
)

func PauseSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.PauseScene] {
	DrawerFactory, err := MakeDrawerFactory(assetsManager)
	if err != nil {
		return func() (*scene.PauseScene, error) {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}
	}

	pausescene, err := scene.NewPauseScene(DrawerFactory, inputHandler, gdata, initScenes)()
	if err != nil {
		return func() (*scene.PauseScene, error) {
			return nil, fmt.Errorf("NewPauseScene() でエラーが発生しました。\n%w", err)
		}
	}

	return func() (*scene.PauseScene, error) {
		return pausescene, nil
	}
}
