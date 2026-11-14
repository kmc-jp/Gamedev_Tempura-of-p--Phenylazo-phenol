package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// PlayCanvas は PlayScene のUI構成である。
// 残り時間、空き重量・所持金、パーツ一覧を持つ。
// 一覧を開く時点はSceneが決めるため、Sceneのポインタを持って観測する。
type PlayCanvas struct {
	scene     *scene.PlayScene
	gdata     *game.Gamedata
	canvas    *canvas.Canvas
	partsList *ui.PartsListUI

	// shownSelectListRevision は反映済みの PlayScene.SelectListRevision。
	shownSelectListRevision int
	// shownPartsRevision は一覧に反映済みの Gamedata.PartsRevision。
	shownPartsRevision int
}

// NewPlayCanvas は PlayCanvas を生成する。
func NewPlayCanvas(s *scene.PlayScene, gdata *game.Gamedata) (*PlayCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("PlayScene が nil です")
	}
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

	return &PlayCanvas{
		scene:                   s,
		gdata:                   gdata,
		canvas:                  c,
		partsList:               partsList,
		shownSelectListRevision: s.SelectListRevision(),
		shownPartsRevision:      gdata.PartsRevision(),
	}, nil
}

// Update implements [scene.PlayCanvasAPI].
// TimeUp後に Play_Select へ変わったら（SelectListRevision の変化）、発見・所持パーツの一覧を表示する。
func (c *PlayCanvas) Update(active bool) error {
	if revision := c.scene.SelectListRevision(); revision != c.shownSelectListRevision {
		if err := c.openPartsList(); err != nil {
			return err
		}
		c.shownSelectListRevision = revision
	}

	// 購入や発見で Found / Parts が変化したら、表示中の一覧を再構築する
	if c.partsList.Visible() && c.gdata.PartsRevision() != c.shownPartsRevision {
		if err := c.partsList.RefreshFoundParts(c.gdata.FoundParts(), c.gdata.OwnedParts()); err != nil {
			return fmt.Errorf("パーツ一覧の更新に失敗しました。\n%w", err)
		}
		c.shownPartsRevision = c.gdata.PartsRevision()
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
// 一覧の非表示中は PartsListUI が何もしない。
func (c *PlayCanvas) MovePartSelectionDown() {
	c.partsList.SelectNext()
}

// MovePartSelectionUp implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) MovePartSelectionUp() {
	c.partsList.SelectPrevious()
}

// ShowNextPartsTab implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) ShowNextPartsTab() error {
	return c.partsList.NextTab()
}

// ShowPreviousPartsTab implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) ShowPreviousPartsTab() error {
	return c.partsList.PreviousTab()
}

// SelectedFoundPart implements [scene.PlayCanvasAPI].
// 一覧の非表示中、発見タブ以外の表示中、未選択のときは nil を返す。
// 非表示中は PartsListUI.SelectedPart が nil を返す。
func (c *PlayCanvas) SelectedFoundPart() *game.Part {
	if !c.partsList.IsShowingFoundParts() {
		return nil
	}
	return c.partsList.SelectedPart()
}

// IsPartsListOpen implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) IsPartsListOpen() bool {
	return c.partsList.Visible()
}

// CurrentPartsTab implements [scene.PlayCanvasAPI].
// UIのタブ型をSceneのタブ型へ変換して返す。一覧の非表示中は PartsListTabNone を返す。
func (c *PlayCanvas) CurrentPartsTab() scene.PartsListTab {
	if !c.partsList.Visible() {
		return scene.PartsListTabNone
	}
	switch c.partsList.ActiveTab() {
	case ui.PartsListFound:
		return scene.PartsListTabFound
	case ui.PartsListOwned:
		return scene.PartsListTabOwned
	}
	return scene.PartsListTabNone
}

func (c *PlayCanvas) openPartsList() error {
	if err := c.partsList.ShowFoundParts(c.gdata.FoundParts(), c.gdata.OwnedParts()); err != nil {
		return fmt.Errorf("パーツ一覧の構築に失敗しました。\n%w", err)
	}
	c.partsList.Show()
	c.shownPartsRevision = c.gdata.PartsRevision()
	return nil
}

var _ scene.PlayCanvasAPI = (*PlayCanvas)(nil)
