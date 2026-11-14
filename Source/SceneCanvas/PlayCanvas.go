package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// PlayCanvas は PlayScene のUI構成である。
// 残り時間、空き重量・所持金、パーツ一覧を持つ。
// シーン固有の情報は不要なため、Sceneではなく Gamedata だけを参照する。
type PlayCanvas struct {
	gdata     *game.Gamedata
	canvas    *canvas.Canvas
	partsList *ui.PartsListUI

	// selectListShown は Play_Select への変化を検知して一覧を表示したか。
	selectListShown bool
	// shownFound・shownOwned は一覧に反映済みの Found / Parts。変化の観測に使う。
	shownFound game.Parts
	shownOwned game.Parts
}

// NewPlayCanvas は PlayCanvas を生成する。
func NewPlayCanvas(gdata *game.Gamedata) (*PlayCanvas, error) {
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}

	timeText, err := ui.NewGameDataText(
		utils.Position{X: 0, Y: 0},
		color.Red(),
		gdata,
		func(g game.Gamedata) string {
			return fmt.Sprintf("探索フェーズ 残り時間: %s", g.Text)
		},
		ui.FontFace(gdata.Font, 15),
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewGameDataText() でエラーが発生しました。\n%w", err)
	}

	gameDataUI, err := ui.NewGameDataUI(
		gdata,
		[]ui.GameDataUIItem{
			{
				Text: func(g game.Gamedata) string {
					return fmt.Sprintf("空き重量: %d", g.Weight)
				},
				Font:  ui.FontFace(gdata.Font, 15),
				Color: color.Black(),
			},
			{
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

	partsList, err := newPartsList(gdata)
	if err != nil {
		return nil, err
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{timeText, gameDataUI, partsList} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("プレイ画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	return &PlayCanvas{gdata: gdata, canvas: c, partsList: partsList}, nil
}

// Update implements [scene.PlayCanvasAPI].
// TimeUp後に Play_Select へ変わったら、発見・所持パーツの一覧を表示する。
func (c *PlayCanvas) Update(active bool) error {
	if c.gdata.State == game.Play_Select {
		if !c.selectListShown {
			if err := c.openPartsList(); err != nil {
				return err
			}
			c.selectListShown = true
		}
	} else {
		c.selectListShown = false
	}

	// 購入などで Found / Parts が変化したら、表示中の一覧を再構築する
	if c.partsList.Visible() &&
		(!slices.Equal(c.gdata.Found, c.shownFound) || !slices.Equal(c.gdata.Parts, c.shownOwned)) {
		if err := c.partsList.RefreshFoundParts(c.gdata.Found, c.gdata.Parts); err != nil {
			return fmt.Errorf("パーツ一覧の更新に失敗しました。\n%w", err)
		}
		c.shownFound = slices.Clone(c.gdata.Found)
		c.shownOwned = slices.Clone(c.gdata.Parts)
	}
	return c.canvas.Update(active)
}

// Draw implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsPlayCanvas implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) IsPlayCanvas() {}

// TogglePartsList implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) TogglePartsList() error {
	if c.partsList.Visible() {
		c.partsList.Hide()
		return nil
	}
	return c.openPartsList()
}

// MovePartSelectionDown implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) MovePartSelectionDown() {
	if c.partsList.Visible() {
		c.partsList.SelectNext()
	}
}

// MovePartSelectionUp implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) MovePartSelectionUp() {
	if c.partsList.Visible() {
		c.partsList.SelectPrevious()
	}
}

// ShowNextPartsTab implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) ShowNextPartsTab() error {
	if !c.partsList.Visible() {
		return nil
	}
	return c.partsList.NextTab()
}

// ShowPreviousPartsTab implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) ShowPreviousPartsTab() error {
	if !c.partsList.Visible() {
		return nil
	}
	return c.partsList.PreviousTab()
}

// SelectedFoundPart implements [scene.PlayCanvasAPI].
// 一覧の非表示中、発見タブ以外の表示中、未選択のときは nil を返す。
func (c *PlayCanvas) SelectedFoundPart() *game.Part {
	if !c.partsList.Visible() || !c.partsList.IsShowingFoundParts() {
		return nil
	}
	return c.partsList.SelectedPart()
}

// IsPartsListOpen implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) IsPartsListOpen() bool {
	return c.partsList.Visible()
}

func (c *PlayCanvas) openPartsList() error {
	if err := c.partsList.ShowFoundParts(c.gdata.Found, c.gdata.Parts); err != nil {
		return fmt.Errorf("パーツ一覧の構築に失敗しました。\n%w", err)
	}
	c.partsList.Show()
	c.shownFound = slices.Clone(c.gdata.Found)
	c.shownOwned = slices.Clone(c.gdata.Parts)
	return nil
}

var _ scene.PlayCanvasAPI = (*PlayCanvas)(nil)
