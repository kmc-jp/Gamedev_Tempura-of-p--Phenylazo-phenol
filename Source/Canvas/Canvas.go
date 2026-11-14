// Package canvas は、個別UIの保持・更新・描画・表示状態を管理する汎用コンテナを提供する。
//
// canvas はSceneを知らない。Scene固有の構成（どのUIをどの順で持つか）は
// scenecanvas packageの各Canvasが決める。
package canvas

import (
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// CanvasItem はCanvasが管理する個別UIの共通抽象である。
type CanvasItem interface {
	Update(active bool) error
	Draw(screen *ebiten.Image)
}

// PositionItem は位置を持つCanvasItemである。
type PositionItem interface {
	CanvasItem
	Position() utils.Position
	SetPosition(utils.Position)
}

// Canvas は追加された順にCanvasItemを更新・描画する汎用UIコンテナである。
// 描画順は追加順であり、後に追加したItemほど手前に描画される。
type Canvas struct {
	items   []CanvasItem
	visible bool
}

// NewCanvas は表示状態の空のCanvasを生成する。
func NewCanvas() *Canvas {
	return &Canvas{
		items:   []CanvasItem{},
		visible: true,
	}
}

// Add はItemを末尾に追加する。
func (c *Canvas) Add(item CanvasItem) error {
	if item == nil {
		return fmt.Errorf("canvas に追加する item が nil です")
	}
	c.items = append(c.items, item)
	return nil
}

// Update は追加順に各Itemを更新する。非表示の間は何もしない。
func (c *Canvas) Update(active bool) error {
	if !c.visible {
		return nil
	}
	for _, item := range c.items {
		if err := item.Update(active); err != nil {
			return fmt.Errorf("canvas の item 更新に失敗しました: %w", err)
		}
	}
	return nil
}

// Draw は追加順に各Itemを描画する。非表示の間は何も描画しない。
func (c *Canvas) Draw(screen *ebiten.Image) {
	if !c.visible {
		return
	}
	for _, item := range c.items {
		item.Draw(screen)
	}
}

// Show はCanvasを表示状態にする。
func (c *Canvas) Show() {
	c.visible = true
}

// Hide はCanvasを非表示状態にする。
func (c *Canvas) Hide() {
	c.visible = false
}

// Visible はCanvasが表示状態かを返す。
func (c *Canvas) Visible() bool {
	return c.visible
}

var _ CanvasItem = (*Canvas)(nil)
