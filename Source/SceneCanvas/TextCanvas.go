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
// Title、PauseのCanvasで共通の構成である。
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
		ui.FontFace(gdata.Font, 16),
		ui.FontFace(gdata.Font, 12),
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewPartsListUI() でエラーが発生しました。\n%w", err)
	}
	return partsList, nil
}

// newGameDataUI は空き重量・所持金をアイコン付きで表示する共通UIを生成する。
func newGameDataUI(gdata *game.Gamedata) (*ui.GameDataUI, error) {
	if gdata.AssetsManager == nil {
		return nil, fmt.Errorf("gdata.AssetsManager が nil です")
	}

	weightIconData, err := gdata.AssetsManager.Image.NewImageData(
		"みんげー_リュックサック_アイコン",
		"UI/みんげー_リュックサック_アイコン.png",
		"UI/みんげー_リュックサック_アイコン.json",
	)
	if err != nil {
		return nil, fmt.Errorf("リュックサックアイコンのImageData作成に失敗しました。\n%w", err)
	}
	_, weightIcon, err := gdata.AssetsManager.Image.Load(*weightIconData)
	if err != nil {
		return nil, fmt.Errorf("リュックサックアイコンの読み込みに失敗しました。\n%w", err)
	}

	moneyIconData, err := gdata.AssetsManager.Image.NewImageData(
		"みんげー_がま口アイコン",
		"UI/みんげー_がま口アイコン.png",
		"UI/みんげー_がま口アイコン.json",
	)
	if err != nil {
		return nil, fmt.Errorf("がま口アイコンのImageData作成に失敗しました。\n%w", err)
	}
	_, moneyIcon, err := gdata.AssetsManager.Image.Load(*moneyIconData)
	if err != nil {
		return nil, fmt.Errorf("がま口アイコンの読み込みに失敗しました。\n%w", err)
	}

	gameDataUI, err := ui.NewGameDataUI(
		gdata,
		[]ui.GameDataUIItem{
			{
				Image: weightIcon,
				Text: func(g game.Gamedata) string {
					return fmt.Sprintf("空き重量: %d", g.Weight)
				},
				Font:  ui.FontFace(gdata.Font, 15),
				Color: color.Black(),
			},
			{
				Image: moneyIcon,
				Text: func(g game.Gamedata) string {
					return fmt.Sprintf("所持金: %d", g.Money)
				},
				Font:  ui.FontFace(gdata.Font, 15),
				Color: color.Black(),
			},
		},
		10,
		20,
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewGameDataUI() でエラーが発生しました。\n%w", err)
	}
	gameDataUI.SetBackgroundColor(color.White())
	return gameDataUI, nil
}
