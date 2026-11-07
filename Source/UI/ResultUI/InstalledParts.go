package resultui

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type InstalledParts struct {
	parts []*ui.AnimatableUI

	pos        *Utils.Position
	recordSize image.Point

	gdata game.Gamedata
}

const RecordDistance = 5

func NewInstalledParts(pos Utils.Position, recordSize image.Point, gdata *game.Gamedata) *InstalledParts {
	return &InstalledParts{
		parts:      []*ui.AnimatableUI{},
		pos:        &pos,
		recordSize: recordSize,
		gdata:      *gdata,
	}
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
