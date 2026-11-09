package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	resultui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/ResultUI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

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

		animText, err := ui.NewText(
			Utils.Position{},
			color.Orange(),
			"Enterキーでベンチマーク開始",
			ui.FontFace(gdata.Font, 20),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
		}
		animUI, err := ui.NewAnimatableUI(
			Utils.Position{},
			image.Point{100, 30},
			nil,
			animText,
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewanimatableUI() でエラーが発生しました。\n%w", err)
		}
		if err := animUI.SetPathAnimation(Utils.Move_Decelerate, Utils.Vec2{X: 150, Y: 150}, 360, 10); err != nil {
			return nil, err
		}
		resultscene.AddUIElement(animUI)

		InstalledPartsUI := resultui.NewInstalledParts(
			Utils.Position{X: 30, Y: 30},
			image.Point{X: 100, Y: 20},
			gdata,
		)
		resultscene.AddUIElement(InstalledPartsUI)

		partsList, err := ui.NewPartsListUI(
			Utils.Position{X: 20, Y: 30},
			ui.FontFace(gdata.Font, 16),
			ui.FontFace(gdata.Font, 12),
		)
		if err != nil {
			return nil, fmt.Errorf("ui.NewPartsListUI() でエラーが発生しました。\n%w", err)
		}
		resultscene.SetPartsList(partsList)

		return resultscene, nil
	}
}
