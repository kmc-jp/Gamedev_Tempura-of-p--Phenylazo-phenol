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

func TitleSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.TitleScene] {
	return func() (*scene.TitleScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		UIDrawerFactory, err := MakeUIDrawerFactory()
		if err != nil {
			return nil, fmt.Errorf("MakeUIDrawerFactory() でエラーが発生しました。\n%w", err)
		}

		titlescene, err := scene.NewTitleScene(DrawerFactory, UIDrawerFactory, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewTitleScene() でエラーが発生しました。\n%w", err)
		}

		text, err := ui.NewText(
			Utils.Position{},
			color.Orange(),
			"Title\nPress Space Key to Start",
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		titlescene.AddUIElement(text)

		return titlescene, nil
	}
}
