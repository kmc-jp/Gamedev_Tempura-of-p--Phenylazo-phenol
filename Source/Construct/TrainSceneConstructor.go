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

func TrainSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.TrainScene] {
	return func() (*scene.TrainScene, error) {
		worldDrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		UIDrawerFactory, err := MakeUIDrawerFactory()
		if err != nil {
			return nil, fmt.Errorf("MakeUIDrawerFactory() でエラーが発生しました。\n%w", err)
		}

		trainscene, err := scene.NewTrainScene(worldDrawerFactory, UIDrawerFactory, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewTrainScene() でエラーが発生しました。\n%w", err)
		}

		text, err := ui.NewText(
			Utils.Position{},
			color.Orange(),
			"電車で移動中\nESCキーでメニューを開く\nEnterキー で Result へ",
			ui.FontFace(gdata.Font, 15),
		)

		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		trainscene.AddUIElement(text)

		partsList, err := ui.NewPartsListUI(
			Utils.Position{X: 20, Y: 30},
			ui.FontFace(gdata.Font, 16),
			ui.FontFace(gdata.Font, 12),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewPartsListUI() でエラーが発生しました。\n%w", err)
		}
		trainscene.SetPartsList(partsList)

		return trainscene, nil
	}
}
