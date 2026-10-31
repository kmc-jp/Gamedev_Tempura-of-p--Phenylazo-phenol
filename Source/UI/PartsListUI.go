package ui

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"
	stdcolor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	partsListWidth  = 600
	partsListHeight = 300
	partCardHeight  = 100
	partCardGap     = 10
)

type PartsListUI struct {
	position    utils.Position
	font        *text.GoTextFace
	captionFont *text.GoTextFace
	scroll      *ScrollUI
	cards       []*PartCard
	selected    int
	visible     bool
	empty       *Text
}

func NewPartsListUI(position utils.Position, font *text.GoTextFace, captionFont *text.GoTextFace) (*PartsListUI, error) {
	if font == nil {
		return nil, fmt.Errorf("font が nil です")
	}

	scroll, err := NewScrollUI(
		image.Point{X: int(position.X) + 20, Y: int(position.Y) + 40},
		image.Point{X: partsListWidth - 40, Y: partsListHeight - 60},
		image.Point{X: partsListWidth - 40, Y: partsListHeight - 60},
		ScrollVertical,
	)
	if err != nil {
		return nil, fmt.Errorf("NewScrollUI() でエラーが発生しました: %w", err)
	}

	return &PartsListUI{position: position, font: font, captionFont: captionFont, scroll: scroll, selected: -1}, nil
}

func (p *PartsListUI) SetParts(parts game.Parts) error {
	contentHeight := partsListHeight - 60
	if len(parts) > 0 {
		contentHeight = len(parts)*(partCardHeight+partCardGap) - partCardGap
	}
	if contentHeight < partsListHeight-60 {
		contentHeight = partsListHeight - 60
	}

	scroll, err := NewScrollUI(
		image.Point{X: int(p.position.X) + 20, Y: int(p.position.Y) + 40},
		image.Point{X: partsListWidth - 40, Y: partsListHeight - 60},
		image.Point{X: partsListWidth - 40, Y: contentHeight},
		ScrollVertical,
	)
	if err != nil {
		return fmt.Errorf("NewScrollUI() でエラーが発生しました: %w", err)
	}

	p.cards = nil
	for i, part := range parts {
		card := NewPartCard(
			image.Point{X: 0, Y: i * (partCardHeight + partCardGap)},
			image.Point{X: partsListWidth - 40, Y: partCardHeight},
			part,
			p.font,
			p.captionFont,
		)
		_ = scroll.AddChild(card)
		p.cards = append(p.cards, card)
	}
	p.scroll = scroll
	p.selected = -1
	if len(p.cards) > 0 {
		p.selected = 0
		p.cards[0].SetSelected(true)
	}
	p.empty = nil
	if len(parts) == 0 {
		empty, err := NewText(
			utils.Position{X: p.position.X + 20, Y: p.position.Y + 80},
			gamecolor.Black(),
			"所持パーツはありません",
			p.font,
		)
		if err == nil {
			p.empty = empty
		}
	}
	return nil
}

func (p *PartsListUI) SetVisible(visible bool) {
	p.visible = visible
}

func (p *PartsListUI) Visible() bool {
	return p.visible
}

func (p *PartsListUI) Scroll(delta int) {
	p.scroll.Scroll(delta)
}

func (p *PartsListUI) MoveSelection(delta int) {
	if len(p.cards) == 0 || delta == 0 {
		return
	}

	next := p.selected + delta
	if next < 0 {
		next = 0
	}
	if next >= len(p.cards) {
		next = len(p.cards) - 1
	}
	if next == p.selected {
		return
	}

	if p.selected >= 0 {
		p.cards[p.selected].SetSelected(false)
	}
	p.selected = next
	p.cards[p.selected].SetSelected(true)
	p.ensureSelectedVisible()
}

func (p *PartsListUI) SelectedPart() *game.Part {
	if p.selected < 0 || p.selected >= len(p.cards) {
		return nil
	}
	return p.cards[p.selected].Part()
}

func (p *PartsListUI) SelectedIndex() int {
	return p.selected
}

func (p *PartsListUI) ensureSelectedVisible() {
	selectedTop := p.selected * (partCardHeight + partCardGap)
	selectedBottom := selectedTop + partCardHeight
	viewportHeight := partsListHeight - 60
	offset := p.scroll.Offset()
	if selectedTop < offset {
		p.scroll.SetOffset(selectedTop)
	} else if selectedBottom > offset+viewportHeight {
		p.scroll.SetOffset(selectedBottom - viewportHeight)
	}
}

func (p *PartsListUI) Update(active bool) error {
	if !p.visible {
		return nil
	}
	return p.scroll.Update(active)
}

func (p *PartsListUI) Draw(screen *ebiten.Image) {
	if !p.visible {
		return
	}

	dim := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	dim.Fill(stdcolor.RGBA{A: 160})
	screen.DrawImage(dim, nil)

	panel := ebiten.NewImage(partsListWidth, partsListHeight)
	panel.Fill(stdcolor.RGBA{R: 235, G: 235, B: 235, A: 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(panel, op)

	DrawText(
		utils.Position{X: p.position.X + 20, Y: p.position.Y + 18},
		gamecolor.Black(),
		"所持パーツ",
		p.font,
		screen,
	)
	if p.empty != nil {
		p.empty.Draw(screen)
		return
	}
	p.scroll.Draw(screen)
}

type PartCard struct {
	position    image.Point
	size        image.Point
	part        *game.Part
	font        *text.GoTextFace
	captionFont *text.GoTextFace
	selected    bool
}

func NewPartCard(position, size image.Point, part *game.Part, font *text.GoTextFace, captionFont *text.GoTextFace) *PartCard {
	return &PartCard{position: position, size: size, part: part, font: font, captionFont: captionFont}
}

func (p *PartCard) Update(active bool) error {
	return nil
}

func (p *PartCard) SetSelected(selected bool) {
	p.selected = selected
}

func (p *PartCard) Selected() bool {
	return p.selected
}

func (p *PartCard) Part() *game.Part {
	return p.part
}

func (p *PartCard) Draw(screen *ebiten.Image) {
	card := ebiten.NewImage(p.size.X, p.size.Y)
	card.Fill(stdcolor.White)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.position.X), float64(p.position.Y))
	screen.DrawImage(card, op)

	if p.selected {
		border := ebiten.NewImage(p.size.X, p.size.Y)
		border.Fill(stdcolor.RGBA{R: 255, G: 165, B: 0, A: 255})
		inner := ebiten.NewImage(p.size.X-6, p.size.Y-6)
		inner.Fill(stdcolor.White)
		borderOp := &ebiten.DrawImageOptions{}
		borderOp.GeoM.Translate(float64(p.position.X), float64(p.position.Y))
		screen.DrawImage(border, borderOp)
		innerOp := &ebiten.DrawImageOptions{}
		innerOp.GeoM.Translate(float64(p.position.X+3), float64(p.position.Y+3))
		screen.DrawImage(inner, innerOp)
	}

	if p.part.Img != nil {
		imageOp := &ebiten.DrawImageOptions{}
		imageOp.GeoM.Scale(2.2, 2.2)
		imageOp.GeoM.Translate(float64(p.position.X+10), float64(p.position.Y+10))
		screen.DrawImage(p.part.Img, imageOp)
	}

	x := float64(p.position.X + 100)
	DrawText(utils.Position{X: x, Y: float64(p.position.Y + 10)}, gamecolor.Black(), fmt.Sprintf("%s  ( %s )", p.part.Name, kindText(p.part)), p.font, screen)
	DrawText(utils.Position{X: x, Y: float64(p.position.Y + 34)}, gamecolor.Orange(), fmt.Sprintf("%s", rarityText(p.part)), p.font, screen)
	DrawText(utils.Position{X: x + 100, Y: float64(p.position.Y + 34)}, gamecolor.Black(), fmt.Sprintf("%s", rankText(p.part)), p.font, screen)
	DrawText(utils.Position{X: x, Y: float64(p.position.Y + 64)}, gamecolor.Black(), p.part.Detail, p.captionFont, screen)
}

func rarityText(part *game.Part) string {
	stars := ""
	for i := 0; i < int(part.Rarity); i++ {
		stars += "★"
	}
	return stars
}

func rankText(part *game.Part) string {
	switch part.Rank {
	case game.A:
		return "状態: 新品"
	case game.B:
		return "状態: 良"
	case game.C:
		return "状態: 普通"
	case game.D:
		return "ジャンク品"
	default:
		return "不明なパーツ状態です"
	}
}

func kindText(part *game.Part) string {
	switch part.Kind {
	case game.CPU:
		return "CPU"
	case game.GPU:
		return "GPU"
	case game.RAM:
		return "RAM"
	case game.MB:
		return "MB"
	case game.Case:
		return "Case"
	case game.LED:
		return "LED"
	default:
		return "その他"
	}
}

var _ scene.Element = (*PartsListUI)(nil)
