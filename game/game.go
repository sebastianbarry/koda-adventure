// game/game.go
package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	// Use your actual module path here
	"github.com/sebastianbarry/koda-adventure/player"
	"github.com/sebastianbarry/koda-adventure/world"
)

// Export ScreenWidth and ScreenHeight so main.go can use them for window setup
const (
	ScreenWidth  = world.ScreenWidth  // Reuse from world or define centrally
	ScreenHeight = world.ScreenHeight // Reuse from world or define centrally
)

// Game implements ebiten.Game interface.
type Game struct {
	player    *player.Player    // Embed the Player struct (pointer)
	platforms []world.Platform  // Slice of platforms from world package
	whiteImage *ebiten.Image    // Helper for drawing shapes
}

// NewGame initializes the game state
func NewGame() (*Game, error) {
	g := &Game{}
	g.player = player.NewPlayer()
	g.platforms = world.GetInitialPlatforms()

	// Create a 1x1 white image for drawing shapes
	g.whiteImage = ebiten.NewImage(1, 1)
	g.whiteImage.Fill(color.White)

	return g, nil
}

// Update proceeds the game state.
func (g *Game) Update() error {
	// 1. Handle Player Input
	g.player.HandleInput()

	// 2. Apply Physics (Gravity)
	g.player.ApplyGravity()

	// 3. Collision Detection & Resolution
	g.player.OnGround = false // Assume not on ground until collision check proves otherwise

	// Calculate potential next bounding boxes from player's perspective
	nextRectX := g.player.NextRectX()
	nextRectY := g.player.NextRectY()

	// Check Horizontal Collision
	collidedX := false
	for _, p := range g.platforms {
		if nextRectX.Overlaps(p) {
			g.player.HandleCollisionX() // Tell player to stop horizontal movement
			collidedX = true
			break
		}
	}

	// Check Vertical Collision
	collidedY := false
	for _, p := range g.platforms {
		if nextRectY.Overlaps(p) {
			// Tell player how to react based on collision point
			if g.player.Vy > 0 { // Landing
				g.player.HandleCollisionY(float64(p.Min.Y))
			} else if g.player.Vy < 0 { // Hitting ceiling
				g.player.HandleCollisionY(float64(p.Max.Y))
			}
			collidedY = true
			break
		}
	}

	// 4. Update Player Position (only apply velocity if no collision occurred on that axis)
	// Note: A more robust approach might involve calculating the exact collision time/point,
	// but for simplicity, we just stop movement if a collision *would* happen on the next step.
	if !collidedX {
		g.player.X += g.player.Vx
	}
	if !collidedY {
		g.player.Y += g.player.Vy
	}


	// 5. Screen Boundaries
	if g.player.X < 0 {
		g.player.X = 0
		g.player.Vx = 0 // Stop movement at edge
	} else if g.player.X+player.Width > ScreenWidth {
		g.player.X = ScreenWidth - player.Width
		g.player.Vx = 0 // Stop movement at edge
	}

	// Reset if falls off bottom
	if g.player.Y > ScreenHeight {
		g.player.ResetPosition()
	}

	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{R: 30, G: 30, B: 50, A: 255})

	// Draw platforms
	platformColor := color.RGBA{R: 100, G: 100, B: 120, A: 255}
	for _, p := range g.platforms {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(p.Dx()), float64(p.Dy())) // Scale the 1x1 white image
		op.GeoM.Translate(float64(p.Min.X), float64(p.Min.Y))
		op.ColorScale.ScaleWithColor(platformColor)
		screen.DrawImage(g.whiteImage, op)
	}

	// Draw player
	playerColor := color.RGBA{R: 255, G: 80, B: 80, A: 255}
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Scale(player.Width, player.Height) // Use constants from player package
	playerOp.GeoM.Translate(g.player.X, g.player.Y)
	playerOp.ColorScale.ScaleWithColor(playerColor)
	screen.DrawImage(g.whiteImage, playerOp)

	// Debug Info
	debugMsg := fmt.Sprintf("X: %.1f, Y: %.1f\nVX: %.1f, VY: %.1f\nOnGround: %t",
		g.player.X, g.player.Y, g.player.Vx, g.player.Vy, g.player.OnGround)
	ebitenutil.DebugPrint(screen, debugMsg)
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}