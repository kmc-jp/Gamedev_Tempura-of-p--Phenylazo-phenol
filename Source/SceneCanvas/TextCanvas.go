// Package scenecanvas は、Sceneごとの具体的なCanvas実装を提供する。
//
// 各Canvasは scene package に定義された Canvas API interface を実装し、
// construct packageからSceneへ注入される。
// SceneはCanvas APIだけを利用し、このpackageの具体型や個別UIを知らない。
package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// newTextCanvas は、案内用の静的なテキスト1個だけを持つCanvasを生成する。
// Title、Pause、SelectStageのCanvasで共通の構成である。
func newTextCanvas(
	position utils.Position,
	value string,
	font *text.GoTextFaceSource,
	size float64,
) (*canvas.Canvas, error) {
	if font == nil {
		return nil, fmt.Errorf("フォントが nil です")
	}

	t, err := ui.NewText(position, color.Orange(), value, ui.FontFace(font, size))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	if err := c.Add(t); err != nil {
		return nil, fmt.Errorf("テキストのCanvasへの追加に失敗しました。\n%w", err)
	}
	return c, nil
}

// newPartsList は、Play・Train・Resultで共通の配置・フォントのパーツ一覧を生成する。
// 生成直後は非表示である。
func newPartsList(gdata *game.Gamedata) (*ui.PartsListUI, error) {
	partsList, err := ui.NewPartsListUI(
		utils.Position{X: 20, Y: 30},
		gdata,
		ui.FontFace(gdata.Font, 16),
		ui.FontFace(gdata.Font, 12),
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewPartsListUI() でエラーが発生しました。\n%w", err)
	}
	return partsList, nil
}
