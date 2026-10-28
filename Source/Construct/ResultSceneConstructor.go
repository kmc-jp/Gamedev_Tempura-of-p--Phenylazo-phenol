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

func ResultSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.ResultScene] {
	return func() (*scene.ResultScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		UIDrawerFactory, err := MakeUIDrawerFactory()
		if err != nil {
			return nil, fmt.Errorf("MakeUIDrawerFactory() でエラーが発生しました。\n%w", err)
		}

		resultscene, err := scene.NewResultScene(DrawerFactory, UIDrawerFactory, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewResultScene() でエラーが発生しました。\n%w", err)
		}

		text, err := ui.NewText(
			Utils.Position{},
			color.Orange(),
			"Result\nEnter キーでタイトルに戻る",
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		resultscene.AddUIElement(text)

		return resultscene, nil
	}
}
