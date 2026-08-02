package main

import (
	"fmt"
	"log"

	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// ゲーム起動
	initScene := PlayScene{} // 起動時に表示されるシーン
	initScene.Init()
	game := Game{ActiveScene: initScene}
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	ActiveScene Scene
	BackScene   Utils.Stack[Scene]
	// 中継用の画像バッファ(毎フレーム生成すると重い)
	bufferScreen *ebiten.Image
}

// Draw implements [ebiten.Game].
func (g Game) Draw(screen *ebiten.Image) {
	g.bufferScreen.Clear()

	// 奥から手前に向けて描写する
	for scene := range g.BackScene.Rev() {
		scene.Draw(g.bufferScreen)
	}
	g.ActiveScene.Draw(g.bufferScreen)

	// 中継画像バッファを画面に描画
	screen.DrawImage(g.bufferScreen, nil)
}

// Layout implements [ebiten.Game].
func (g Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}

// Update implements [ebiten.Game].
func (g Game) Update() error {
	nextScene, asNewScene, err := g.ActiveScene.Update(true)
	if err != nil {
		return fmt.Errorf("ActiveScene %s でエラーが発生しました \n %w", g.ActiveScene.Name(), err)
	}
	for scene := range g.BackScene.All() {
		ns, _, err := scene.Update(false)
		if err != nil {
			return fmt.Errorf("BackScene %s でエラーが発生しました \n %w", scene.Name(), err)
		}
		if ns != nil {
			return fmt.Errorf("BackScene %s がシーン切り替え要求を行いました。\n正常動作ならこのエラーを削除してください。", scene.Name())
		}
	}
	if nextScene != nil {
		if asNewScene {
			g.BackScene.Clear()
		} else {
			g.BackScene.Push(g.ActiveScene)
		}
		g.ActiveScene = nextScene
		if err := g.ActiveScene.Init(); err != nil {
			return fmt.Errorf("Scene %s の初期化中にエラーが発生しました \n %w", g.ActiveScene.Name(), err)
		}
	}
	return nil
}

type Scene interface {
	Name() string
	Init() error
	Update(active bool) (nextScene Scene, asNewScene bool, err error)
	Draw(screen *ebiten.Image)
}

// メインの遊べるシーン
type PlayScene struct {
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Init implements [Scene].
func (s PlayScene) Init() error {
	// なにもしない
	return nil
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	// なにもしない
}

// Update implements [Scene].
func (s PlayScene) Update(active bool) (nextScene Scene, asNewScene bool, err error) {
	// なにもしない
	return
}
