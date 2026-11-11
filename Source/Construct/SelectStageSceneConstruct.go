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

func SelectStageSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.SelectStageScene] {
	return func() (*scene.SelectStageScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		UIDrawerFactory, err := MakeUIDrawerFactory()
		if err != nil {
			return nil, fmt.Errorf("MakeUIDrawerFactory() でエラーが発生しました。\n%w", err)
		}

		selectStagescene, err := scene.NewSelectStageScene(DrawerFactory, UIDrawerFactory, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewSelectStageScene() でエラーが発生しました。\n%w", err)
		}

		text, err := ui.NewText(
			Utils.Position{},
			color.Orange(),
			"ステージ選択\nPress 1 Key to Start Stage1\nPress 2 Key to Start RedStage\nPress 3 Key to Start BlueStage\nPress 4 Key to Start GreenStage\nPress 5 Key to Start OrangeStage\nPress 6 Key to Start PurpleStage",
			ui.FontFace(gdata.Font, 15),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		selectStagescene.AddUIElement(text)

		return selectStagescene, nil
	}
}
