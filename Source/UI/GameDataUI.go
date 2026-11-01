package ui

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameDataUIItem struct {
	Image *ebiten.Image
	Text  func(game.Gamedata) string
	Font  *text.GoTextFace
	Color *color.Color
}

type GameDataUI struct {
	gdata  *game.Gamedata
	items  []GameDataUIItem
	margin float64
	gap    float64
}

func NewGameDataUI(
	gdata *game.Gamedata,
	items []GameDataUIItem,
	margin float64,
	gap float64,
) (*GameDataUI, error) {
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}
	if margin < 0 || gap < 0 {
		return nil, fmt.Errorf("margin と gap は 0 以上である必要があります")
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("items が空です")
	}
	for _, item := range items {
		if item.Text == nil {
			return nil, fmt.Errorf("item の Text が nil です")
		}
		if item.Font == nil {
			return nil, fmt.Errorf("item の Font が nil です")
		}
		if item.Color == nil {
			return nil, fmt.Errorf("item の Color が nil です")
		}
	}

	return &GameDataUI{
		gdata:  gdata,
		items:  items,
		margin: margin,
		gap:    gap,
	}, nil
}

func (u *GameDataUI) Update(active bool) error {
	return nil
}

func (u *GameDataUI) Draw(screen *ebiten.Image) {
	data := *u.gdata
	width := 0.0
	height := 0.0
	for _, item := range u.items {
		value := item.Text(data)
		textWidth, textHeight := text.Measure(value, item.Font, item.Font.Size*1.5)
		itemWidth := textWidth
		if item.Image != nil {
			itemWidth += float64(item.Image.Bounds().Dx()) + u.gap
			if float64(item.Image.Bounds().Dy()) > textHeight {
				textHeight = float64(item.Image.Bounds().Dy())
			}
		}
		width += itemWidth + u.gap
		if textHeight > height {
			height = textHeight
		}
	}
	width -= u.gap

	x := float64(screen.Bounds().Dx()) - u.margin - width
	y := u.margin
	for _, item := range u.items {
		value := item.Text(data)
		textWidth, textHeight := text.Measure(value, item.Font, item.Font.Size*1.5)
		if item.Image != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y+(height-float64(item.Image.Bounds().Dy()))/2)
			screen.DrawImage(item.Image, op)
			x += float64(item.Image.Bounds().Dx()) + u.gap
		}
		DrawText(utils.Position{X: x, Y: y + (height-textHeight)/2}, item.Color, value, item.Font, screen)
		x += textWidth + u.gap
	}
}

var _ scene.Element = (*GameDataUI)(nil)
