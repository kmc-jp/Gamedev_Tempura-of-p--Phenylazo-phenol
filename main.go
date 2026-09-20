package main

import (
	"log"

	construct "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Construct"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// ゲーム起動
	initScene, err := construct.PlaySceneConstruct()() // 起動時に表示されるシーン
	if err != nil {
		log.Fatalf("初期シーン %s の読み込みに失敗しました: %v", initScene.Name(), err)
	}
	game := game.Game{
		ActiveScene: initScene,
	}
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
