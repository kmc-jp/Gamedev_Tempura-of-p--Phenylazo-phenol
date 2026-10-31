package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"

	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	construct "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Construct"

	"github.com/hajimehoshi/ebiten/v2"
)

// 「Assets」フォルダ配下のすべてのファイルを埋め込む
//
//go:embed Assets/*
var assetsFS embed.FS

func main() {
	var fullScreenMode bool
	fmt.Println("To run in FULL SCREEN mode, input y.\nOthers will be ignored.")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		text := scanner.Text()
		if text == "y" {
			fullScreenMode = true
		} else {
			fullScreenMode = false
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("読み込みエラー:", err)
	}

	AssetsManager, err := assets.NewAssetsManager(&assetsFS)
	if err != nil {
		log.Fatalf("assets.NewAssetsManager(&assetsFS) でエラーが発生しました。\n%v", err)
	}
	// ゲーム起動
	game, err := construct.GameConstruct(AssetsManager)()
	if err != nil {
		log.Fatalf("ゲームの読み込みに失敗しました。\n%v", err)
	}
	ebiten.SetFullscreen(fullScreenMode)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("ゲームがクラッシュしました。\n%w", err)
	}
}
