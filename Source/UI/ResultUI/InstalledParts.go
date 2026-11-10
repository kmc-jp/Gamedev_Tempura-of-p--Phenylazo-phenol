package resultui

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type InstalledParts struct {
	parts []installedPart

	pos        *Utils.Position
	recordSize image.Point

	gdata *game.Gamedata

	timer *tick.Timer

	active []bool
	newest uint

	ready bool
}

type installedPart struct {
	category *ui.AnimatableUI
	name     *ui.AnimatableUI
}

func (p installedPart) Update(active bool) error {
	if p.category != nil {
		if err := p.category.Update(active); err != nil {
			return err
		}
	}
	return p.name.Update(active)
}

func (p installedPart) IsAnimating() bool {
	if p.category != nil && p.category.IsAnimating() {
		return true
	}
	return p.name.IsAnimating()
}

const (
	RecordDistance  = 50
	TextSize        = 20
	CategorySize    = 14
	NameIndent      = 20
	UIsOffset       = 200 //ui移動開始前の位置差分; 正数値で左からスライドイン
	SlideinInterval = 4   //unit: 0.1s
)

func NewInstalledParts(pos Utils.Position, recordSize image.Point, gdata *game.Gamedata) *InstalledParts {
	return &InstalledParts{
		parts:      []installedPart{},
		pos:        &pos,
		recordSize: recordSize,
		gdata:      gdata,
		timer:      tick.NewTimer(),
		active:     []bool{},
		newest:     0,
		ready:      false,
	}
}

func (ips *InstalledParts) SetPartsInfo() error {
	for i, v := range ips.gdata.Assemble {
		targetPosition := Utils.Position{
			X: ips.pos.X,
			Y: ips.pos.Y + float64(i*RecordDistance),
		}
		startPosition := Utils.Position{
			X: targetPosition.X - UIsOffset,
			Y: targetPosition.Y,
		}

		categoryName := partCategoryName(uint(v.Kind))
		var categoryUI *ui.AnimatableUI
		if categoryName != "" {
			categoryText, err := ui.NewText(
				startPosition,
				color.Gold(),
				categoryName,
				ui.FontFace(ips.gdata.Font, CategorySize),
			)
			if err != nil {
				return fmt.Errorf("SetPartsInfo()の種別テキスト生成に失敗しました\n%w\n", err)
			}
			categoryUI, err = ui.NewAnimatableUI(
				startPosition,
				ips.recordSize,
				nil,
				categoryText,
			)
			if err != nil {
				return fmt.Errorf("SetPartsInfo()の種別AnimatableUI生成に失敗しました\n%w\n", err)
			}
			if err := categoryUI.SetPathAnimation(
				Utils.Move_Walk_timeBased,
				targetPosition.PtoV(),
				0,
				SlideinInterval,
			); err != nil {
				return fmt.Errorf("%d 番目の種別アニメーション設定に失敗しました\n%w\n", i+1, err)
			}
		}

		namePosition := Utils.Position{
			X: startPosition.X + NameIndent,
			Y: startPosition.Y + float64(CategorySize),
		}
		nameTargetPosition := Utils.Position{
			X: targetPosition.X + NameIndent,
			Y: targetPosition.Y + float64(CategorySize),
		}
		t, err := ui.NewText(namePosition, color.Orange(), v.Name, ui.FontFace(ips.gdata.Font, TextSize))
		if err != nil {
			return fmt.Errorf("SetPartsInfo()のテキスト生成NewText()でエラーが発生しました\n%w\n", err)
		}
		aui, err := ui.NewAnimatableUI(
			namePosition,
			ips.recordSize,
			nil,
			t,
		)
		if err != nil {
			return fmt.Errorf("SetPartsInfo()のAnimatableUI生成NewAnimatableUI()でエラーが発生しました\n%w\n", err)
		}
		if err := aui.SetPathAnimation(
			Utils.Move_Walk_timeBased,
			nameTargetPosition.PtoV(),
			0,
			SlideinInterval,
		); err != nil {
			return fmt.Errorf("%d 番目のアニメーション設定に失敗しました\n%w\n", i+1, err)
		}
		ips.parts = append(ips.parts, installedPart{
			category: categoryUI,
			name:     aui,
		})
	}
	for i := 0; i < len(ips.parts); i++ {
		ips.active = append(ips.active, false)
	}
	return nil
}

func (i *InstalledParts) Draw(screen *ebiten.Image) {
	for _, part := range i.parts {
		if part.category != nil {
			part.category.Draw(screen)
		}
		part.name.Draw(screen)
	}
}

func (i *InstalledParts) Update(active bool) error {
	if i.ready == false {
		if err := i.SetPartsInfo(); err != nil {
			return fmt.Errorf("InstalledParts の初期化に失敗しました: %w", err)
		}
		i.ready = true
		return nil
	}
	if i.newest < uint(len(i.active)) {
		i.timer.Set(SlideinInterval)
		if i.timer.Expired() {
			i.active[i.newest] = true
			i.newest++
		}
	}
	for idx, actv := range i.active {
		if err := i.parts[idx].Update(actv); err != nil {
			return fmt.Errorf("InstalledParts.Update()でエラーが発生しました\n%w\n", err)
		}
	}

	if i.newest == uint(len(i.active)) && i.allAnimationsFinished() {
		if err := i.gdata.NextState(); err != nil {
			return err
		}
	}
	return nil
}

func (i *InstalledParts) allAnimationsFinished() bool {
	for _, part := range i.parts {
		if part.IsAnimating() {
			return false
		}
	}
	return true
}

func partCategoryName(kind uint) string {
	switch kind {
	case uint(game.CPU):
		return "CPU"
	case uint(game.GPU):
		return "グラフィックボード"
	case uint(game.RAM):
		return "メモリ"
	case uint(game.MB):
		return "マザーボード"
	case uint(game.Case):
		return "ケース"
	default:
		return ""
	}
}
