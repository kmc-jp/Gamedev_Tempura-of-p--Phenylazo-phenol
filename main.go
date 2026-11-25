package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"

	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	construct "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Construct"
	gamepkg "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"

	"github.com/hajimehoshi/ebiten/v2"
)

// 「Assets」フォルダ配下のすべてのファイルを埋め込む
//
//go:embed Assets/*
var assetsFS embed.FS

const release = false

func main() {
	var fullScreenMode bool

	if release {
		fullScreenMode = true
	} else {
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
	// ウィンドウモードの初期サイズは、論理解像度の整数倍のうちモニターに収まる最大のもの。
	// 表示スケール150%の 1080p（モニターの論理サイズ 約1280x720）などで、固定倍率だとはみ出すため。
	// Size() はタスクバーを含むモニター全体の大きさなので、余白はタスクバーとタイトルバーの分を見込んだ目安の値。
	// Monitor() は RunGame 前にメインスレッドで呼ぶ必要があり、nil もありうる。
	const marginW, marginH = 80, 120
	scale := 1
	if m := ebiten.Monitor(); m != nil {
		mw, mh := m.Size() // DIP 単位
		scale = max(1, min((mw-marginW)/gamepkg.ScreenWidth, (mh-marginH)/gamepkg.ScreenHeight))
	}
	ebiten.SetWindowSize(gamepkg.ScreenWidth*scale, gamepkg.ScreenHeight*scale)
	// リサイズ可能にし、拡大・縮小は Ebitengine に任せる。
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(fullScreenMode)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("ゲームがクラッシュしました。\n%w", err)
	}
}
