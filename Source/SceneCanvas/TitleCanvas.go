package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	titleui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/TitleUI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// 以下はタイトル画面の見た目の調整値。人の手で編集する。

// titleBackgroundScale は背景画像の描画倍率。BackGroundCamera.Draw の倍率と合わせる。
const titleBackgroundScale = 5

// titleMenuAnchors は、タイトル背景画像上のメニュー文字の位置（画像ピクセル座標）。
// カーソルの先端（三角形の右端・縦中央）がこの点を指す。
var titleMenuAnchors = map[scene.TitleMenuItem]utils.Position{
	scene.TitleMenuStart:     {X: 62, Y: 55}, // 仮値。実機で調整する
	scene.TitleMenuHowToPlay: {X: 62, Y: 65},
	scene.TitleMenuSettings:  {X: 62, Y: 72},
}

// titleCursorStyle はカーソルの色・大きさ・揺れ方。
func titleCursorStyle() titleui.CursorStyle {
	return titleui.DefaultCursorStyle()
}

// TitleCanvas は TitleScene のUI構成である。メニューの選択カーソルを持つ。
// 選択中の項目はSceneが持つため、Sceneのポインタを持って観測する。
type TitleCanvas struct {
	scene  *scene.TitleScene
	canvas *canvas.Canvas
	cursor *titleui.MenuCursor

	items         []scene.TitleMenuItem
	shownRevision int
}

// NewTitleCanvas は TitleCanvas を生成する。
func NewTitleCanvas(s *scene.TitleScene) (*TitleCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("TitleScene が nil です")
	}

	items := scene.TitleMenuItems()
	anchors := make([]utils.Position, 0, len(items))
	for _, item := range items {
		pos, ok := titleMenuAnchors[item]
		if !ok {
			return nil, fmt.Errorf("メニュー項目 %d のカーソル位置が titleMenuAnchors にありません", item)
		}
		anchors = append(anchors, pos)
	}

	cursor, err := titleui.NewMenuCursor(anchors, titleBackgroundScale, titleCursorStyle())
	if err != nil {
		return nil, fmt.Errorf("titleui.NewMenuCursor() でエラーが発生しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	if err := c.Add(cursor); err != nil {
		return nil, fmt.Errorf("タイトル画面のCanvas構築に失敗しました。\n%w", err)
	}

	tc := &TitleCanvas{scene: s, canvas: c, cursor: cursor, items: items}
	tc.syncWithScene()
	return tc, nil
}

// syncWithScene は Scene の選択中の項目をカーソルへ反映する。
func (c *TitleCanvas) syncWithScene() {
	selected := c.scene.SelectedMenuItem()
	for i, item := range c.items {
		if item == selected {
			c.cursor.SetSelected(i)
			break
		}
	}
	c.shownRevision = c.scene.SelectionRevision()
}

// Update implements [scene.TitleCanvasAPI].
func (c *TitleCanvas) Update(active bool) error {
	if c.scene.SelectionRevision() != c.shownRevision {
		c.syncWithScene()
	}
	return c.canvas.Update(active)
}

// Draw implements [scene.TitleCanvasAPI].
func (c *TitleCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsTitleCanvas implements [scene.TitleCanvasAPI].
func (c *TitleCanvas) IsTitleCanvas() {}

var _ scene.TitleCanvasAPI = (*TitleCanvas)(nil)
