package scenecanvas

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	ranking "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Ranking"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	resultui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/ResultUI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// 以下はランキング画面の見た目の調整値。
const (
	rankingTitleMessage = "ランキング"
	rankingEmptyMessage = "まだ記録がありません"
	rankingHintMessage  = "↑↓：えらぶ　Enter：タイトルへ"
	rankingNoImage      = "画像を表示できません"

	rankingMargin     = 24
	rankingTitleY     = 12
	rankingListTop    = 56
	rankingRowHeight  = 30
	rankingDetailX    = 340
	rankingImageScale = 3
	rankingImageTop   = 48
	rankingPartsTop   = 248
	rankingPartsLine  = 19
	rankingHintBottom = 28

	rankingTitleSize = 24
	rankingRowSize   = 18
	rankingPartsSize = 14
)

// RankingCanvas は RankingScene のUI構成である。
// 左に順位一覧、右に選択中の記録のPC画像と構成を描く。
// 一覧の選択はSceneが持つため、Sceneのポインタを持って観測する。
type RankingCanvas struct {
	scene *scene.RankingScene

	titleFont *text.GoTextFace
	rowFont   *text.GoTextFace
	partsFont *text.GoTextFace
}

// NewRankingCanvas は RankingCanvas を生成する。
func NewRankingCanvas(s *scene.RankingScene, gdata *game.Gamedata) (*RankingCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("RankingScene が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}
	if gdata.Font == nil {
		return nil, fmt.Errorf("gdata.Font が nil です")
	}
	return &RankingCanvas{
		scene:     s,
		titleFont: ui.FontFace(gdata.Font, rankingTitleSize),
		rowFont:   ui.FontFace(gdata.Font, rankingRowSize),
		partsFont: ui.FontFace(gdata.Font, rankingPartsSize),
	}, nil
}

// Update implements [scene.RankingCanvasAPI].
// 表示は毎フレームSceneから読み取って描くため、保持する状態はない。
func (c *RankingCanvas) Update(active bool) error { return nil }

// Draw implements [scene.RankingCanvasAPI].
func (c *RankingCanvas) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	ui.DrawText(utils.Position{X: rankingMargin, Y: rankingTitleY}, color.Orange(), rankingTitleMessage, c.titleFont, screen)

	records := c.scene.Records()
	if len(records) == 0 {
		ui.DrawText(utils.Position{X: rankingMargin, Y: rankingListTop}, color.LightGray(), rankingEmptyMessage, c.rowFont, screen)
	}

	selected := c.scene.SelectedIndex()
	for i, e := range records {
		y := float64(rankingListTop + i*rankingRowHeight)
		rowColor := color.White()
		if i == c.scene.HighlightIndex() {
			rowColor = color.Gold()
		}
		if i == selected {
			ui.DrawText(utils.Position{X: rankingMargin - 18, Y: y}, rowColor, "▶", c.rowFont, screen)
		}
		ui.DrawText(utils.Position{X: rankingMargin, Y: y}, rowColor, rankingRowText(i, e), c.rowFont, screen)
	}

	if selected >= 0 && selected < len(records) {
		c.drawDetail(screen, records[selected])
	}

	hintY := float64(bounds.Dy() - rankingHintBottom)
	ui.DrawText(utils.Position{X: rankingMargin, Y: hintY}, color.LightGray(), rankingHintMessage, c.partsFont, screen)
}

// drawDetail は選択中の記録のPC画像と構成を描く。
func (c *RankingCanvas) drawDetail(screen *ebiten.Image, e ranking.Entry) {
	if img := c.scene.SelectedImage(); img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(rankingImageScale, rankingImageScale)
		op.GeoM.Translate(rankingDetailX, rankingImageTop)
		screen.DrawImage(img, op)
	} else {
		ui.DrawText(utils.Position{X: rankingDetailX, Y: rankingImageTop}, color.Gray(), rankingNoImage, c.partsFont, screen)
	}

	for i, p := range e.Parts {
		category := resultui.PartCategoryName(p.Kind)
		switch {
		case p.Kind == uint(game.GPU):
			category = "GPU" // 「グラフィックボード」は欄に収まらないため短縮する
		case category == "":
			category = "そのほか"
		}
		y := float64(rankingPartsTop + i*rankingPartsLine)
		ui.DrawText(utils.Position{X: rankingDetailX, Y: y}, color.LightGray(), category+"：", c.partsFont, screen)
		ui.DrawText(utils.Position{X: rankingDetailX + 110, Y: y}, color.White(), p.Name, c.partsFont, screen)
	}
}

// rankingRowText は一覧1行分の文字列を返す。
func rankingRowText(index int, e ranking.Entry) string {
	return fmt.Sprintf("%2d位  %6d pts  %s", index+1, e.Score, e.PlayedAt.Format("01/02"))
}

// IsRankingCanvas implements [scene.RankingCanvasAPI].
func (c *RankingCanvas) IsRankingCanvas() {}

var _ scene.RankingCanvasAPI = (*RankingCanvas)(nil)
