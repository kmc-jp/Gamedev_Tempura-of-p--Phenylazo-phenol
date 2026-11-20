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

// howToPlayImagePath は遊び方画面に表示する画像の /Assets/Images からの相対パス
const howToPlayImagePath = "遊び方.png"

func HowToPlaySceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes) Utils.Factory[*scene.HowToPlayScene] {
	return func() (*scene.HowToPlayScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeWorldDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		pageImage, err := assetsManager.Image.LoadPNG(howToPlayImagePath)
		if err != nil {
			return nil, fmt.Errorf("assetsManager.Image.LoadPNG(\"%s\") でエラーが発生しました。\n%w", howToPlayImagePath, err)
		}

		newCanvas := func(s *scene.HowToPlayScene) (scene.HowToPlayCanvasAPI, error) {
			return scenecanvas.NewHowToPlayCanvas(gdata, pageImage)
		}

		howToPlayScene, err := scene.NewHowToPlayScene(DrawerFactory, newCanvas, inputHandler, gdata, initScenes)()
		if err != nil {
			return nil, fmt.Errorf("NewHowToPlayScene() でエラーが発生しました。\n%w", err)
		}
		return howToPlayScene, nil
	}
}
