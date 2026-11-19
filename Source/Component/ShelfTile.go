package component

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	parts "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/Parts"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type ShelfTile struct {
	scene *scene.PlayScene
}

func (s ShelfTile) Interact(player *Player) error {
	foundPart := parts.RandomPartsFromAll()
	if foundPart != nil {
		if foundPart.Kind != game.Trash {
			s.scene.MessageWindowDisplayTimer.Set(20)
			s.scene.Canvas().SetMessageWindowText(fmt.Sprintf("%s を見つけた", foundPart.Name))

		}
	}
	s.scene.Gamedata.AddPartToFoundList(foundPart)
	s.scene.ShowFoundPartLog(foundPart)
	s.scene.ShowPartPopup(foundPart, player.tw.tf.Position())

	return nil

}

func (s ShelfTile) Update(active bool) error {
	// なにもしない
	return nil
}

func NewShelfTile(scene *scene.PlayScene) Utils.Factory[*ShelfTile] {
	return func() (*ShelfTile, error) {
		return &ShelfTile{
			scene: scene,
		}, nil
	}
}

// test
var _ InteractableTile = ShelfTile{}
