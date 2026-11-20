package ui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"
	stdcolor "image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	partsListWidth = 600
	// partsListGuideHeight は一覧パネル下端の操作ガイド欄の高さ。
	partsListGuideHeight = 25
	partsListHeight      = 325 + partsListGuideHeight
	partCardHeight       = 100
	// actionCardHeight は末尾のアクション項目（店を出る等）の高さ。
	actionCardHeight   = 50
	partCardGap        = 10
	scrollBarWidth     = 8
	partsListTabHeight = 32
	// emptyMessageHeight は一覧が空のときに「〜はありません」の表示用に確保する高さ。
	emptyMessageHeight = 40
	// partsListBorderWidth は選択中のタブ・パーツカードの枠の太さ。
	partsListBorderWidth = 3
	// guideFlashFrames は FlashGuide で操作ガイドを強調する長さ（フレーム数）。
	guideFlashFrames = 90
)

type PartsListTab uint8

const (
	PartsListFound PartsListTab = iota
	PartsListOwned
)

// listItem は一覧に並ぶ選択可能な項目（パーツカード・アクション項目）である。
type listItem interface {
	canvas.CanvasItem
	SetSelected(selected bool)
}

type PartsListUI struct {
	position      utils.Position
	font          *text.GoTextFace
	captionFont   *text.GoTextFace
	scroll        *ScrollUI
	cards         []*PartCard
	action        *ActionCard
	actionLabel   string
	contentHeight int
	// selected は項目の位置。0..len(cards)-1 がパーツ、len(cards) がアクション項目（ある場合）。
	selected   int
	lists      [2]game.Parts
	activeTab  PartsListTab
	showTabs   bool
	tabLabels  []string
	tabCount   int
	visible    bool
	empty      *Text
	guide      string
	guideFlash int
}

// NewPartsListUI は PartsListUI を生成する。
func NewPartsListUI(position utils.Position, font *text.GoTextFace, captionFont *text.GoTextFace) (*PartsListUI, error) {
	if font == nil {
		return nil, fmt.Errorf("font が nil です")
	}

	viewportHeight := partsListHeight - 60 - partsListGuideHeight
	scroll, err := NewScrollUI(
		image.Point{X: int(position.X) + 20, Y: int(position.Y) + 40},
		image.Point{X: partsListWidth - 40, Y: viewportHeight},
		image.Point{X: partsListWidth - 40, Y: viewportHeight},
		ScrollVertical,
	)
	if err != nil {
		return nil, fmt.Errorf("NewScrollUI() でエラーが発生しました: %w", err)
	}

	return &PartsListUI{
		position:      position,
		font:          font,
		captionFont:   captionFont,
		scroll:        scroll,
		contentHeight: viewportHeight,
		selected:      -1,
	}, nil
}

// ---- Sceneに公開する意味ベースAPI ----
//
// 選択・タブ切り替え・選択中パーツの取得（SelectNext / SelectPrevious / NextTab / PreviousTab / SelectedPart / IsActionSelected）は、
// 非表示中は受け付けない（何もしない、SelectedPart は nil、IsActionSelected は false を返す）。
// 一覧の構築（ShowFoundParts / ShowParts / RefreshFoundParts）は、非表示のまま構築してから Show できるよう、表示状態によらず行う。

// ShowFoundParts は発見したパーツ・所持パーツのタブ付き一覧を構築する。
// どちらのタブにも末尾に「店を出る」のアクション項目を置く。
// 表示状態は変更しないため、表示するには Show を呼ぶ。
func (p *PartsListUI) ShowFoundParts(found, owned game.Parts) error {
	p.showTabs = true
	p.tabLabels = []string{"発見したパーツ", "所持パーツ"}
	p.tabCount = len(p.tabLabels)
	p.actionLabel = "店を出る"
	p.lists[PartsListFound] = found
	p.lists[PartsListOwned] = owned
	p.activeTab = PartsListFound
	return p.rebuild(found, 0)
}

// ShowParts はタブなしの単一パーツ一覧を構築する。
// actionLabel が空でなければ、末尾にそのラベルのアクション項目を置く。
// 表示状態は変更しないため、表示するには Show を呼ぶ。
func (p *PartsListUI) ShowParts(parts game.Parts, actionLabel string) error {
	p.showTabs = false
	p.tabLabels = nil
	p.tabCount = 0
	p.actionLabel = actionLabel
	p.lists[PartsListFound] = parts
	p.activeTab = PartsListFound
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

// SetGuide は一覧パネル下端に表示する操作ガイドを設定する。強調中の状態は変えない。
func (p *PartsListUI) SetGuide(guide string) {
	p.guide = guide
}

// FlashGuide は操作ガイドを一定時間だけ強調色にして、押したキーの意味に気づかせる。
func (p *PartsListUI) FlashGuide() {
	p.guideFlash = guideFlashFrames
}

// SelectNext は選択を1つ下へ移す。非表示中は何もしない。
func (p *PartsListUI) SelectNext() {
	if !p.visible {
		return
	}
	p.moveSelection(1)
}

// SelectPrevious は選択を1つ上へ移す。非表示中は何もしない。
func (p *PartsListUI) SelectPrevious() {
	if !p.visible {
		return
	}
	p.moveSelection(-1)
}

// NextTab はタブを次へ切り替える。非表示中は何もしない。
func (p *PartsListUI) NextTab() error {
	if !p.visible {
		return nil
	}
	return p.moveTab(1)
}

// PreviousTab はタブを前へ切り替える。非表示中は何もしない。
func (p *PartsListUI) PreviousTab() error {
	if !p.visible {
		return nil
	}
	return p.moveTab(-1)
}

// SelectedPart は選択中のパーツを返す。非表示中・未選択・アクション項目の選択中は nil を返す。
func (p *PartsListUI) SelectedPart() *game.Part {
	if !p.visible || p.selected < 0 || p.selected >= len(p.cards) {
		return nil
	}
	return p.cards[p.selected].Part()
}

// IsActionSelected は末尾のアクション項目（店を出る等）を選択中かを返す。非表示中は false を返す。
func (p *PartsListUI) IsActionSelected() bool {
	return p.visible && p.actionSelected()
}

// IsShowingFoundParts はタブ付き一覧で「発見したパーツ」タブを表示中かを返す。
// ShowParts による単一一覧の表示中は false を返す。
func (p *PartsListUI) IsShowingFoundParts() bool {
	return p.showTabs && p.activeTab == PartsListFound
}

// ActiveTab は現在のタブを返す。表示状態やタブ表示の有無によらず、最後に選ばれたタブを返す。
func (p *PartsListUI) ActiveTab() PartsListTab {
	return p.activeTab
}

// RefreshFoundParts はタブ付き一覧のデータを差し替え、現在のタブと選択位置を維持して再構築する。
// 購入などで Found / Parts が変化したときに、Canvasが観測して呼ぶ。
// アクション項目を選んでいた場合は、再構築後もアクション項目を選んだままにする。
func (p *PartsListUI) RefreshFoundParts(found, owned game.Parts) error {
	p.lists[PartsListFound] = found
	p.lists[PartsListOwned] = owned
	selected := p.selected
	if p.actionSelected() {
		selected = math.MaxInt32
	}
	return p.rebuild(p.lists[p.activeTab], selected)
}

// ---- 内部実装 ----

func (p *PartsListUI) actionSelected() bool {
	return p.action != nil && p.selected == len(p.cards)
}

func (p *PartsListUI) moveTab(delta int) error {
	if !p.showTabs || delta == 0 {
		return nil
	}
	next := max(0, min(int(p.activeTab)+delta, p.tabCount-1))
	if next == int(p.activeTab) {
		return nil
	}
	p.activeTab = PartsListTab(next)
	return p.rebuild(p.lists[p.activeTab], 0)
}

// listTop は先頭の項目の上端（スクロール内座標）。一覧が空のときは空表示のための余白を確保する。
func (p *PartsListUI) listTop() int {
	if len(p.cards) == 0 {
		return emptyMessageHeight
	}
	return 0
}

// itemTop は項目 i の上端（スクロール内座標）。i == len(cards) はアクション項目。
func (p *PartsListUI) itemTop(i int) int {
	return p.listTop() + i*(partCardHeight+partCardGap)
}

func (p *PartsListUI) itemHeight(i int) int {
	if i < len(p.cards) {
		return partCardHeight
	}
	return actionCardHeight
}

func (p *PartsListUI) itemCount() int {
	if p.action != nil {
		return len(p.cards) + 1
	}
	return len(p.cards)
}

// rebuild はパーツカードとアクション項目を作り直し、selected を範囲内にクランプして選択する。
func (p *PartsListUI) rebuild(parts game.Parts, selected int) error {
	viewportHeight := p.viewportHeight()

	// 先にカードを作り、高さを確定させてからスクロール領域を作る
	p.cards = nil
	p.action = nil
	for _, part := range parts {
		price := fmt.Sprintf("価値: %d", part.Value())
		if p.showTabs && p.activeTab == PartsListFound {
			price = fmt.Sprintf("価格: %d", part.Val)
		}
		p.cards = append(p.cards, NewPartCard(
			image.Point{},
			image.Point{X: partsListWidth - 40, Y: partCardHeight},
			part,
			price,
			p.font,
			p.captionFont,
		))
	}
	if p.actionLabel != "" {
		p.action = NewActionCard(
			image.Point{X: 0, Y: p.itemTop(len(p.cards))},
			image.Point{X: partsListWidth - 40, Y: actionCardHeight},
			p.actionLabel,
			p.font,
		)
	}

	contentHeight := viewportHeight
	if n := p.itemCount(); n > 0 {
		last := n - 1
		contentHeight = max(contentHeight, p.itemTop(last)+p.itemHeight(last))
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

	p.contentHeight = contentHeight
	for i, card := range p.cards {
		card.position = image.Point{X: 0, Y: p.itemTop(i)}
		_ = scroll.AddChild(card)
	}
	if p.action != nil {
		_ = scroll.AddChild(p.action)
	}
	p.scroll = scroll

	p.selected = -1
	if n := p.itemCount(); n > 0 {
		p.selected = max(0, min(selected, n-1))
		p.item(p.selected).SetSelected(true)
		p.ensureSelectedVisible()
	}

	p.empty = nil
	if len(parts) == 0 {
		emptyText := "所持パーツはありません"
		if p.showTabs && p.activeTab == PartsListFound {
			emptyText = "発見したパーツはありません"
		}
		empty, err := NewText(
			utils.Position{X: p.position.X + 20, Y: p.position.Y + 48 + float64(p.tabOffset())},
			gamecolor.Black(),
			emptyText,
			p.font,
		)
		if err == nil {
			p.empty = empty
		}
	}
	return nil
}

// item は項目 i を返す。i は 0 以上 itemCount 未満であること。
func (p *PartsListUI) item(i int) listItem {
	if i < len(p.cards) {
		return p.cards[i]
	}
	return p.action
}

func (p *PartsListUI) tabOffset() int {
	if p.showTabs {
		return partsListTabHeight
	}
	return 0
}

func (p *PartsListUI) viewportHeight() int {
	return partsListHeight - 60 - partsListGuideHeight - p.tabOffset()
}

// Scroll はSceneにスクロール量を操作させない方針のため無効化している。
// スクロール位置は ensureSelectedVisible により PartsListUI 内部で完結させる。
// 将来ホイール等の相対スクロールやヘルパー実装で再利用する可能性があるため残す。
// func (p *PartsListUI) Scroll(delta int) {
// 	p.scroll.Scroll(delta)
// }

func (p *PartsListUI) moveSelection(delta int) {
	count := p.itemCount()
	if count == 0 || delta == 0 {
		return
	}

	next := max(0, min(p.selected+delta, count-1))
	if next == p.selected {
		return
	}

	if p.selected >= 0 {
		p.item(p.selected).SetSelected(false)
	}
	p.selected = next
	p.item(p.selected).SetSelected(true)
	p.ensureSelectedVisible()
}

func (p *PartsListUI) ensureSelectedVisible() {
	selectedTop := p.itemTop(p.selected)
	selectedBottom := selectedTop + p.itemHeight(p.selected)
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
	if p.guideFlash > 0 {
		p.guideFlash--
	}
	return p.scroll.Update(active)
}

func (p *PartsListUI) Draw(screen *ebiten.Image) {
	if !p.visible {
		return
	}

	// 背景を暗くする
	uiutils.FillRect(screen, screen.Bounds(), stdcolor.RGBA{A: 160})

	panelX, panelY := int(p.position.X), int(p.position.Y)
	uiutils.FillRect(
		screen,
		image.Rect(panelX, panelY, panelX+partsListWidth, panelY+partsListHeight),
		stdcolor.RGBA{R: 235, G: 235, B: 235, A: 255},
	)

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
	p.drawGuide(screen)
	if p.empty != nil {
		p.empty.Draw(screen)
	}
	p.scroll.Draw(screen)
	p.drawScrollBar(screen)
}

// drawGuide はパネル下端に操作ガイドを描画する。FlashGuide 中は強調色にする。
func (p *PartsListUI) drawGuide(screen *ebiten.Image) {
	if p.guide == "" {
		return
	}
	guideColor := gamecolor.Black()
	if p.guideFlash > 0 {
		guideColor = gamecolor.Orange()
	}
	DrawText(
		utils.Position{X: p.position.X + 20, Y: p.position.Y + partsListHeight - partsListGuideHeight - 14},
		guideColor,
		p.guide,
		p.font,
		screen,
	)
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

	uiutils.FillRect(
		screen,
		image.Rect(x, y, x+scrollBarWidth, y+trackHeight),
		stdcolor.RGBA{R: 190, G: 190, B: 190, A: 255},
	)
	uiutils.FillRect(
		screen,
		image.Rect(x, y+thumbOffset, x+scrollBarWidth, y+thumbOffset+thumbHeight),
		stdcolor.RGBA{R: 60, G: 60, B: 60, A: 255},
	)
}

func (p *PartsListUI) drawTabs(screen *ebiten.Image) {
	tabWidth := (partsListWidth - 40) / p.tabCount
	for i, label := range p.tabLabels {
		x := int(p.position.X) + 20 + i*tabWidth
		y := int(p.position.Y) + 35
		tabRect := image.Rect(x, y, x+tabWidth, y+partsListTabHeight)
		uiutils.FillRect(screen, tabRect, stdcolor.RGBA{R: 220, G: 220, B: 220, A: 255})
		if PartsListTab(i) == p.activeTab {
			uiutils.DrawBorder(screen, tabRect, partsListBorderWidth, uiutils.SelectionBorderColor)
		}
		DrawText(utils.Position{X: float64(x + 10), Y: float64(y + 2)}, gamecolor.Black(), label, p.font, screen)
	}
}

// ActionCard は一覧の末尾に置く「店を出る」等のアクション項目である。
type ActionCard struct {
	position image.Point
	size     image.Point
	label    string
	font     *text.GoTextFace
	selected bool
}

func NewActionCard(position, size image.Point, label string, font *text.GoTextFace) *ActionCard {
	return &ActionCard{position: position, size: size, label: label, font: font}
}

func (a *ActionCard) Update(active bool) error {
	return nil
}

func (a *ActionCard) SetSelected(selected bool) {
	a.selected = selected
}

func (a *ActionCard) Draw(screen *ebiten.Image) {
	rect := image.Rectangle{Min: a.position, Max: a.position.Add(a.size)}
	uiutils.FillRect(screen, rect, stdcolor.White)
	if a.selected {
		uiutils.DrawBorder(screen, rect, partsListBorderWidth, uiutils.SelectionBorderColor)
	}
	width, height := text.Measure(a.label, a.font, a.font.Size*1.5)
	DrawText(
		utils.Position{
			X: float64(a.position.X) + (float64(a.size.X)-width)/2,
			Y: float64(a.position.Y) + (float64(a.size.Y)-height)/2,
		},
		gamecolor.Black(),
		a.label,
		a.font,
		screen,
	)
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
	cardRect := image.Rectangle{Min: p.position, Max: p.position.Add(p.size)}
	uiutils.FillRect(screen, cardRect, stdcolor.White)
	if p.selected {
		uiutils.DrawBorder(screen, cardRect, partsListBorderWidth, uiutils.SelectionBorderColor)
	}

	if p.part.Img != nil {
		imageOp := &ebiten.DrawImageOptions{}
		imageOp.GeoM.Scale(p.part.PartsImgScaleOnPartsList(), p.part.PartsImgScaleOnPartsList())
		imageOp.GeoM.Translate(float64(p.position.X+10)+p.part.ImgOffset.X, float64(p.position.Y+10)+p.part.ImgOffset.Y)
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
