package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type RegisterTile struct {
	scene *scene.PlayScene
}

func (s RegisterTile) Interact(player *Player) error {
	fmt.Printf("Interact with Regester.\n")
	if s.scene.Gamedata.State == game.Play_Explore {
		s.scene.Gamedata.State = game.Play_Register
	} else {
		return fmt.Errorf("Interactした時点のstateが不正な値です\n")
	}
	return nil
}

func (s RegisterTile) Update(active bool) error {
	// なにもしない
	return nil
}

func NewRegisterTile(scene *scene.PlayScene, render gameobject.Render, position Utils.Position) Utils.Factory[*RegisterTile] {
	return func() (*RegisterTile, error) {
		// レジ台の上に「お会計」吹き出しを浮かべる
		img, _, _ := render.Draw()
		if img != nil {
			top := Utils.Position{X: position.X, Y: position.Y - float64(img.Bounds().Dy())/2}
			if err := scene.AddRegisterBalloon(top); err != nil {
				return nil, fmt.Errorf("AddRegisterBalloon() でエラーが発生しました。\n%w", err)
			}
		}
		return &RegisterTile{
			scene: scene,
		}, nil
	}
}

// test
var _ InteractableTile = RegisterTile{}
