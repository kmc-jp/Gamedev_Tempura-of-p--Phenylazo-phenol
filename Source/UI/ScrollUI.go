package ui

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScrollMode uint8

const (
	ScrollVertical ScrollMode = iota
	ScrollHorizontal
)

type ScrollUI struct {
	position image.Point
	viewport image.Point
	content  image.Point
	mode     ScrollMode
	offset   int
	children []Element
	buffer   *ebiten.Image
}

func NewScrollUI(
	position image.Point,
	viewport image.Point,
	content image.Point,
	mode ScrollMode,
) (*ScrollUI, error) {
	if viewport.X <= 0 || viewport.Y <= 0 {
		return nil, fmt.Errorf("viewport は正のサイズである必要があります: %v", viewport)
	}
	if content.X < viewport.X || content.Y < viewport.Y {
		return nil, fmt.Errorf("content は viewport 以上のサイズである必要があります")
	}
	if mode != ScrollVertical && mode != ScrollHorizontal {
		return nil, fmt.Errorf("未対応のスクロールモードです: %d", mode)
	}

	return &ScrollUI{
		position: position,
		viewport: viewport,
		content:  content,
		mode:     mode,
		children: []Element{},
	}, nil
}

func (s *ScrollUI) AddChild(child Element) error {
	if child == nil {
		return fmt.Errorf("child が nil です")
	}
	s.children = append(s.children, child)
	return nil
}

func (s *ScrollUI) Update(active bool) error {
	if !active {
		return nil
	}
	for _, child := range s.children {
		if err := child.Update(active); err != nil {
			return fmt.Errorf("スクロール UI の子要素更新に失敗しました: %w", err)
		}
	}
	return nil
}

func (s *ScrollUI) Draw(screen *ebiten.Image) {
	if s.buffer == nil || s.buffer.Bounds().Dx() != s.content.X || s.buffer.Bounds().Dy() != s.content.Y {
		s.buffer = ebiten.NewImage(s.content.X, s.content.Y)
	}
	s.buffer.Clear()
	for _, child := range s.children {
		child.Draw(s.buffer)
	}

	source := image.Rect(0, 0, s.viewport.X, s.viewport.Y)
	if s.mode == ScrollVertical {
		source = image.Rect(0, s.offset, s.viewport.X, s.offset+s.viewport.Y)
	} else {
		source = image.Rect(s.offset, 0, s.offset+s.viewport.X, s.viewport.Y)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.position.X), float64(s.position.Y))
	screen.DrawImage(s.buffer.SubImage(source).(*ebiten.Image), op)
}

func (s *ScrollUI) Scroll(delta int) {
	s.SetOffset(s.offset + delta)
}

func (s *ScrollUI) SetOffset(offset int) {
	maxOffset := s.maxOffset()
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	s.offset = offset
}

func (s *ScrollUI) Offset() int {
	return s.offset
}

func (s *ScrollUI) Mode() ScrollMode {
	return s.mode
}

func (s *ScrollUI) Bounds() image.Rectangle {
	return image.Rect(s.position.X, s.position.Y, s.position.X+s.viewport.X, s.position.Y+s.viewport.Y)
}

func (s *ScrollUI) maxOffset() int {
	if s.mode == ScrollVertical {
		return s.content.Y - s.viewport.Y
	}
	return s.content.X - s.viewport.X
}

var _ Element = (*ScrollUI)(nil)
