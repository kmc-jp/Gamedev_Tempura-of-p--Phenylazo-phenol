package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	resultui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/ResultUI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// ResultCanvas は ResultScene のUI構成である。
// 案内テキスト、ベンチマーク案内、組み立て済みパーツ、パーツ一覧を持つ。
// 選択候補はSceneのゲームルールで決まるため、Sceneのポインタを持って観測する。
type ResultCanvas struct {
	scene *scene.ResultScene
	gdata *game.Gamedata

	// canvas は描画順の管理だけに使う。更新は状態に応じて個別に行う。
	canvas          *canvas.Canvas
	benchmarkPrompt *ui.AnimatableUI
	installedParts  *resultui.InstalledParts
	partsList       *ui.PartsListUI

	// shownRevision は一覧に反映済みの ResultScene.CandidatesRevision。
	shownRevision int
	// installedStarted は組み立て済みパーツの表示を開始したか。
	installedStarted bool
}

// NewResultCanvas は ResultCanvas を生成する。
func NewResultCanvas(s *scene.ResultScene, gdata *game.Gamedata) (*ResultCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("ResultScene が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}

	guide, err := ui.NewText(
		utils.Position{},
		color.Orange(),
		"Result\nEnter キーでタイトルに戻る",
		ui.FontFace(gdata.Font, 15),
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
	}

	promptText, err := ui.NewText(
		utils.Position{},
		color.Orange(),
		"Enterキーでベンチマーク開始",
		ui.FontFace(gdata.Font, 20),
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
	}
	benchmarkPrompt, err := ui.NewAnimatableUI(
		utils.Position{},
		image.Point{100, 30},
		nil,
		promptText,
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewAnimatableUI() でエラーが発生しました。\n%w", err)
	}
	if err := benchmarkPrompt.SetPathAnimation(utils.Move_Decelerate, utils.Vec2{X: 150, Y: 150}, 360, 10); err != nil {
		return nil, fmt.Errorf("ベンチマーク案内のアニメーション設定に失敗しました。\n%w", err)
	}

	installedParts := resultui.NewInstalledParts(
		utils.Position{X: 30, Y: 30},
		image.Point{X: 100, Y: 20},
		gdata,
	)

	partsList, err := newPartsList(gdata)
	if err != nil {
		return nil, err
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{guide, benchmarkPrompt, installedParts, partsList} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("結果画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	return &ResultCanvas{
		scene:           s,
		gdata:           gdata,
		canvas:          c,
		benchmarkPrompt: benchmarkPrompt,
		installedParts:  installedParts,
		partsList:       partsList,
	}, nil
}

// Update implements [scene.ResultCanvasAPI].
// SceneとGamedataの状態を観測して一覧と組み立て表示を切り替え、状態に応じたUIだけを更新する。
func (c *ResultCanvas) Update(active bool) error {
	switch c.gdata.State {
	case game.Result_SelectParts:
		if revision := c.scene.CandidatesRevision(); revision != c.shownRevision {
			parts, err := c.scene.CandidateParts()
			if err != nil {
				return fmt.Errorf("結果画面のパーツ候補の取得に失敗しました。\n%w", err)
			}
			if err := c.partsList.ShowParts(parts); err != nil {
				return fmt.Errorf("結果画面のパーツ一覧の構築に失敗しました。\n%w", err)
			}
			c.partsList.Show()
			c.shownRevision = revision
		}

	case game.Result_Assembled:
		if !c.installedStarted {
			c.partsList.Hide()
			if err := c.installedParts.Start(); err != nil {
				return fmt.Errorf("InstalledPartsUIの開始に失敗しました。\n%w", err)
			}
			c.installedStarted = true
		}
		if err := c.installedParts.Update(true); err != nil {
			return fmt.Errorf("InstalledPartsUIの更新に失敗しました。\n%w", err)
		}

	case game.Result_Benchmark:
		if err := c.benchmarkPrompt.Update(true); err != nil {
			return fmt.Errorf("ベンチマーク案内UIの更新に失敗しました。\n%w", err)
		}
	}
	return nil
}

// Draw implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsResultCanvas implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) IsResultCanvas() {}

// MovePartSelectionDown implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) MovePartSelectionDown() {
	if c.partsList.Visible() {
		c.partsList.SelectNext()
	}
}

// MovePartSelectionUp implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) MovePartSelectionUp() {
	if c.partsList.Visible() {
		c.partsList.SelectPrevious()
	}
}

// SelectedPart implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) SelectedPart() *game.Part {
	return c.partsList.SelectedPart()
}

// IsInstalledPartsFinished implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) IsInstalledPartsFinished() bool {
	return c.installedParts.IsFinished()
}

var _ scene.ResultCanvasAPI = (*ResultCanvas)(nil)
