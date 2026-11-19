package resultui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ScoreDisplayConfig はスコア表示の見た目とアニメーション速度を設定する。
type ScoreDisplayConfig struct {
	Font          *text.GoTextFace
	TextColor     *gamecolor.Color
	ShuffleFrames int
	DigitFrames   int
}

// ScoreDisplay はスコアを桁ごとにシャッフル表示から確定表示へ変更する。
type ScoreDisplay struct {
	text          *ui.Text
	font          *text.GoTextFace
	center        utils.Position
	target        []byte
	display       []byte
	fixedCount    int
	frame         int
	shuffleFrames int
	digitFrames   int
	started       bool
	finished      bool
}

func NewScoreDisplay(config ScoreDisplayConfig) (*ScoreDisplay, error) {
	if config.Font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	if config.TextColor == nil {
		return nil, fmt.Errorf("textColor が nil です")
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
		text:          textUI,
		font:          config.Font,
		shuffleFrames: config.ShuffleFrames,
		digitFrames:   config.DigitFrames,
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
	s.text.SetText(string(s.display))
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
	s.text.SetText(string(s.display))
	s.syncTextPosition()
	return nil
}

func (s *ScoreDisplay) Draw(screen *ebiten.Image) {
	if s.started {
		s.text.Draw(screen)
	}
}

func (s *ScoreDisplay) syncTextPosition() {
	if s.font == nil || s.text == nil {
		return
	}
	value := s.textText()
	width, height := text.Measure(value, s.font, s.font.Size*1.5)
	s.text.SetPosition(utils.Position{
		X: s.center.X - width/2,
		Y: s.center.Y - height/2,
	})
}

func (s *ScoreDisplay) textText() string {
	if len(s.display) == 0 {
		return "0"
	}
	return string(s.display)
}

func randomDigit() byte {
	return byte('0' + rand.Intn(10))
}

var _ canvas.CanvasItem = (*ScoreDisplay)(nil)
