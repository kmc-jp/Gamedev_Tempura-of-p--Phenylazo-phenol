package main

import (
	"embed"
	"log"

	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	construct "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Construct"

	"github.com/hajimehoshi/ebiten/v2"
)

// 「Assets」フォルダ配下のすべてのファイルを埋め込む
//
//go:embed Assets/*
var assetsFS embed.FS

func main() {
	AssetsManager, err := assets.NewAssetsManager(&assetsFS)
	if err != nil {
		log.Fatalf("assets.NewAssetsManager(&assetsFS) でエラーが発生しました。\n%v", err)
	}
	// ゲーム起動
	game, err := construct.GameConstruct(AssetsManager)()
	if err != nil {
		log.Fatalf("ゲームの読み込みに失敗しました。\n%v", err)
	}
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("ゲームがクラッシュしました。\n%w", err)
	}
}
