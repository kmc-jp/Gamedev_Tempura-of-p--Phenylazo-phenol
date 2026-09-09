package scene

import (
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	transition "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/SceneTransitionType"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PauseScene struct{}

func (s PauseScene) Update(active bool) (Utils.Factory[game.Scene], transition.Type, error) {
	// Escキーでポーズ解除（Pop）
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return nil, transition.Pop(), nil
	}
	return nil, transition.None(), nil
}

func (s PauseScene) Name() string {
	return "PauseScene"
}

func (s PauseScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	ebitenutil.DebugPrint(screen, "pause\nEsc to return")
}
