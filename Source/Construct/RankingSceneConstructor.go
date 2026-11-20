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

func RankingSceneConstruct(assetsManager *assets.AssetsManager, inputHandler *input.Handler, gdata *game.Gamedata, initScenes *scene.InitializedScenes, store scene.RankingStore) Utils.Factory[*scene.RankingScene] {
	return func() (*scene.RankingScene, error) {
		DrawerFactory, err := MakeWorldDrawerFactory(assetsManager)
		if err != nil {
			return nil, fmt.Errorf("MakeWorldDrawerFactory(assetsManager) でエラーが発生しました。\n%w", err)
		}

		newCanvas := func(s *scene.RankingScene) (scene.RankingCanvasAPI, error) {
			return scenecanvas.NewRankingCanvas(s, gdata)
		}

		rankingScene, err := scene.NewRankingScene(DrawerFactory, newCanvas, inputHandler, gdata, initScenes, store)()
		if err != nil {
			return nil, fmt.Errorf("NewRankingScene() でエラーが発生しました。\n%w", err)
		}

		return rankingScene, nil
	}
}
