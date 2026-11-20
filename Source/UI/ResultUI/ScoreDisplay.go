package resultui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ScoreDisplayConfig はスコア表示の見た目とアニメーション速度を設定する。
type ScoreDisplayConfig struct {
	Font            *text.GoTextFace
	TextColor       *gamecolor.Color
	BackgroundColor *gamecolor.Color
	// 以下は省略可能。Label を指定すると数字の上に見出しを、Unit を指定すると数字の後ろに単位を表示する。
	BorderColor *gamecolor.Color
	Label       string
	LabelFont   *text.GoTextFace
	LabelColor  *gamecolor.Color
	Unit        string
	Padding     int
	ShuffleFrames   int
	DigitFrames     int
}

// ScoreDisplay はスコアを桁ごとにシャッフル表示から確定表示へ変更する。
type ScoreDisplay struct {
	text            *ui.Text
	font            *text.GoTextFace
	center          utils.Position
	backgroundColor *gamecolor.Color
	borderColor     *gamecolor.Color
	label           string
	labelFont       *text.GoTextFace
	labelColor      *gamecolor.Color
	unit            string
	padding         int
	target          []byte
	display         []byte
	fixedCount      int
	frame           int
	shuffleFrames   int
	digitFrames     int
	started         bool
	finished        bool
}

func NewScoreDisplay(config ScoreDisplayConfig) (*ScoreDisplay, error) {
	if config.Font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	if config.TextColor == nil {
		return nil, fmt.Errorf("textColor が nil です")
	}
	if config.Padding < 0 {
		config.Padding = 0
	}
	if config.ShuffleFrames <= 0 {
		config.ShuffleFrames = 2
	}
	if config.DigitFrames <= 0 {
		config.DigitFrames = 6
	}

	textUI, err := ui.NewText(utils.Position{}, config.TextColor, "", config.Font)
	if err != nil {
		return nil, fmt.Errorf("スコア表示テキストの生成に失敗しました: %w", err)
	}
	return &ScoreDisplay{
		text:            textUI,
		font:            config.Font,
		backgroundColor: config.BackgroundColor,
		borderColor:     config.BorderColor,
		label:           config.Label,
		labelFont:       config.LabelFont,
		labelColor:      config.LabelColor,
		unit:            config.Unit,
		padding:         config.Padding,
		shuffleFrames:   config.ShuffleFrames,
		digitFrames:     config.DigitFrames,
	}, nil
}

// SetCenter はスコアを中央配置する基準位置を設定する。
func (s *ScoreDisplay) SetCenter(center utils.Position) {
	s.center = center
	s.syncTextPosition()
}

// Start はスコアを確定し、シャッフルアニメーションを開始する。
func (s *ScoreDisplay) Start(score int) {
	if score < 0 {
		score = 0
	}
	s.target = []byte(strconv.Itoa(score))
	s.display = make([]byte, len(s.target))
	for i := range s.display {
		s.display[i] = randomDigit()
	}
	s.fixedCount = 0
	s.frame = 0
	s.started = true
	s.finished = false
	s.text.SetText(s.textText())
	s.syncTextPosition()
}

// IsStarted はアニメーション開始済みかを返す。
func (s *ScoreDisplay) IsStarted() bool {
	return s.started
}

// IsFinished は全桁の確定表示が完了したかを返す。
func (s *ScoreDisplay) IsFinished() bool {
	return s.started && s.finished
}

func (s *ScoreDisplay) Update(active bool) error {
	if !active || !s.started || s.finished {
		return nil
	}
	s.frame++

	if s.frame%s.shuffleFrames == 0 {
		for i := len(s.display) - 1 - s.fixedCount; i >= 0; i-- {
			s.display[i] = randomDigit()
		}
	}

	if s.frame%s.digitFrames == 0 && s.fixedCount < len(s.target) {
		index := len(s.target) - 1 - s.fixedCount
		s.display[index] = s.target[index]
		s.fixedCount++
	}

	if s.fixedCount == len(s.target) {
		copy(s.display, s.target)
		s.finished = true
	}
	s.text.SetText(s.textText())
	s.syncTextPosition()
	return nil
}

func (s *ScoreDisplay) Draw(screen *ebiten.Image) {
	if !s.started {
		return
	}
	panel := s.panelBounds()
	if s.backgroundColor != nil {
		uiutils.FillRect(screen, panel, s.backgroundColor.ToRGBA())
	}
	if s.borderColor != nil {
		const t = 3
		c := s.borderColor.ToRGBA()
		uiutils.FillRect(screen, image.Rect(panel.Min.X, panel.Min.Y, panel.Max.X, panel.Min.Y+t), c)
		uiutils.FillRect(screen, image.Rect(panel.Min.X, panel.Max.Y-t, panel.Max.X, panel.Max.Y), c)
		uiutils.FillRect(screen, image.Rect(panel.Min.X, panel.Min.Y, panel.Min.X+t, panel.Max.Y), c)
		uiutils.FillRect(screen, image.Rect(panel.Max.X-t, panel.Min.Y, panel.Max.X, panel.Max.Y), c)
	}
	if s.hasLabel() {
		lw, _ := text.Measure(s.label, s.labelFont, s.labelFont.Size*1.5)
		_, _, labelTop, _ := s.layout()
		ui.DrawText(utils.Position{X: s.center.X - lw/2, Y: labelTop}, s.labelColor, s.label, s.labelFont, screen)
	}
	s.text.Draw(screen)
}

func (s *ScoreDisplay) hasLabel() bool {
	return s.label != "" && s.labelFont != nil && s.labelColor != nil
}

// layout は見出し込みの全体サイズと、見出し・数字それぞれの上端Yを返す。
func (s *ScoreDisplay) layout() (width, height, labelTop, numberTop float64) {
	numW, numH := text.Measure(s.textText(), s.font, s.font.Size*1.5)
	width, height = numW, numH
	numberOffset := 0.0
	if s.hasLabel() {
		lw, lh := text.Measure(s.label, s.labelFont, s.labelFont.Size*1.5)
		if lw > width {
			width = lw
		}
		height += lh
		numberOffset = lh
	}
	top := s.center.Y - height/2
	return width, height, top, top + numberOffset
}

func (s *ScoreDisplay) panelBounds() image.Rectangle {
	width, height, _, _ := s.layout()
	return image.Rect(
		int(s.center.X-width/2)-s.padding,
		int(s.center.Y-height/2)-s.padding,
		int(s.center.X+width/2)+s.padding,
		int(s.center.Y+height/2)+s.padding,
	)
}

func (s *ScoreDisplay) syncTextPosition() {
	if s.font == nil || s.text == nil {
		return
	}
	numW, _ := text.Measure(s.textText(), s.font, s.font.Size*1.5)
	_, _, _, numberTop := s.layout()
	s.text.SetPosition(utils.Position{
		X: s.center.X - numW/2,
		Y: numberTop,
	})
}

// textText は数字と単位を連結した表示文字列を返す。
func (s *ScoreDisplay) textText() string {
	if len(s.display) == 0 {
		return "0" + s.unit
	}
	return string(s.display) + s.unit
}

func randomDigit() byte {
	return byte('0' + rand.Intn(10))
}

var _ canvas.CanvasItem = (*ScoreDisplay)(nil)
