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
	spTimer     *tick.Timer
	normalGet   int
}

// レア棚の制御
const (
	prob_Shelf         = 0.001
	particletime_Shelf = 200
	normalmax_Shelf    = 5
)

func (s *ShelfTile) Interact(player *Player) error {
	var foundPart *game.Part = nil

	if s.isSparkling {
		// キラキラ表示があれば停止
		s.render.EnableSubUpdate(false)
		foundPart = parts.RandomPartsFromAll(false)
		s.normalGet = 0
		s.isSparkling = false
		s.spTimer.Set(0)
	} else if s.normalGet < normalmax_Shelf {
		foundPart = parts.RandomPartsFromAll(true)
		s.normalGet++

		if s.normalGet == normalmax_Shelf {
			// 棚が品切れになった時の処理をどうするか
		}
	}

	s.scene.Gamedata.AddPartToFoundList(foundPart)
	s.scene.ShowFoundPartLog(foundPart)
	s.scene.ShowPartPopup(foundPart, player.tw.tf.Position())

	return nil
}

func (s *ShelfTile) Update(active bool) error {
	if !active {
		s.spTimer.Suspend()
		return nil
	}
	s.spTimer.Resume()

	if !s.isSparkling && Utils.WithProbability(prob_Shelf) {
		s.render.EnableSubUpdate(true)

		s.isSparkling = true
		s.spTimer.Set(particletime_Shelf)
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
			spTimer:     tick.NewTimer(),
			isSparkling: false,
		}, nil
	}
}

// test
var _ InteractableTile = (*ShelfTile)(nil)
