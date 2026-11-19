// Package pauseui は、PauseScene で使う固有のUI部品を提供する。
package pauseui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// PauseMenuConfig は PauseMenu の生成時設定である。
type PauseMenuConfig struct {
	Position              utils.Position
	Labels                []string // 上から順に並べるボタンの文言
	Font                  *text.GoTextFace
	TextColor             *gamecolor.Color
	ButtonBackgroundColor *gamecolor.Color
	ButtonSize            image.Point
	Spacing               int // ボタン間の縦の間隔
}

// PauseMenu は選択式ボタンを縦に並べる受動的なUIである。
// 入力は扱わず、選択中のボタンは SetSelected で外から渡す。
type PauseMenu struct {
	position              utils.Position
	buttons               []*ui.ButtonUI
	buttonSize            image.Point
	spacing               int
	buttonBackgroundColor *gamecolor.Color
}

// NewPauseMenu は設定に基づいて PauseMenu を生成する。
func NewPauseMenu(config PauseMenuConfig) (*PauseMenu, error) {
	if len(config.Labels) == 0 {
		return nil, fmt.Errorf("labels が空です")
	}
	if config.Font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	if config.TextColor == nil {
		return nil, fmt.Errorf("textColor が nil です")
	}
	if config.ButtonBackgroundColor == nil {
		return nil, fmt.Errorf("buttonBackgroundColor が nil です")
	}
	if config.ButtonSize.X <= 0 || config.ButtonSize.Y <= 0 {
		return nil, fmt.Errorf("buttonSize は 0 より大きい必要があります: %v", config.ButtonSize)
	}
	if config.Spacing < 0 {
		return nil, fmt.Errorf("spacing は 0 以上である必要があります: %d", config.Spacing)
	}

	m := &PauseMenu{
		position:              config.Position,
		buttonSize:            config.ButtonSize,
		spacing:               config.Spacing,
		buttonBackgroundColor: config.ButtonBackgroundColor,
	}
	for i, label := range config.Labels {
		pos := m.buttonPosition(i)
		t, err := ui.NewText(pos, config.TextColor, label, config.Font)
		if err != nil {
			return nil, fmt.Errorf("ボタン %q のテキスト生成に失敗しました: %w", label, err)
		}
		content, err := ui.NewRichUI(pos, config.ButtonSize, nil, t)
		if err != nil {
			return nil, fmt.Errorf("ボタン %q の生成に失敗しました: %w", label, err)
		}
		button, err := ui.NewButtonUI(content)
		if err != nil {
			return nil, fmt.Errorf("ボタン %q の選択UI生成に失敗しました: %w", label, err)
		}
		m.buttons = append(m.buttons, button)
	}
	m.SetSelected(0)
	return m, nil
}

// buttonPosition は i 番目のボタンの左上位置を返す。
func (m *PauseMenu) buttonPosition(i int) utils.Position {
	return utils.Position{
		X: m.position.X,
		Y: m.position.Y + float64(i*(m.buttonSize.Y+m.spacing)),
	}
}

// SetSelected は i 番目のボタンだけを選択状態にする。範囲外なら何もしない。
func (m *PauseMenu) SetSelected(i int) {
	if i < 0 || i >= len(m.buttons) {
		return
	}
	for j, b := range m.buttons {
		b.SetSelected(j == i)
	}
}

// Update は子要素を更新する。入力処理は行わない。
func (m *PauseMenu) Update(active bool) error {
	for _, b := range m.buttons {
		if err := b.Update(active); err != nil {
			return fmt.Errorf("ポーズメニューのボタン更新に失敗しました: %w", err)
		}
	}
	return nil
}

// Draw はボタンの背景、文言、選択枠を描画する。
func (m *PauseMenu) Draw(screen *ebiten.Image) {
	for _, b := range m.buttons {
		uiutils.FillRect(screen, b.RichUI().Bounds(), m.buttonBackgroundColor.ToRGBA())
		b.Draw(screen)
	}
}

// Position はメニューの左上位置を返す。
func (m *PauseMenu) Position() utils.Position {
	return m.position
}

// SetPosition はメニュー全体を移動する。
func (m *PauseMenu) SetPosition(position utils.Position) {
	m.position = position
	for i, b := range m.buttons {
		b.RichUI().SetPosition(m.buttonPosition(i))
	}
}

// Size はメニュー全体のサイズを返す。
func (m *PauseMenu) Size() image.Point {
	n := len(m.buttons)
	return image.Point{X: m.buttonSize.X, Y: n*m.buttonSize.Y + (n-1)*m.spacing}
}

var _ canvas.PositionItem = (*PauseMenu)(nil)
