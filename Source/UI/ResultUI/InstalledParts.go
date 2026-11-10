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
	parts []*ui.AnimatableUI

	pos        *Utils.Position
	recordSize image.Point

	gdata *game.Gamedata

	timer *tick.Timer

	active []bool
	newest uint

	ready bool
}

const (
	RecordDistance  = 30
	TextSize        = 20
	UIsOffset       = 200 //ui移動開始前の位置差分; 正数値で左からスライドイン
	SlideinInterval = 4   //unit: 0.1s
)

func NewInstalledParts(pos Utils.Position, recordSize image.Point, gdata *game.Gamedata) *InstalledParts {
	return &InstalledParts{
		parts:      []*ui.AnimatableUI{},
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

		t, err := ui.NewText(startPosition, color.Orange(), v.Name, ui.FontFace(ips.gdata.Font, TextSize))
		if err != nil {
			return fmt.Errorf("SetPartsInfo()のテキスト生成NewText()でエラーが発生しました\n%w\n", err)
		}
		aui, err := ui.NewAnimatableUI(
			startPosition,
			ips.recordSize,
			nil,
			t,
		)
		if err != nil {
			return fmt.Errorf("SetPartsInfo()のAnimatableUI生成NewAnimatableUI()でエラーが発生しました\n%w\n", err)
		}
		if err := aui.SetPathAnimation(
			Utils.Move_Walk_timeBased,
			targetPosition.PtoV(),
			0,
			SlideinInterval,
		); err != nil {
			return fmt.Errorf("%d 番目のアニメーション設定に失敗しました\n%w\n", i+1, err)
		}
		ips.parts = append(ips.parts, aui)
	}
	for i := 0; i < len(ips.parts); i++ {
		ips.active = append(ips.active, false)
	}
	return nil
}

func (i *InstalledParts) Draw(screen *ebiten.Image) {
	for _, v := range i.parts {
		v.Draw(screen)
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
	if i.newest == uint(len(i.active)) {
		if err := i.gdata.NextState(); err != nil {
			return err
		}
		return nil
	}
	i.timer.Set(SlideinInterval)
	if i.timer.Expired() {
		i.active[i.newest] = true
		i.newest++
	}
	for idx, actv := range i.active {
		if err := i.parts[idx].Update(actv); err != nil {
			return fmt.Errorf("InstalledParts.Update()でエラーが発生しました\n%w\n", err)
		}
	}
	return nil
}
