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

// MessageWindowState はMessageWindowの表示状態である。
type MessageWindowState uint8

const (
	MessageWindowTextHidden MessageWindowState = iota
	MessageWindowTextVisible
	MessageWindowHidden
)

// MessageWindowConfig はMessageWindowの生成時設定である。
type MessageWindowConfig struct {
	Position        utils.Position
	Size            image.Point
	Font            *text.GoTextFace
	TextColor       *gamecolor.Color
	BackgroundColor *gamecolor.Color
	BackgroundImage *ebiten.Image
	Padding         image.Point
}

// MessageWindow は背景とテキストを持つ常時表示用のウィンドウである。
// テキストの折り返しは行わず、文字列と表示領域を分離して保持する。
type MessageWindow struct {
	position        utils.Position
	size            image.Point
	font            *text.GoTextFace
	textColor       *gamecolor.Color
	backgroundColor *gamecolor.Color
	backgroundImage *ebiten.Image
	padding         image.Point
	text            *ui.Text
	state           MessageWindowState
}

// NewMessageWindow は設定に基づいてMessageWindowを生成する。
func NewMessageWindow(config MessageWindowConfig) (*MessageWindow, error) {
	if config.Size.X <= 0 || config.Size.Y <= 0 {
		return nil, fmt.Errorf("size は 0 より大きい必要があります: %v", config.Size)
	}
	if config.Font == nil {
		return nil, fmt.Errorf("font が nil です")
	}
	if config.TextColor == nil {
		return nil, fmt.Errorf("textColor が nil です")
	}
	if config.Padding.X < 0 || config.Padding.Y < 0 {
		return nil, fmt.Errorf("padding は 0 以上である必要があります: %v", config.Padding)
	}
	if config.Padding.X > config.Size.X || config.Padding.Y > config.Size.Y {
		return nil, fmt.Errorf("padding が size を超えています: padding=%v, size=%v", config.Padding, config.Size)
	}
	if config.BackgroundImage != nil && config.BackgroundImage.Bounds().Empty() {
		return nil, fmt.Errorf("backgroundImage のサイズが空です")
	}

	textPosition := utils.Position{
		X: config.Position.X + float64(config.Padding.X),
		Y: config.Position.Y + float64(config.Padding.Y),
	}
	text, err := ui.NewText(textPosition, config.TextColor, "", config.Font)
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました: %w", err)
	}

	return &MessageWindow{
		position:        config.Position,
		size:            config.Size,
		font:            config.Font,
		textColor:       config.TextColor,
		backgroundColor: config.BackgroundColor,
		backgroundImage: config.BackgroundImage,
		padding:         config.Padding,
		text:            text,
		state:           MessageWindowTextHidden,
	}, nil
}

// Update implements [canvas.CanvasItem].
func (w *MessageWindow) Update(active bool) error {
	return w.text.Update(active)
}

// Draw implements [canvas.CanvasItem].
func (w *MessageWindow) Draw(screen *ebiten.Image) {
	if w.state == MessageWindowHidden {
		return
	}

	bounds := w.Bounds()
	if w.backgroundColor != nil {
		uiutils.FillRect(screen, bounds, w.backgroundColor.ToRGBA())
	}
	if w.backgroundImage != nil {
		imageBounds := w.backgroundImage.Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(
			float64(bounds.Dx())/float64(imageBounds.Dx()),
			float64(bounds.Dy())/float64(imageBounds.Dy()),
		)
		op.GeoM.Translate(w.position.X, w.position.Y)
		screen.DrawImage(w.backgroundImage, op)
	}
	if w.state == MessageWindowTextVisible {
		w.text.Draw(screen)
	}
}

// SetText は表示するテキストを差し替える。表示状態は変更しない。
func (w *MessageWindow) SetText(value string) {
	w.text.SetText(value)
}

// ClearText は表示するテキストを空にする。表示状態は変更しない。
func (w *MessageWindow) ClearText() {
	w.SetText("")
}

// ShowText はウィンドウとテキストを表示する。
func (w *MessageWindow) ShowText() {
	w.state = MessageWindowTextVisible
}

// HideText はウィンドウを表示したままテキストを非表示にする。
func (w *MessageWindow) HideText() {
	w.state = MessageWindowTextHidden
}

// Show はウィンドウを表示し、テキストを非表示にする。
func (w *MessageWindow) Show() {
	w.HideText()
}

// Hide はウィンドウとテキストを非表示にする。
func (w *MessageWindow) Hide() {
	w.state = MessageWindowHidden
}

// State は現在の表示状態を返す。
func (w *MessageWindow) State() MessageWindowState {
	return w.state
}

// Position はウィンドウの左上位置を返す。
func (w *MessageWindow) Position() utils.Position {
	return w.position
}

// Size はウィンドウサイズを返す。
func (w *MessageWindow) Size() image.Point {
	return w.size
}

// Bounds はウィンドウの描画範囲を返す。
func (w *MessageWindow) Bounds() image.Rectangle {
	x := int(w.position.X)
	y := int(w.position.Y)
	return image.Rect(x, y, x+w.size.X, y+w.size.Y)
}

// SetPosition はウィンドウとテキストの位置を移動する。
func (w *MessageWindow) SetPosition(position utils.Position) {
	w.position = position
	w.text.SetPosition(utils.Position{
		X: position.X + float64(w.padding.X),
		Y: position.Y + float64(w.padding.Y),
	})
}

var _ canvas.CanvasItem = (*MessageWindow)(nil)
var _ canvas.PositionItem = (*MessageWindow)(nil)
