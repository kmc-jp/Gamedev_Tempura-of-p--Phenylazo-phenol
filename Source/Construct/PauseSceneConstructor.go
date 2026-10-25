package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	input "github.com/quasilyte/ebitengine-input"
)

func PauseSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.PauseScene] {
	return func() (*scene.PauseScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		UIDrawerFactory, err := MakeUIDrawerFactory()
		if err != nil {
			return nil, fmt.Errorf("MakeUIDrawerFactory() でエラーが発生しました。\n%w", err)
		}

		pausescene, err := scene.NewPauseScene(DrawerFactory, UIDrawerFactory, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewPauseScene() でエラーが発生しました。\n%w", err)
		}

		text, err := ui.NewText(
			Utils.Position{X: 0, Y: 50},
			color.Orange(),
			"Pause\nPress ESC Key to Return the Game",
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		pausescene.AddUIElement(text)

		return pausescene, nil
	}
}
