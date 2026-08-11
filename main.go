package main

import (
	"fmt"
	"log"
	"reflect"

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

type GameObjectUtils struct {
	gameObject scene.GameObject
}

func NewGameObjectUtils(obj scene.GameObject) (*GameObjectUtils, error) {
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
	var sceneTest game.Scene
	sceneTest, _ = NewPlayScene()
	sceneTest = sceneTest
	var GOTest GameObject
	GOTest, _ = NewHogeObject()
	GOTest = GOTest
	var componentTest Component
	componentTest, _ = NewHogeComponent()
	componentTest = componentTest
}
