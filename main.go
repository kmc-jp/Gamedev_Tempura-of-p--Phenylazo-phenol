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
	initScene, err := NewPlayScene() // 起動時に表示されるシーン
	if err != nil {
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
		g.ActiveScene, err = nextScene()
		if err != nil {
			return fmt.Errorf("Scene %s の初期化中にエラーが発生しました \n %w", g.ActiveScene.Name(), err)
		}
	}
	return nil
}

type Scene interface {
	Name() string
	Update(active bool) (nextScene SceneFactory, asNewScene bool, err error)
	Draw(screen *ebiten.Image)
}

type SceneFactory func() (Scene, error)

// メインの遊べるシーン
type PlayScene struct {
}

func NewPlayScene() (PlayScene, error) {
	ps := PlayScene{}
	// なにもしない
	return ps, nil
}

// Name implements [Scene].
func (s PlayScene) Name() string {
	return fmt.Sprintf("%T", s)
}

// Draw implements [Scene].
func (s PlayScene) Draw(screen *ebiten.Image) {
	// なにもしない
}

// Update implements [Scene].
func (s PlayScene) Update(active bool) (nextScene SceneFactory, asNewScene bool, err error) {
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
	gameObject GameObject
}

func NewGameObjectUtils(obj GameObject) (*GameObjectUtils, error) {
	if obj == nil {
		return nil, fmt.Errorf("obj が空です。 実装先の GameObject を代入してください。")
	}
	return &GameObjectUtils{gameObject: obj}, nil
}

func (u GameObjectUtils) GetComponent(target reflect.Type) (component Component, err error) {
	if u.gameObject == nil {
		return nil, fmt.Errorf("u.obj が未登録です")
	}
	for _, c := range u.gameObject.components() {
		if target == reflect.TypeOf(c) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%s に %s が存在しません。", u.gameObject.Name(), target.Name())
}

func (u GameObjectUtils) Update(active bool) (err error) {
	obj := u.gameObject
	for _, c := range obj.components() {
		if err := c.Update(obj, active); err != nil {
			return fmt.Errorf("コンポーネント %s でエラーが発生しました。\n%w", reflect.TypeOf(c), err)
		}
	}
	return nil
}

func (u GameObjectUtils) Draw(screen *ebiten.Image) {
	obj := u.gameObject
	obj.render().Draw(obj, screen)
}

type Transform struct {
}

func NewTransform() (*Transform, error) {
	return &Transform{}, nil
}

type Render struct {
}

func NewRender() (*Render, error) {
	return &Render{}, nil
}

func (r Render) Draw(obj GameObject, scene *ebiten.Image) {

}

type HogeObject struct {
	GameObjectUtils
	tf   Transform
	rd   Render
	cmps []Component
}

func NewHogeObject() (*HogeObject, error) {
	t, err := NewTransform()
	if err != nil {
		return nil, fmt.Errorf("NewTransform() でエラーが発生しました。\n %w", err)
	}
	r, err := NewRender()
	if err != nil {
		return nil, fmt.Errorf("NewRender() でエラーが発生しました。\n %w", err)
	}
	h := HogeObject{
		tf: *t,
		rd: *r,
	}
	h.GameObjectUtils.gameObject = h
	return &h, nil
}

func (h HogeObject) transform() Transform    { return h.tf }
func (h HogeObject) render() Render          { return h.rd }
func (h HogeObject) components() []Component { return h.cmps }
func (h HogeObject) Name() string            { return "Hoge" }

type Component interface {
	Update(obj GameObject, active bool) (err error)
}

type HogeComponent struct {
}

func (h HogeComponent) Update(obj GameObject, active bool) (err error) {
	// なにもしない
	return nil
}

func NewHogeComponent() (*HogeComponent, error) {
	return &HogeComponent{}, nil
}

func test() {
	var sceneTest Scene
	sceneTest, _ = NewPlayScene()
	sceneTest = sceneTest
	var GOTest GameObject
	GOTest, _ = NewHogeObject()
	GOTest = GOTest
	var componentTest Component
	componentTest, _ = NewHogeComponent()
	componentTest = componentTest
}
