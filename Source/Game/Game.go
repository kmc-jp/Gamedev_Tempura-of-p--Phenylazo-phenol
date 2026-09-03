package game

import (
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	ActiveScene Scene
	BackScenes  Utils.Stack[Scene]
	// 中継用の画像バッファ(毎フレーム生成すると重い)
	bufferScreen *ebiten.Image
}

// 実装テスト
var _ ebiten.Game = &Game{}

// Draw implements [ebiten.Game].
func (g *Game) Draw(screen *ebiten.Image) {
	// バッファ画像の初期化
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if g.bufferScreen == nil || g.bufferScreen.Bounds().Dx() != w || g.bufferScreen.Bounds().Dy() != h {
		g.bufferScreen = ebiten.NewImage(w, h)
	}

	g.bufferScreen.Clear()

	// 奥から手前に向けて描写する
	for scene := range g.BackScenes.Rev() {
		scene.Draw(g.bufferScreen)
	}
	g.ActiveScene.Draw(g.bufferScreen)

	// 中継画像バッファを画面に描画
	screen.DrawImage(g.bufferScreen, nil)
}

// Layout implements [ebiten.Game].
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}

// Update implements [ebiten.Game].
func (g *Game) Update() error {
	nextScene, transitionType, err := g.ActiveScene.Update(true)
	if err != nil {
		return fmt.Errorf("ActiveScene %s でエラーが発生しました \n %w", g.ActiveScene.Name(), err)
	}
	for scene := range g.BackScenes.All() {
		ns, _, err := scene.Update(false)
		if err != nil {
			return fmt.Errorf("BackScene %s でエラーが発生しました \n %w", scene.Name(), err)
		}
		if ns != nil {
			return fmt.Errorf("BackScene %s がシーン切り替え要求を行いました。\n正常動作ならこのエラーを削除してください。", scene.Name())
		}
	}
	if transitionType == transition.None() {
		// 何もせず返す
	} else if transitionType == transition.Pop() {
		// 画面復帰
		newactive, ok := g.BackScenes.Pop()
		if !ok {
			return fmt.Errorf("背景シーンがないので読み込めません。")
		}
		g.ActiveScene = newactive
	} else {
		// 新規生成
		if transitionType == transition.NewScene() {
			g.BackScenes.Clear()
		} else if transitionType == transition.Push() {
			g.BackScenes.Push(g.ActiveScene)
		}
		g.ActiveScene, err = nextScene()
		if err != nil {
			return fmt.Errorf("Scene %s の初期化中にエラーが発生しました \n %w", g.ActiveScene.Name(), err)
		}
	}
	return nil
}
