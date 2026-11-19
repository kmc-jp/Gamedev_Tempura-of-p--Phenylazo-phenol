package component

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	parts "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/Parts"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
)

type ShelfTile struct {
	scene       *scene.PlayScene
	render      *MultiRender
	isSparkling bool
}

func (s ShelfTile) Interact(player *Player) error {
	// キラキラ表示があれば停止
	s.render.EnableSubUpdate(false)

	foundPart := parts.RandomPartsFromAll()
	s.scene.Gamedata.AddPartToFoundList(foundPart)
	s.scene.ShowFoundPartLog(foundPart)
	s.scene.ShowPartPopup(foundPart, player.tw.tf.Position())

	return nil
}

func (s *ShelfTile) Update(active bool) error {
	const prob = 0.001
	if Utils.WithProbability(prob) {
		s.render.EnableSubRendering(true)
		s.render.EnableSubUpdate(true)
		s.isSparkling = true
	}

	return nil
}

func NewShelfTile(scene *scene.PlayScene, render *MultiRender) Utils.Factory[*ShelfTile] {
	return func() (*ShelfTile, error) {
		return &ShelfTile{
			scene:       scene,
			render:      render,
			isSparkling: false,
		}, nil
	}
}

// test
var _ InteractableTile = (*ShelfTile)(nil)
