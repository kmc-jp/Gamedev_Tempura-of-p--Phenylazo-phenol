package ui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"
	stdcolor "image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	partsListWidth     = 600
	partsListHeight    = 340
	partCardHeight     = 100
	partCardGap        = 10
	scrollBarWidth     = 8
	partsListTabHeight = 32
)

type PartsListTab uint8

const (
	PartsListFound PartsListTab = iota
	PartsListOwned
)

type PartsListUI struct {
	position      utils.Position
	gamedata      *game.Gamedata
	font          *text.GoTextFace
	captionFont   *text.GoTextFace
	scroll        *ScrollUI
	cards         []*PartCard
	contentHeight int
	selected      int
	lists         [2]game.Parts
	activeTab     PartsListTab
	showTabs      bool
	visible       bool
	empty         *Text
}

func NewPartsListUI(position utils.Position, gamedata *game.Gamedata, font *text.GoTextFace, captionFont *text.GoTextFace) (*PartsListUI, error) {
	if gamedata == nil {
		return nil, fmt.Errorf("gamedata が nil です")
	}
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

	return &PartsListUI{
		position:      position,
		gamedata:      gamedata,
		font:          font,
		captionFont:   captionFont,
		scroll:        scroll,
		contentHeight: partsListHeight - 60,
		selected:      -1,
	}, nil
}

// ---- Sceneに公開する意味ベースAPI ----

// ShowFoundParts は発見したパーツ・所持パーツのタブ付き一覧を構築する。
// 表示状態は変更しないため、表示するには Show を呼ぶ。
func (p *PartsListUI) ShowFoundParts(found, owned game.Parts) error {
	p.showTabs = true
	p.lists[PartsListFound] = found
	p.lists[PartsListOwned] = owned
	p.activeTab = PartsListFound
	return p.rebuild(found, 0)
}

// ShowParts はタブなしの単一パーツ一覧を構築する。
// 表示状態は変更しないため、表示するには Show を呼ぶ。
func (p *PartsListUI) ShowParts(parts game.Parts) error {
	p.showTabs = false
	p.lists[PartsListFound] = parts
	return p.rebuild(parts, 0)
}

func (p *PartsListUI) Show() {
	p.visible = true
}

func (p *PartsListUI) Hide() {
	p.visible = false
}

func (p *PartsListUI) Visible() bool {
	return p.visible
}

func (p *PartsListUI) SelectNext() {
	p.moveSelection(1)
}

func (p *PartsListUI) SelectPrevious() {
	p.moveSelection(-1)
}

func (p *PartsListUI) NextTab() error {
	return p.moveTab(1)
}

func (p *PartsListUI) PreviousTab() error {
	return p.moveTab(-1)
}

func (p *PartsListUI) SelectedPart() *game.Part {
	if p.selected < 0 || p.selected >= len(p.cards) {
		return nil
	}
	return p.cards[p.selected].Part()
}

// IsShowingFoundParts はタブ付き一覧で「発見したパーツ」タブを表示中かを返す。
// ShowParts による単一一覧の表示中は false を返す。
func (p *PartsListUI) IsShowingFoundParts() bool {
	return p.showTabs && p.activeTab == PartsListFound
}

// BuySelectedPart は選択中の発見済みパーツの購入をModelへ依頼する。
// 発見タブ以外の表示中や未選択時は game.NoItem を返す。
// 購入に成功した場合は Gamedata の Found / Parts で一覧を再構築し、選択位置を維持する。
func (p *PartsListUI) BuySelectedPart() (game.PurchaseResult, error) {
	if !p.IsShowingFoundParts() {
		return game.NoItem, nil
	}
	part := p.SelectedPart()
	if part == nil {
		return game.NoItem, nil
	}
	index := slices.Index(p.gamedata.Found, part)
	if index < 0 {
		return game.NoItem, nil
	}

	result := p.gamedata.BuyPart(index)
	if result != game.Buy {
		return result, nil
	}

	p.lists[PartsListFound] = p.gamedata.Found
	p.lists[PartsListOwned] = p.gamedata.Parts
	if err := p.rebuild(p.lists[PartsListFound], p.selected); err != nil {
		return result, fmt.Errorf("パーツ一覧の再構築に失敗しました: %w", err)
	}
	return result, nil
}

// ---- 互換API（フェーズB-8でSceneの呼び出しを置換後に削除する） ----

// Deprecated: ShowParts を使用する。
func (p *PartsListUI) SetParts(parts game.Parts) error {
	return p.ShowParts(parts)
}

// Deprecated: ShowFoundParts を使用する。
func (p *PartsListUI) SetLists(found, owned game.Parts) error {
	return p.ShowFoundParts(found, owned)
}

// Deprecated: Show / Hide を使用する。
func (p *PartsListUI) SetVisible(visible bool) {
	p.visible = visible
}

// Deprecated: NextTab / PreviousTab を使用する。
func (p *PartsListUI) MoveTab(delta int) error {
	return p.moveTab(delta)
}

// Deprecated: SelectNext / SelectPrevious を使用する。
func (p *PartsListUI) MoveSelection(delta int) {
	p.moveSelection(delta)
}

// Deprecated: 購入処理は BuySelectedPart を使用する。
func (p *PartsListUI) SelectedIndex() int {
	return p.selected
}

// Deprecated: IsShowingFoundParts を使用する。
func (p *PartsListUI) CurrentTab() PartsListTab {
	return p.activeTab
}

// Deprecated: IsShowingFoundParts を使用する。
// 単一一覧の表示中も true を返す点が IsShowingFoundParts と異なる。
func (p *PartsListUI) IsFoundTab() bool {
	return !p.showTabs || p.activeTab == PartsListFound
}

// ---- 内部実装 ----

func (p *PartsListUI) moveTab(delta int) error {
	if !p.showTabs || delta == 0 {
		return nil
	}
	next := int(p.activeTab) + delta
	if next < int(PartsListFound) {
		next = int(PartsListFound)
	}
	if next > int(PartsListOwned) {
		next = int(PartsListOwned)
	}
	if next == int(p.activeTab) {
		return nil
	}
	p.activeTab = PartsListTab(next)
	return p.rebuild(p.lists[p.activeTab], 0)
}

// rebuild はパーツカードを作り直し、selected を範囲内にクランプして選択する。
func (p *PartsListUI) rebuild(parts game.Parts, selected int) error {
	viewportHeight := p.viewportHeight()
	contentHeight := viewportHeight
	if len(parts) > 0 {
		contentHeight = len(parts) * (partCardHeight + partCardGap)
	}
	if contentHeight < viewportHeight {
		contentHeight = viewportHeight
	}

	scroll, err := NewScrollUI(
		image.Point{X: int(p.position.X) + 20, Y: int(p.position.Y) + 40 + p.tabOffset()},
		image.Point{X: partsListWidth - 40, Y: viewportHeight},
		image.Point{X: partsListWidth - 40, Y: contentHeight},
		ScrollVertical,
	)
	if err != nil {
		return fmt.Errorf("NewScrollUI() でエラーが発生しました: %w", err)
	}

	p.cards = nil
	p.contentHeight = contentHeight
	for i, part := range parts {
		price := fmt.Sprintf("価値: %d", part.Value())
		if p.showTabs && p.activeTab == PartsListFound {
			price = fmt.Sprintf("価格: %d", part.Val)
		}
		card := NewPartCard(
			image.Point{X: 0, Y: i * (partCardHeight + partCardGap)},
			image.Point{X: partsListWidth - 40, Y: partCardHeight},
			part,
			price,
			p.font,
			p.captionFont,
		)
		_ = scroll.AddChild(card)
		p.cards = append(p.cards, card)
	}
	p.scroll = scroll
	p.selected = -1
	if len(p.cards) > 0 {
		p.selected = max(0, min(selected, len(p.cards)-1))
		p.cards[p.selected].SetSelected(true)
		p.ensureSelectedVisible()
	}
	p.empty = nil
	if len(parts) == 0 {
		empty, err := NewText(
			utils.Position{X: p.position.X + 20, Y: p.position.Y + 80 + float64(p.tabOffset())},
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

func (p *PartsListUI) tabOffset() int {
	if p.showTabs {
		return partsListTabHeight
	}
	return 0
}

func (p *PartsListUI) viewportHeight() int {
	return partsListHeight - 60 - p.tabOffset()
}

// Scroll はSceneにスクロール量を操作させない方針のため無効化している。
// スクロール位置は ensureSelectedVisible により PartsListUI 内部で完結させる。
// 将来ホイール等の相対スクロールやヘルパー実装で再利用する可能性があるため残す。
// func (p *PartsListUI) Scroll(delta int) {
// 	p.scroll.Scroll(delta)
// }

func (p *PartsListUI) moveSelection(delta int) {
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

func (p *PartsListUI) ensureSelectedVisible() {
	selectedTop := p.selected * (partCardHeight + partCardGap)
	selectedBottom := selectedTop + partCardHeight
	viewportHeight := p.viewportHeight()
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
		utils.Position{X: p.position.X + 20, Y: p.position.Y + 12},
		gamecolor.Black(),
		"パーツ一覧",
		p.font,
		screen,
	)
	if p.showTabs {
		p.drawTabs(screen)
	}
	if p.empty != nil {
		p.empty.Draw(screen)
		return
	}
	p.scroll.Draw(screen)
	p.drawScrollBar(screen)
}

func (p *PartsListUI) drawScrollBar(screen *ebiten.Image) {
	viewportHeight := p.viewportHeight()
	maxOffset := p.contentHeight - viewportHeight
	if maxOffset <= 0 {
		return
	}

	trackHeight := viewportHeight
	thumbHeight := trackHeight * viewportHeight / p.contentHeight
	if thumbHeight < 20 {
		thumbHeight = 20
	}
	thumbTravel := trackHeight - thumbHeight
	thumbOffset := thumbTravel * p.scroll.Offset() / maxOffset
	x := int(p.position.X) + partsListWidth - 20
	y := int(p.position.Y) + 40 + p.tabOffset()

	track := ebiten.NewImage(scrollBarWidth, trackHeight)
	track.Fill(stdcolor.RGBA{R: 190, G: 190, B: 190, A: 255})
	trackOp := &ebiten.DrawImageOptions{}
	trackOp.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(track, trackOp)

	thumb := ebiten.NewImage(scrollBarWidth, thumbHeight)
	thumb.Fill(stdcolor.RGBA{R: 60, G: 60, B: 60, A: 255})
	thumbOp := &ebiten.DrawImageOptions{}
	thumbOp.GeoM.Translate(float64(x), float64(y+thumbOffset))
	screen.DrawImage(thumb, thumbOp)
}

func (p *PartsListUI) drawTabs(screen *ebiten.Image) {
	tabWidth := (partsListWidth - 40) / 2
	for i, label := range []string{"発見したパーツ", "所持パーツ"} {
		x := int(p.position.X) + 20 + i*tabWidth
		y := int(p.position.Y) + 35
		tab := ebiten.NewImage(tabWidth, partsListTabHeight)
		tab.Fill(stdcolor.RGBA{R: 220, G: 220, B: 220, A: 255})
		tabOp := &ebiten.DrawImageOptions{}
		tabOp.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(tab, tabOp)
		if PartsListTab(i) == p.activeTab {
			borderColor := stdcolor.RGBA{R: 255, G: 165, B: 0, A: 255}
			for _, edge := range []struct {
				x, y, width, height int
			}{
				{x, y, tabWidth, 3},
				{x, y + partsListTabHeight - 3, tabWidth, 3},
				{x, y, 3, partsListTabHeight},
				{x + tabWidth - 3, y, 3, partsListTabHeight},
			} {
				line := ebiten.NewImage(edge.width, edge.height)
				line.Fill(borderColor)
				lineOp := &ebiten.DrawImageOptions{}
				lineOp.GeoM.Translate(float64(edge.x), float64(edge.y))
				screen.DrawImage(line, lineOp)
			}
		}
		DrawText(utils.Position{X: float64(x + 10), Y: float64(y + 2)}, gamecolor.Black(), label, p.font, screen)
	}
}

type PartCard struct {
	position    image.Point
	size        image.Point
	part        *game.Part
	price       string
	font        *text.GoTextFace
	captionFont *text.GoTextFace
	selected    bool
}

func NewPartCard(position, size image.Point, part *game.Part, price string, font *text.GoTextFace, captionFont *text.GoTextFace) *PartCard {
	return &PartCard{position: position, size: size, part: part, price: price, font: font, captionFont: captionFont}
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
	DrawText(utils.Position{X: x + 100, Y: float64(p.position.Y + 34)}, gamecolor.Black(), fmt.Sprintf("%s   %s   重さ: %d", rankText(p.part), p.price, p.part.Weight), p.font, screen)
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

var _ canvas.CanvasItem = (*PartsListUI)(nil)
