package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

func GameConstruct(assetsManager *assets.AssetsManager) Utils.Factory[*game.Game] {
	initScene, err := PlaySceneConstruct(assetsManager)() // 起動時に表示されるシーン
	if err != nil {
		return func() (*game.Game, error) {
			return nil, fmt.Errorf("初期シーンの読み込みに失敗しました\n%v", err)
		}
	}
	created, err := game.NewGame(initScene)()
	if err != nil {
		return func() (*game.Game, error) {
			return nil, fmt.Errorf("ゲームの構築に失敗しました。\n%v", err)
		}
	}

	return func() (*game.Game, error) {
		return created, nil
	}
}
