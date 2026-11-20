package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	settingsui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/SettingsUI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// 以下は設定画面の見た目の調整値。
const (
	settingsTitleTextSize = 24
	settingsTitleMarginY  = 30
	settingsTextSize      = 20
	settingsRowHeight     = 36
	settingsRowSpacing    = 14
	settingsBackWidth     = 200
	settingsBackHeight    = 30
	settingsBackMarginTop = 24
	settingsHintTextSize  = 14
	settingsHintMarginTop = 18
	settingsHintMessage   = "↑↓：えらぶ　←→：おんりょう　Esc：もどる"
	settingsBGMLabel      = "BGM"
	settingsSELabel       = "SE"
	settingsBackLabel     = "もどる"
	settingsTitleMessage  = "せってい"
)

// SettingsCanvas は SettingsScene のUI構成である。見出し、BGM/SE の音量行、「もどる」ボタン、操作の案内を持つ。
// 選択中の項目と音量はSceneが持つため、Sceneのポインタを持って観測する。
type SettingsCanvas struct {
	scene  *scene.SettingsScene
	canvas *canvas.Canvas
	title  *ui.Text
	hint   *ui.Text
	bgm    *settingsui.VolumeSlider
	se     *settingsui.VolumeSlider
	back   *ui.ButtonUI

	shownRevision int
}

// NewSettingsCanvas は SettingsCanvas を生成する。
func NewSettingsCanvas(s *scene.SettingsScene, gdata *game.Gamedata) (*SettingsCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("SettingsScene が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}
	if gdata.Font == nil {
		return nil, fmt.Errorf("gdata.Font が nil です")
	}

	title, err := ui.NewText(utils.Position{}, color.Orange(), settingsTitleMessage, ui.FontFace(gdata.Font, settingsTitleTextSize))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() -> title でエラーが発生しました。\n%w", err)
	}
	hint, err := ui.NewText(utils.Position{}, color.LightGray(), settingsHintMessage, ui.FontFace(gdata.Font, settingsHintTextSize))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() -> hint でエラーが発生しました。\n%w", err)
	}

	newSlider := func(label string, value int) (*settingsui.VolumeSlider, error) {
		return settingsui.NewVolumeSlider(settingsui.VolumeSliderConfig{
			Label:           label,
			Max:             game.MaxVolume,
			Value:           value,
			Font:            ui.FontFace(gdata.Font, settingsTextSize),
			TextColor:       color.White(),
			FilledColor:     color.Orange(),
			EmptyColor:      color.Gray(),
			BackgroundColor: color.DarkGray(),
			Height:          settingsRowHeight,
		})
	}
	bgm, err := newSlider(settingsBGMLabel, s.BGMVolume())
	if err != nil {
		return nil, fmt.Errorf("BGM音量UIの生成に失敗しました。\n%w", err)
	}
	se, err := newSlider(settingsSELabel, s.SEVolume())
	if err != nil {
		return nil, fmt.Errorf("SE音量UIの生成に失敗しました。\n%w", err)
	}

	backText, err := ui.NewText(utils.Position{}, color.White(), settingsBackLabel, ui.FontFace(gdata.Font, settingsTextSize))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() -> backText でエラーが発生しました。\n%w", err)
	}
	backImage := ebiten.NewImage(settingsBackWidth, settingsBackHeight)
	backImage.Fill(color.DarkGray().ToRGBA())
	backContent, err := ui.NewRichUI(utils.Position{}, image.Point{X: settingsBackWidth, Y: settingsBackHeight}, backImage, backText)
	if err != nil {
		return nil, fmt.Errorf("ui.NewRichUI() -> back でエラーが発生しました。\n%w", err)
	}
	back, err := ui.NewButtonUI(backContent)
	if err != nil {
		return nil, fmt.Errorf("ui.NewButtonUI() -> back でエラーが発生しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{title, bgm, se, back, hint} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("設定画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	sc := &SettingsCanvas{scene: s, canvas: c, title: title, hint: hint, bgm: bgm, se: se, back: back}
	sc.sync()
	return sc, nil
}

// sync は Scene の選択中の項目と音量を各UIへ反映する。
func (c *SettingsCanvas) sync() {
	selected := c.scene.SelectedMenuItem()
	c.bgm.SetSelected(selected == scene.SettingsBGM)
	c.se.SetSelected(selected == scene.SettingsSE)
	c.back.SetSelected(selected == scene.SettingsBack)
	c.bgm.SetValue(c.scene.BGMVolume())
	c.se.SetValue(c.scene.SEVolume())
	c.shownRevision = c.scene.SelectionRevision()
}

// Update implements [scene.SettingsCanvasAPI].
func (c *SettingsCanvas) Update(active bool) error {
	// 音量は選択が動かなくても変わるため、毎フレーム反映する（反映は軽い）
	c.sync()
	return c.canvas.Update(active)
}

// Draw implements [scene.SettingsCanvasAPI].
// 画面サイズに合わせて各UIを中央へ縦に配置する。
func (c *SettingsCanvas) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	sliderSize := c.bgm.Size()

	totalH := sliderSize.Y*2 + settingsRowSpacing + settingsBackMarginTop + settingsBackHeight
	top := float64(bounds.Dy()-totalH) / 2
	centerX := float64(bounds.Dx()) / 2

	c.title.SetPosition(utils.Position{
		X: centerX - settingsTitleTextSize*2,
		Y: top - settingsTitleMarginY - settingsTitleTextSize,
	})
	c.bgm.SetPosition(utils.Position{X: centerX - float64(sliderSize.X)/2, Y: top})
	c.se.SetPosition(utils.Position{X: centerX - float64(sliderSize.X)/2, Y: top + float64(sliderSize.Y+settingsRowSpacing)})
	backY := top + float64(sliderSize.Y*2+settingsRowSpacing+settingsBackMarginTop)
	c.back.RichUI().SetPosition(utils.Position{X: centerX - settingsBackWidth/2, Y: backY})
	c.hint.SetPosition(utils.Position{
		X: centerX - float64(len([]rune(settingsHintMessage)))*settingsHintTextSize/2,
		Y: backY + settingsBackHeight + settingsHintMarginTop,
	})
	c.canvas.Draw(screen)
}

// IsSettingsCanvas implements [scene.SettingsCanvasAPI].
func (c *SettingsCanvas) IsSettingsCanvas() {}

var _ scene.SettingsCanvasAPI = (*SettingsCanvas)(nil)
