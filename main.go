package main

import (
	"log"

	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// ゲーム起動
	initScene, err := scene.NewPlayScene() // 起動時に表示されるシーン
	if err != nil {
		log.Fatalf("初期シーン %s の読み込みに失敗しました: %v", initScene.Name(), err)
	}
	game := game.Game{
		ActiveScene: &initScene,
	}
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
