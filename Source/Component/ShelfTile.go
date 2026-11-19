package component

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	parts "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/Parts"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	tick "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Tick"
)

type ShelfTile struct {
	scene       *scene.PlayScene
	render      *MultiRender
	isSparkling bool
	spTimer     tick.Timer
}

// レア棚の制御
const (
	prob         = 0.001
	particletime = 200
)

func (s ShelfTile) Interact(player *Player) error {
	var foundPart *game.Part

	if s.isSparkling {
		// キラキラ表示があれば停止
		s.render.EnableSubUpdate(false)
		foundPart = parts.RandomPartsFromAll(false)
	} else {
		foundPart = parts.RandomPartsFromAll(true)
	}

	s.scene.Gamedata.AddPartToFoundList(foundPart)
	s.scene.ShowFoundPartLog(foundPart)
	s.scene.ShowPartPopup(foundPart, player.tw.tf.Position())

	return nil
}

func (s *ShelfTile) Update(active bool) error {
	if !s.isSparkling && Utils.WithProbability(prob) {
		s.render.EnableSubUpdate(true)

		s.isSparkling = true
		s.spTimer.Set(particletime)
	} else if s.isSparkling && s.spTimer.Expired() {
		s.render.EnableSubUpdate(false)
		s.isSparkling = false
	}

	return nil
}

func NewShelfTile(scene *scene.PlayScene, render *MultiRender) Utils.Factory[*ShelfTile] {
	return func() (*ShelfTile, error) {
		render.EnableSubRendering(true)
		return &ShelfTile{
			scene:       scene,
			render:      render,
			isSparkling: false,
		}, nil
	}
}

// test
var _ InteractableTile = (*ShelfTile)(nil)
