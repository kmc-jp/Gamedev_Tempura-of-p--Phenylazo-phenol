package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	commonui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/CommonUI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// PlayCanvas は PlayScene のUI構成である。
// 残り時間、空き重量・所持金、パーツ一覧、発見パーツのログ、メッセージウィンドウを持つ。
// 一覧を開く時点はSceneが決めるため、Sceneのポインタを持って観測する。
type PlayCanvas struct {
	scene            *scene.PlayScene
	gdata            *game.Gamedata
	canvas           *canvas.Canvas
	partsList        *ui.PartsListUI
	foundPartLog     *commonui.LogWindow
	messageWindow    *commonui.MessageWindow
	purchaseConfirm  *commonui.Confirm
	exitStoreConfirm *commonui.Confirm

	// shownSelectListRevision は反映済みの PlayScene.SelectListRevision。
	shownSelectListRevision int
	// shownPartsRevision は一覧に反映済みの Gamedata.PartsRevision。
	shownPartsRevision int
}

const (
	foundPartLogDuration      = 10 // unit: 0.1s
	foundPartLogOffsetX       = 150
	foundPartLogOffsetY       = 50
	foundPartLogTextSize      = 10
	messageWindowTextSize     = 20
	messageWindowPaddingX     = 10
	messageWindowPaddingY     = 5
	messageWindowBottomMargin = 20
	confirmWidth              = 280
	confirmHeight             = 110
	confirmPadding            = 12
	confirmButtonWidth        = 90
	confirmButtonHeight       = 30
)

var (
	foundPartLogTextColor        = color.Black()
	messageWindowTextColor       = color.White()
	messageWindowBackgroundColor = color.SemiTransparentBlack()
)

// NewPlayCanvas は PlayCanvas を生成する。
func NewPlayCanvas(s *scene.PlayScene, gdata *game.Gamedata) (*PlayCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("PlayScene が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}

	timeText, err := ui.NewGameDataText(
		utils.Position{X: 0, Y: 20},
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

	worldTimeText, err := commonui.NewClock(
		utils.Position{X: 0, Y: 0},
		color.White(),
		gdata,
		ui.FontFace(gdata.Font, 25),
	)
	if err != nil {
		return nil, fmt.Errorf("commonui.NewClock() でエラーが発生しました。\n%w", err)
	}

	if gdata.AssetsManager == nil {
		return nil, fmt.Errorf("gdata.AssetsManager が nil です")
	}

	weightIconData, err := gdata.AssetsManager.Image.NewImageData(
		"みんげー_リュックサック_アイコン",
		"UI/みんげー_リュックサック_アイコン.png",
		"UI/みんげー_リュックサック_アイコン.json",
	)
	if err != nil {
		return nil, fmt.Errorf("リュックサックアイコンのImageData作成に失敗しました。\n%w", err)
	}
	_, weightIcon, err := gdata.AssetsManager.Image.Load(*weightIconData)
	if err != nil {
		return nil, fmt.Errorf("リュックサックアイコンの読み込みに失敗しました。\n%w", err)
	}

	moneyIconData, err := gdata.AssetsManager.Image.NewImageData(
		"みんげー_がま口アイコン",
		"UI/みんげー_がま口アイコン.png",
		"UI/みんげー_がま口アイコン.json",
	)
	if err != nil {
		return nil, fmt.Errorf("がま口アイコンのImageData作成に失敗しました。\n%w", err)
	}
	_, moneyIcon, err := gdata.AssetsManager.Image.Load(*moneyIconData)
	if err != nil {
		return nil, fmt.Errorf("がま口アイコンの読み込みに失敗しました。\n%w", err)
	}

	gameDataUI, err := ui.NewGameDataUI(
		gdata,
		[]ui.GameDataUIItem{
			{
				Image: weightIcon,
				Text: func(g game.Gamedata) string {
					return fmt.Sprintf("空き重量: %d", g.Weight)
				},
				Font:  ui.FontFace(gdata.Font, 15),
				Color: color.Black(),
			},
			{
				Image: moneyIcon,
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
	gameDataUI.SetBackgroundColor(color.White())

	partsList, err := newPartsList(gdata)
	if err != nil {
		return nil, err
	}

	foundPartLog, err := commonui.NewLogWindow(
		utils.Position{},
		gdata.Font,
		foundPartLogTextSize,
		foundPartLogTextColor,
		foundPartLogDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("commonui.NewLogWindow() でエラーが発生しました。\n%w", err)
	}

	messageFont := ui.FontFace(gdata.Font, messageWindowTextSize)
	messageWidth, messageHeight := text.Measure("探索終了", messageFont, messageFont.Size*1.5)
	messageWindow, err := commonui.NewMessageWindow(commonui.MessageWindowConfig{
		Size: image.Point{
			X: int(math.Ceil(messageWidth)) + messageWindowPaddingX*2,
			Y: int(math.Ceil(messageHeight)) + messageWindowPaddingY*2,
		},
		Font:            messageFont,
		TextColor:       messageWindowTextColor,
		BackgroundColor: messageWindowBackgroundColor,
		Padding:         image.Point{X: messageWindowPaddingX, Y: messageWindowPaddingY},
	})
	if err != nil {
		return nil, fmt.Errorf("commonui.NewMessageWindow() でエラーが発生しました。\n%w", err)
	}
	messageWindow.SetText("探索中")
	messageWindow.ShowText()

	purchaseConfirm, err := commonui.NewConfirm(commonui.ConfirmConfig{
		Size:                  image.Point{X: confirmWidth, Y: confirmHeight},
		Font:                  ui.FontFace(gdata.Font, messageWindowTextSize),
		TextColor:             messageWindowTextColor,
		BackgroundColor:       messageWindowBackgroundColor,
		ButtonBackgroundColor: color.DarkGray(),
		Padding:               image.Point{X: confirmPadding, Y: confirmPadding},
		ButtonSize:            image.Point{X: confirmButtonWidth, Y: confirmButtonHeight},
		Message:               "購入しますか？",
		YesLabel:              "はい",
		NoLabel:               "いいえ",
	})
	if err != nil {
		return nil, fmt.Errorf("commonui.NewConfirm() -> purchaseConfirm でエラーが発生しました。\n%w", err)
	}

	exitStoreConfirm, err := commonui.NewConfirm(commonui.ConfirmConfig{
		Size:                  image.Point{X: confirmWidth, Y: confirmHeight},
		Font:                  ui.FontFace(gdata.Font, messageWindowTextSize),
		TextColor:             messageWindowTextColor,
		BackgroundColor:       messageWindowBackgroundColor,
		ButtonBackgroundColor: color.DarkGray(),
		Padding:               image.Point{X: confirmPadding, Y: confirmPadding},
		ButtonSize:            image.Point{X: confirmButtonWidth, Y: confirmButtonHeight},
		Message:               "店を出ますか？",
		YesLabel:              "はい",
		NoLabel:               "いいえ",
	})
	if err != nil {
		return nil, fmt.Errorf("commonui.NewConfirm() -> exitStoreConfirm でエラーが発生しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{worldTimeText, timeText, partsList, gameDataUI, foundPartLog, messageWindow, purchaseConfirm, exitStoreConfirm} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("プレイ画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	return &PlayCanvas{
		scene:                   s,
		gdata:                   gdata,
		canvas:                  c,
		partsList:               partsList,
		foundPartLog:            foundPartLog,
		messageWindow:           messageWindow,
		purchaseConfirm:         purchaseConfirm,
		exitStoreConfirm:        exitStoreConfirm,
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

// ShowFoundPartLog は発見したパーツ名をプレイヤー頭上に表示する。
func (c *PlayCanvas) ShowFoundPartLog(text string) {
	c.foundPartLog.SetText(text)
}

// Draw implements [scene.PlayCanvasAPI].
func (c *PlayCanvas) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	c.foundPartLog.SetPosition(utils.Position{
		X: float64(bounds.Dx()/2 - foundPartLogOffsetX),
		Y: float64(bounds.Dy()/2 - foundPartLogOffsetY),
	})
	messageSize := c.messageWindow.Size()
	c.messageWindow.SetPosition(utils.Position{
		X: float64((bounds.Dx() - messageSize.X) / 2),
		Y: float64(bounds.Dy() - messageSize.Y - messageWindowBottomMargin),
	})
	confirmSize := c.purchaseConfirm.Size()
	c.purchaseConfirm.SetPosition(utils.Position{
		X: float64((bounds.Dx() - confirmSize.X) / 2),
		Y: float64((bounds.Dy() - confirmSize.Y) / 2),
	})
	exitStoreConfirmSize := c.exitStoreConfirm.Size()
	c.exitStoreConfirm.SetPosition(utils.Position{
		X: float64((bounds.Dx() - exitStoreConfirmSize.X) / 2),
		Y: float64((bounds.Dy() - exitStoreConfirmSize.Y) / 2),
	})
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

func (c *PlayCanvas) ShowPurchaseConfirmation() bool {
	if !c.partsList.IsShowingFoundParts() || c.partsList.SelectedPart() == nil {
		return false
	}
	c.purchaseConfirm.Show()
	return true
}

func (c *PlayCanvas) MovePurchaseConfirmationLeft() {
	c.purchaseConfirm.MoveLeft()
}

func (c *PlayCanvas) MovePurchaseConfirmationRight() {
	c.purchaseConfirm.MoveRight()
}

func (c *PlayCanvas) ConfirmPurchase() {
	c.purchaseConfirm.Confirm()
}

func (c *PlayCanvas) IsPurchaseConfirmationOpen() bool {
	return c.purchaseConfirm.Visible()
}

func (c *PlayCanvas) PurchaseConfirmationResult() (bool, bool) {
	return c.purchaseConfirm.Result()
}

func (c *PlayCanvas) HidePurchaseConfirmation() {
	c.purchaseConfirm.Hide()
}

func (c *PlayCanvas) ShowExitStoreConfirmation() bool {
	c.exitStoreConfirm.Show()
	return true
}

func (c *PlayCanvas) MoveExitStoreConfirmationLeft() {
	c.exitStoreConfirm.MoveLeft()
}

func (c *PlayCanvas) MoveExitStoreConfirmationRight() {
	c.exitStoreConfirm.MoveRight()
}

func (c *PlayCanvas) ConfirmExitStore() {
	c.exitStoreConfirm.Confirm()
}

func (c *PlayCanvas) IsExitStoreConfirmationOpen() bool {
	return c.exitStoreConfirm.Visible()
}

func (c *PlayCanvas) ExitStoreConfirmationResult() (bool, bool) {
	return c.exitStoreConfirm.Result()
}

func (c *PlayCanvas) HideExitStoreConfirmation() {
	c.exitStoreConfirm.Hide()
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
	case ui.PartsListExit:
		return scene.PartsListTabExit
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

func (c *PlayCanvas) SetMessageWindowText(s string) {
	c.messageWindow.SetText(s)
}

var _ scene.PlayCanvasAPI = (*PlayCanvas)(nil)
