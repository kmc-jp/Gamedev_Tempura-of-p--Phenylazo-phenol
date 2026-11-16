package commonui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	gamecolor "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ConfirmChoice は確認UIで現在選択されている項目である。
type ConfirmChoice uint8

const (
	ConfirmNo ConfirmChoice = iota
	ConfirmYes
)

// ConfirmConfig は確認UIの生成時設定である。
type ConfirmConfig struct {
	Position              utils.Position
	Size                  image.Point
	Font                  *text.GoTextFace
	TextColor             *gamecolor.Color
	BackgroundColor       *gamecolor.Color
	ButtonBackgroundColor *gamecolor.Color
	Padding               image.Point
	ButtonSize            image.Point
	Message               string
	YesLabel              string
	NoLabel               string
}

// Confirm はメッセージと「はい」「いいえ」の選択肢を表示する受動的な確認UIである。
//
// 入力は扱わず、Scene・Componentが MoveLeft、MoveRight、Confirm を呼び出す。
// 確定結果は Result で問い合わせる。
type Confirm struct {
	position              utils.Position
	size                  image.Point
	backgroundColor       *gamecolor.Color
	buttonBackgroundColor *gamecolor.Color
	padding               image.Point
	message               *ui.Text
	yes                   *ui.ButtonUI
	no                    *ui.ButtonUI
	choice                ConfirmChoice
	visible               bool
	confirmed             bool
	result                bool
}

// NewConfirm は設定に基づいて確認UIを生成する。
func NewConfirm(config ConfirmConfig) (*Confirm, error) {
	if config.Size.X <= 0 || config.Size.Y <= 0 {
		return nil, fmt.Errorf("size は 0 より大きい必要があります: %v", config.Size)
	}
	if config.Font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	if config.TextColor == nil {
		return nil, fmt.Errorf("textColor が nil です")
	}
	if config.BackgroundColor == nil {
		return nil, fmt.Errorf("backgroundColor が nil です")
	}
	if config.ButtonBackgroundColor == nil {
		return nil, fmt.Errorf("buttonBackgroundColor が nil です")
	}
	if config.Padding.X < 0 || config.Padding.Y < 0 {
		return nil, fmt.Errorf("padding は 0 以上である必要があります: %v", config.Padding)
	}
	if config.ButtonSize.X <= 0 || config.ButtonSize.Y <= 0 {
		return nil, fmt.Errorf("buttonSize は 0 より大きい必要があります: %v", config.ButtonSize)
	}
	if config.Padding.X*2+config.ButtonSize.X*2 > config.Size.X ||
		config.Padding.Y*2+config.ButtonSize.Y > config.Size.Y {
		return nil, fmt.Errorf("padding または buttonSize が size を超えています: padding=%v, buttonSize=%v, size=%v",
			config.Padding, config.ButtonSize, config.Size)
	}
	if config.YesLabel == "" || config.NoLabel == "" {
		return nil, fmt.Errorf("yesLabel と noLabel は空にできません")
	}

	message, err := ui.NewText(
		utils.Position{
			X: config.Position.X + float64(config.Padding.X),
			Y: config.Position.Y + float64(config.Padding.Y),
		},
		config.TextColor,
		config.Message,
		config.Font,
	)
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました: %w", err)
	}

	buttonY := config.Position.Y + float64(config.Size.Y-config.Padding.Y-config.ButtonSize.Y)
	yesButton, err := newConfirmButton(
		utils.Position{X: config.Position.X + float64(config.Padding.X), Y: buttonY},
		config.ButtonSize, config.YesLabel, config.Font, config.TextColor,
	)
	if err != nil {
		return nil, err
	}
	noButton, err := newConfirmButton(
		utils.Position{
			X: config.Position.X + float64(config.Size.X-config.Padding.X-config.ButtonSize.X),
			Y: buttonY,
		},
		config.ButtonSize, config.NoLabel, config.Font, config.TextColor,
	)
	if err != nil {
		return nil, err
	}

	confirm := &Confirm{
		position:              config.Position,
		size:                  config.Size,
		backgroundColor:       config.BackgroundColor,
		buttonBackgroundColor: config.ButtonBackgroundColor,
		padding:               config.Padding,
		message:               message,
		yes:                   yesButton,
		no:                    noButton,
		choice:                ConfirmYes,
		visible:               false,
	}
	confirm.syncSelection()
	return confirm, nil
}

func newConfirmButton(
	position utils.Position,
	size image.Point,
	label string,
	font *text.GoTextFace,
	textColor *gamecolor.Color,
) (*ui.ButtonUI, error) {
	text, err := ui.NewText(position, textColor, label, font)
	if err != nil {
		return nil, fmt.Errorf("確認ボタンのテキスト生成に失敗しました: %w", err)
	}
	content, err := ui.NewRichUI(position, size, nil, text)
	if err != nil {
		return nil, fmt.Errorf("確認ボタンの生成に失敗しました: %w", err)
	}
	button, err := ui.NewButtonUI(content)
	if err != nil {
		return nil, fmt.Errorf("確認ボタンの選択UI生成に失敗しました: %w", err)
	}
	return button, nil
}

// Update は確認UIの子要素を更新する。入力処理は行わない。
func (c *Confirm) Update(active bool) error {
	if !c.visible {
		return nil
	}
	if err := c.message.Update(active); err != nil {
		return fmt.Errorf("確認メッセージの更新に失敗しました: %w", err)
	}
	if err := c.no.Update(active); err != nil {
		return fmt.Errorf("いいえボタンの更新に失敗しました: %w", err)
	}
	if err := c.yes.Update(active); err != nil {
		return fmt.Errorf("はいボタンの更新に失敗しました: %w", err)
	}
	return nil
}

// Draw は確認UIを描画する。
func (c *Confirm) Draw(screen *ebiten.Image) {
	if !c.visible {
		return
	}
	uiutils.FillRect(screen, c.Bounds(), c.backgroundColor.ToRGBA())
	uiutils.FillRect(screen, c.no.RichUI().Bounds(), c.buttonBackgroundColor.ToRGBA())
	uiutils.FillRect(screen, c.yes.RichUI().Bounds(), c.buttonBackgroundColor.ToRGBA())
	c.message.Draw(screen)
	c.no.Draw(screen)
	c.yes.Draw(screen)
}

// SetMessage は上部に表示するメッセージを差し替える。
func (c *Confirm) SetMessage(message string) {
	c.message.SetText(message)
}

// Show は確認UIを表示し、選択を「はい」に戻して未確定にする。
func (c *Confirm) Show() {
	c.visible = true
	c.choice = ConfirmYes
	c.confirmed = false
	c.result = false
	c.syncSelection()
}

// Hide は確認UIを非表示にする。確定結果は保持する。
func (c *Confirm) Hide() {
	c.visible = false
}

// Visible は確認UIが表示中かを返す。
func (c *Confirm) Visible() bool {
	return c.visible
}

// MoveLeft は選択を「はい」へ移す。
func (c *Confirm) MoveLeft() {
	if c.visible && !c.confirmed {
		c.choice = ConfirmYes
		c.syncSelection()
	}
}

// MoveRight は選択を「いいえ」へ移す。
func (c *Confirm) MoveRight() {
	if c.visible && !c.confirmed {
		c.choice = ConfirmNo
		c.syncSelection()
	}
}

// Selection は現在選択されている項目を返す。
func (c *Confirm) Selection() ConfirmChoice {
	return c.choice
}

// Confirm は現在の選択を確定する。非表示中または確定済みなら何もしない。
func (c *Confirm) Confirm() {
	if !c.visible || c.confirmed {
		return
	}
	c.confirmed = true
	c.result = c.choice == ConfirmYes
}

// Result は確定結果と、確定済みかどうかを返す。
func (c *Confirm) Result() (bool, bool) {
	return c.result, c.confirmed
}

func (c *Confirm) syncSelection() {
	c.no.SetSelected(c.choice == ConfirmNo)
	c.yes.SetSelected(c.choice == ConfirmYes)
}

// Position は確認UIの左上位置を返す。
func (c *Confirm) Position() utils.Position {
	return c.position
}

// SetPosition は確認UI全体を移動する。
func (c *Confirm) SetPosition(position utils.Position) {
	deltaX := position.X - c.position.X
	deltaY := position.Y - c.position.Y
	c.position = position
	c.message.SetPosition(offsetPosition(c.message.Position(), deltaX, deltaY))
	c.no.RichUI().SetPosition(offsetPosition(c.no.RichUI().Position(), deltaX, deltaY))
	c.yes.RichUI().SetPosition(offsetPosition(c.yes.RichUI().Position(), deltaX, deltaY))
}

// Size は確認UIのサイズを返す。
func (c *Confirm) Size() image.Point {
	return c.size
}

// Bounds は確認UIの描画範囲を返す。
func (c *Confirm) Bounds() image.Rectangle {
	x, y := int(c.position.X), int(c.position.Y)
	return image.Rect(x, y, x+c.size.X, y+c.size.Y)
}

func offsetPosition(position utils.Position, x, y float64) utils.Position {
	return utils.Position{X: position.X + x, Y: position.Y + y}
}

var _ canvas.CanvasItem = (*Confirm)(nil)
var _ canvas.PositionItem = (*Confirm)(nil)
