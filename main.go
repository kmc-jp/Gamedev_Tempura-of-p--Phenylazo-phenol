package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// ゲーム起動
	game := Game{}
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
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
