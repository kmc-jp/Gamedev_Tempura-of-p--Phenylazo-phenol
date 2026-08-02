package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Game が ebiten.Game を実装していることを保証するためだけのコード
	var game ebiten.Game = Game{}
	game = game
}

type Game struct {
	ActiveScene Scene
	BackScene   []Scene
}

// Draw implements [ebiten.Game].
func (g Game) Draw(screen *ebiten.Image) {
	panic("unimplemented")
}

// Layout implements [ebiten.Game].
func (g Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	panic("unimplemented")
}

// Update implements [ebiten.Game].
func (g Game) Update() error {
	panic("unimplemented")
}

type Scene interface {
	Init()
	Update(active bool) (nextScene Scene, err error)
	Draw(screen *ebiten.Image)
}
