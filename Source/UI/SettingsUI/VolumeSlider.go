// Package settingsui は、SettingsScene で使う固有のUI部品を提供する。
package settingsui

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

const (
	sliderBorderWidth = 3
	sliderPadding     = 8
	sliderLabelWidth  = 80
	sliderCellWidth   = 14
	sliderCellGap     = 3
	sliderValueWidth  = 40
)

// VolumeSliderConfig は VolumeSlider の生成時設定である。
type VolumeSliderConfig struct {
	Position        utils.Position
	Label           string
	Max             int // 段階数（マスの数）
	Value           int
	Font            *text.GoTextFace
	TextColor       *gamecolor.Color
	FilledColor     *gamecolor.Color
	EmptyColor      *gamecolor.Color
	BackgroundColor *gamecolor.Color
	Height          int
}

// VolumeSlider は「ラベル + 段階バー + 数値」を横に並べた受動的なUIである。
// 入力は扱わず、値と選択状態は外から渡す。
type VolumeSlider struct {
	position   utils.Position
	label      *ui.Text
	valueText  *ui.Text
	max        int
	value      int
	selected   bool
	height     int
	filled     *gamecolor.Color
	empty      *gamecolor.Color
	background *gamecolor.Color
}

// NewVolumeSlider は設定に基づいて VolumeSlider を生成する。
func NewVolumeSlider(config VolumeSliderConfig) (*VolumeSlider, error) {
	if config.Max <= 0 {
		return nil, fmt.Errorf("max は 0 より大きい必要があります: %d", config.Max)
	}
	if config.Height <= 0 {
		return nil, fmt.Errorf("height は 0 より大きい必要があります: %d", config.Height)
	}
	if config.Font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	if config.TextColor == nil || config.FilledColor == nil || config.EmptyColor == nil || config.BackgroundColor == nil {
		return nil, fmt.Errorf("色の指定に nil が含まれています")
	}
	label, err := ui.NewText(utils.Position{}, config.TextColor, config.Label, config.Font)
	if err != nil {
		return nil, fmt.Errorf("ラベルの生成に失敗しました: %w", err)
	}
	valueText, err := ui.NewText(utils.Position{}, config.TextColor, "", config.Font)
	if err != nil {
		return nil, fmt.Errorf("数値の生成に失敗しました: %w", err)
	}
	s := &VolumeSlider{
		label:      label,
		valueText:  valueText,
		max:        config.Max,
		height:     config.Height,
		filled:     config.FilledColor,
		empty:      config.EmptyColor,
		background: config.BackgroundColor,
	}
	s.SetValue(config.Value)
	s.SetPosition(config.Position)
	return s, nil
}

// SetValue は表示する値を設定する（0〜max に丸める）。
func (s *VolumeSlider) SetValue(v int) {
	s.value = min(max(v, 0), s.max)
	s.valueText.SetText(fmt.Sprintf("%d", s.value))
}

// SetSelected は選択枠の表示を切り替える。
func (s *VolumeSlider) SetSelected(selected bool) { s.selected = selected }

func (s *VolumeSlider) barWidth() int {
	return s.max*sliderCellWidth + (s.max-1)*sliderCellGap
}

// Size は選択枠を含むUI全体のサイズを返す。
func (s *VolumeSlider) Size() image.Point {
	return image.Point{
		X: sliderPadding*2 + sliderLabelWidth + s.barWidth() + sliderValueWidth,
		Y: s.height,
	}
}

// Bounds は画面上の矩形を返す。
func (s *VolumeSlider) Bounds() image.Rectangle {
	size := s.Size()
	x, y := int(s.position.X), int(s.position.Y)
	return image.Rect(x, y, x+size.X, y+size.Y)
}

func (s *VolumeSlider) Update(active bool) error { return nil }

func (s *VolumeSlider) Draw(screen *ebiten.Image) {
	b := s.Bounds()
	uiutils.FillRect(screen, b, s.background.ToRGBA())

	barX := b.Min.X + sliderPadding + sliderLabelWidth
	cellH := s.height / 2
	cellY := b.Min.Y + (s.height-cellH)/2
	for i := 0; i < s.max; i++ {
		x := barX + i*(sliderCellWidth+sliderCellGap)
		c := s.empty
		if i < s.value {
			c = s.filled
		}
		uiutils.FillRect(screen, image.Rect(x, cellY, x+sliderCellWidth, cellY+cellH), c.ToRGBA())
	}

	s.label.Draw(screen)
	s.valueText.Draw(screen)

	if s.selected {
		uiutils.DrawBorder(screen, b, sliderBorderWidth, uiutils.SelectionBorderColor)
	}
}

func (s *VolumeSlider) Position() utils.Position { return s.position }

// SetPosition はUI全体を移動し、子要素の位置を更新する。
func (s *VolumeSlider) SetPosition(position utils.Position) {
	s.position = position
	b := s.Bounds()
	textY := position.Y + float64(s.height)/2 - 10
	s.label.SetPosition(utils.Position{X: float64(b.Min.X + sliderPadding), Y: textY})
	s.valueText.SetPosition(utils.Position{X: float64(b.Min.X + sliderPadding + sliderLabelWidth + s.barWidth() + 10), Y: textY})
}

var _ canvas.PositionItem = (*VolumeSlider)(nil)
