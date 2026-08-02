package main

import (
	"log"

	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"

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
	BackScene   Utils.Stack[Scene]
}

// Draw implements [ebiten.Game].
func (g Game) Draw(screen *ebiten.Image) {
	panic("unimplemented")
}

// Layout implements [ebiten.Game].
func (g Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}

// Update implements [ebiten.Game].
func (g Game) Update() error {
	panic("unimplemented")
}

type Scene interface {
	Name() string
	Init()
	Update(active bool) (nextScene Scene, asNewScene bool, err error)
	Draw(screen *ebiten.Image)
}
