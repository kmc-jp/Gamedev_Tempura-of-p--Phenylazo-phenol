package resultui

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type InstalledParts struct {
	parts []*ui.AnimatableUI

	pos        *Utils.Position
	recordSize image.Point

	gdata game.Gamedata
}

const (
	RecordDistance = 5
	TextSize       = 20
	UIsOffset      = 300 //ui移動開始前の位置差分
)

func NewInstalledParts(pos Utils.Position, recordSize image.Point, gdata *game.Gamedata) *InstalledParts {
	return &InstalledParts{
		parts:      []*ui.AnimatableUI{},
		pos:        &pos,
		recordSize: recordSize,
		gdata:      *gdata,
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
		}
	}
	return nil
}

func (i *InstalledParts) Draw(screen *ebiten.Image) {
	for _, v := range i.parts {
		v.Draw(screen)
	}
}

func (i *InstalledParts) Update(active bool) error {
	for _, v := range i.parts {
		v.Update(true)
	}
	return nil
}
