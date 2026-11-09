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
	RecordDistance  = 5
	TextSize        = 20
	UIsOffset       = 300 //ui移動開始前の位置差分; 正数値で左からスライドイン
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
	}
}

func (ips *InstalledParts) SetPartsInfo() error {
	for i, v := range ips.gdata.Assemble {
		//**未実装**初回のみ通常追加、以降Position.Yにインデックスi-1のPosition.Y + RecordDistanceを加算
		if i == 0 {
			t, err := ui.NewText(*ips.pos, color.Orange(), v.Name, ui.FontFace(ips.gdata.Font, TextSize))
			if err != nil {
				return fmt.Errorf("SetPartsInfo()のテキスト生成NewText()でエラーが発生しました\n%w\n", err)
			}
			aui, err := ui.NewAnimatableUI(
				Utils.Position{X: ips.pos.X - UIsOffset, Y: ips.pos.Y},
				ips.recordSize,
				nil,
				t,
			)
			if err != nil {
				return fmt.Errorf("SetPartsInfo()のAnimatableUI生成NewAnimatableUI()でエラーが発生しました\n%w\n", err)
			}
			ips.parts = append(ips.parts, aui)
		} else {
			t, err := ui.NewText(*ips.pos, color.Orange(), v.Name, ui.FontFace(ips.gdata.Font, TextSize))
			if err != nil {
				return fmt.Errorf("SetPartsInfo()のテキスト生成NewText()でエラーが発生しました\n%w\n", err)
			}
			aui, err := ui.NewAnimatableUI(
				Utils.Position{X: ips.pos.X - UIsOffset, Y: ips.parts[i-1].Position().Y + RecordDistance},
				ips.recordSize,
				nil,
				t,
			)
			if err != nil {
				return fmt.Errorf("SetPartsInfo()のAnimatableUI生成NewAnimatableUI()でエラーが発生しました\n%w\n", err)
			}
			ips.parts = append(ips.parts, aui)
		}
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
		i.SetPartsInfo()
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
