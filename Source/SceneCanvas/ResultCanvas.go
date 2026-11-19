package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	commonui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/CommonUI"
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

	// canvas は全UIを毎フレーム更新・描画する。
	// 各UIは開始前・非表示中の Update では何もしないため、状態ごとの更新の振り分けは不要。
	canvas          *canvas.Canvas
	benchmarkPrompt *ui.AnimatableUI
	installedParts  *resultui.InstalledParts
	partsList       *ui.PartsListUI
	partUseConfirm  *commonui.Confirm
	scoreDisplay    *resultui.ScoreDisplay

	// shownRevision は一覧に反映済みの ResultScene.CandidatesRevision。
	shownRevision int
	// installedStarted は組み立て済みパーツの表示を開始したか。
	installedStarted bool
	// benchmarkStarted はベンチマーク案内のアニメーションを開始したか。
	benchmarkStarted bool
}

// ベンチマーク案内のアニメーション設定
var (
	benchmarkPromptMotion   = utils.Move_Decelerate
	benchmarkPromptTo       = utils.Vec2{X: 150, Y: 150}
	benchmarkPromptSpeed    = 360.0
	benchmarkPromptDuration = uint(10)
)

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
	// アニメーションは Result_Benchmark に変わったときに Update で開始する

	installedParts := resultui.NewInstalledParts(
		utils.Position{X: 30, Y: 30},
		image.Point{X: 100, Y: 20},
		gdata,
	)

	partsList, err := newPartsList(gdata)
	if err != nil {
		return nil, err
	}

	partUseConfirm, err := commonui.NewConfirm(commonui.ConfirmConfig{
		Size:                  image.Point{X: 280, Y: 110},
		Font:                  ui.FontFace(gdata.Font, 20),
		TextColor:             color.White(),
		BackgroundColor:       color.SemiTransparentBlack(),
		ButtonBackgroundColor: color.DarkGray(),
		Padding:               image.Point{X: 12, Y: 12},
		ButtonSize:            image.Point{X: 90, Y: 30},
		Message:               "このパーツをつかいますか？",
		YesLabel:              "はい",
		NoLabel:               "いいえ",
	})
	if err != nil {
		return nil, fmt.Errorf("commonui.NewConfirm() でエラーが発生しました。\n%w", err)
	}

	scoreDisplay, err := resultui.NewScoreDisplay(resultui.ScoreDisplayConfig{
		Font:            ui.FontFace(gdata.Font, 60),
		TextColor:       color.Black(),
		BackgroundColor: color.LightGray(),
		Padding:         12,
		ShuffleFrames:   2,
		DigitFrames:     6,
	})
	if err != nil {
		return nil, fmt.Errorf("スコア表示UIの生成に失敗しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{guide, benchmarkPrompt, installedParts, partsList, partUseConfirm, scoreDisplay} {
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
		partUseConfirm:  partUseConfirm,
		scoreDisplay:    scoreDisplay,
	}, nil
}

// Update implements [scene.ResultCanvasAPI].
// SceneとGamedataの状態を観測して一覧・組み立て表示・ベンチマーク案内を切り替え、全UIを更新する。
func (c *ResultCanvas) Update(active bool) error {
	switch c.gdata.State {
	case game.Result_SelectParts:
		if revision := c.scene.CandidatesRevision(); revision != c.shownRevision {
			if err := c.partsList.ShowParts(c.scene.CandidateParts()); err != nil {
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

	case game.Result_Benchmark:
		if !c.benchmarkStarted {
			if err := c.benchmarkPrompt.SetPathAnimation(
				benchmarkPromptMotion,
				benchmarkPromptTo,
				benchmarkPromptSpeed,
				benchmarkPromptDuration,
			); err != nil {
				return fmt.Errorf("ベンチマーク案内のアニメーション設定に失敗しました。\n%w", err)
			}
			c.benchmarkStarted = true
		}
	}

	// installedParts は Start 前、benchmarkPrompt はアニメーション開始前、
	// partsList は非表示中の Update で何もしないため、毎フレーム全UIを更新してよい。
	return c.canvas.Update(active)
}

// Draw implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	c.scoreDisplay.SetCenter(utils.Position{
		X: float64(bounds.Dx() / 2),
		Y: float64(bounds.Dy() / 2),
	})
	if c.partUseConfirm.Visible() {
		confirmSize := c.partUseConfirm.Size()
		c.partUseConfirm.SetPosition(utils.Position{
			X: float64((bounds.Dx() - confirmSize.X) / 2),
			Y: float64((bounds.Dy() - confirmSize.Y) / 2),
		})
	}
	c.canvas.Draw(screen)
}

// IsResultCanvas implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) IsResultCanvas() {}

// MovePartSelectionDown implements [scene.ResultCanvasAPI].
// 一覧の非表示中は PartsListUI が何もしない。
func (c *ResultCanvas) MovePartSelectionDown() {
	c.partsList.SelectNext()
}

// MovePartSelectionUp implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) MovePartSelectionUp() {
	c.partsList.SelectPrevious()
}

func (c *ResultCanvas) ShowPartUseConfirmation() bool {
	if c.partsList.SelectedPart() == nil {
		return false
	}
	c.partUseConfirm.Show()
	return true
}

func (c *ResultCanvas) MovePartUseConfirmationLeft() {
	c.partUseConfirm.MoveLeft()
}

func (c *ResultCanvas) MovePartUseConfirmationRight() {
	c.partUseConfirm.MoveRight()
}

func (c *ResultCanvas) ConfirmPartUse() {
	c.partUseConfirm.Confirm()
}

func (c *ResultCanvas) IsPartUseConfirmationOpen() bool {
	return c.partUseConfirm.Visible()
}

func (c *ResultCanvas) PartUseConfirmationResult() (bool, bool) {
	return c.partUseConfirm.Result()
}

func (c *ResultCanvas) HidePartUseConfirmation() {
	c.partUseConfirm.Hide()
}

func (c *ResultCanvas) StartScoreDisplay(score int) {
	c.scoreDisplay.Start(score)
}

func (c *ResultCanvas) IsScoreDisplayFinished() bool {
	return c.scoreDisplay.IsFinished()
}

// SelectedPart implements [scene.ResultCanvasAPI].
// 一覧の非表示中は PartsListUI が nil を返す。
func (c *ResultCanvas) SelectedPart() *game.Part {
	return c.partsList.SelectedPart()
}

// IsInstalledPartsFinished implements [scene.ResultCanvasAPI].
func (c *ResultCanvas) IsInstalledPartsFinished() bool {
	return c.installedParts.IsFinished()
}

var _ scene.ResultCanvasAPI = (*ResultCanvas)(nil)
