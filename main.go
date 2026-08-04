package main

import (
	"fmt"
	"log"
	"reflect"

	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// ゲーム起動
	initScene := PlayScene{} // 起動時に表示されるシーン
	if err := initScene.Init(); err != nil {
		log.Fatalf("初期シーン %s の読み込みに失敗しました: %v", initScene.Name(), err)
	}
	game := Game{
		ActiveScene: &initScene,
	}
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	ActiveScene Scene
	BackScenes  Utils.Stack[Scene]
	// 中継用の画像バッファ(毎フレーム生成すると重い)
	bufferScreen *ebiten.Image
}

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
	nextScene, asNewScene, err := g.ActiveScene.Update(true)
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
	if nextScene != nil {
		if asNewScene {
			g.BackScenes.Clear()
		} else {
			g.BackScenes.Push(g.ActiveScene)
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
func (s *PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Init implements [Scene].
func (s *PlayScene) Init() error {
	// なにもしない
	return nil
}

// Draw implements [Scene].
func (s *PlayScene) Draw(screen *ebiten.Image) {
	// なにもしない
}

// Update implements [Scene].
func (s *PlayScene) Update(active bool) (nextScene Scene, asNewScene bool, err error) {
	// なにもしない
	return
}

type GameObject interface {
	transform() Transform
	render() Render
	components() []Component

	Name() string
	Update(active bool) (err error)
	Draw(screen *ebiten.Image)
}

type GameObjectUtils struct {
	obj *GameObject
}

func (u GameObjectUtils) GetComponent(target reflect.Type) (component Component, err error) {
	if u.obj == nil {
		return nil, fmt.Errorf("u.obj が未登録です")
	}
	for _, c := range (*u.obj).components() {
		if target == reflect.TypeOf(c) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%s に %s が存在しません。", (*u.obj).Name(), target.Name())
}

type Transform struct {
}

type Render struct {
}

func (r Render) Draw(obj GameObject, scene *ebiten.Image) {

}

type Component interface {
	Updatea(obj GameObject, active bool) (err error)
}
