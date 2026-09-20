package main

import (
	"embed"
	"log"

	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	construct "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Construct"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

	"github.com/hajimehoshi/ebiten/v2"
)

// 「Assets」フォルダ配下のすべてのファイルを埋め込む
//
//go:embed Assets/*
var assetsFS embed.FS

func main() {
	AssetsManager, err := assets.NewAssetsManager(&assetsFS)
	if err != nil {
		log.Fatal(err)
	}
	// ゲーム起動
	initScene, err := construct.PlaySceneConstruct(AssetsManager)() // 起動時に表示されるシーン
	if err != nil {
		log.Fatalf("初期シーン %s の読み込みに失敗しました\n%v", initScene.Name(), err)
	}
	game := game.Game{
		ActiveScene: initScene,
	}
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
