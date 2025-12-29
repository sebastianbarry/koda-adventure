// main.go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	// Use your actual module path here
	"github.com/sebastianbarry/koda-adventure/game"
)

const (
	// Keep window scaling factor here if desired, or move to game/config
	windowScaleFactor = 2
)

func main() {
	// Initialize the game state from the 'game' package
	g, err := game.NewGame()
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	// Set window properties (using constants from the game package)
	ebiten.SetWindowSize(game.ScreenWidth*windowScaleFactor, game.ScreenHeight*windowScaleFactor)
	ebiten.SetWindowTitle("Simple Go Platformer (Refactored)")

	// Run the game loop
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}