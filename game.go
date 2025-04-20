package main

import (
	"image"
	"image/color"
	"log"
	//"math"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil" // Use inpututil for key presses
)

const (
	screenWidth  = 640
	screenHeight = 480

	playerWidth   = 16
	playerHeight  = 32
	playerMoveSpeed = 3
	playerJumpPower = 7 // Initial upward velocity on jump
	gravity       = 0.4 // Acceleration due to gravity per frame
	maxFallSpeed  = 8   // Terminal velocity
)

// Player represents the character controlled by the user
type Player struct {
	x, y   float64 // Position (top-left corner)
	vx, vy float64 // Velocity
	onGround bool    // Is the player currently standing on something?
}

// Rect returns the player's bounding box as an image.Rectangle
func (p *Player) Rect() image.Rectangle {
	return image.Rect(int(p.x), int(p.y), int(p.x)+playerWidth, int(p.y)+playerHeight)
}

// Game implements ebiten.Game interface.
type Game struct {
	player    Player
	platforms []image.Rectangle // Use image.Rectangle for simple AABB collision
	// A simple white image for drawing player/platforms
	// In a real game, you'd use sprites (loaded images)
	whiteImage *ebiten.Image
}

// NewGame initializes the game state
func NewGame() *Game {
	g := &Game{}
	g.player = Player{x: 100, y: 100} // Initial position

	// Create platforms (x, y, width, height)
	g.platforms = []image.Rectangle{
		image.Rect(0, screenHeight-40, screenWidth, screenHeight), // Floor
		image.Rect(200, screenHeight-100, 300, screenHeight-80),   // Mid platform 1
		image.Rect(400, screenHeight-180, 500, screenHeight-160),  // Mid platform 2
		image.Rect(50, screenHeight-260, 150, screenHeight-240),   // High platform
	}

	// Create a 1x1 white image for drawing shapes
	g.whiteImage = ebiten.NewImage(1, 1)
	g.whiteImage.Fill(color.White)

	return g
}

// Update proceeds the game state.
// Update is called every frame (1/60 second).
func (g *Game) Update() error {
	// --- Input Handling ---
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.vx = -playerMoveSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.vx = playerMoveSpeed
	} else {
		g.player.vx = 0 // Stop horizontal movement if no key is pressed
	}

	// Jumping - Use inpututil.IsKeyJustPressed for single trigger
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.onGround {
		g.player.vy = -playerJumpPower
		g.player.onGround = false // Player is now airborne
	}

	// --- Physics ---
	// Apply gravity
	g.player.vy += gravity
	// Clamp fall speed
	if g.player.vy > maxFallSpeed {
		g.player.vy = maxFallSpeed
	}

	// --- Collision Detection & Resolution ---
	// Reset onGround status before checking collisions for the current frame
	g.player.onGround = false

	// Calculate potential next position
	nextX := g.player.x + g.player.vx
	nextY := g.player.y + g.player.vy

	// Create potential future bounding boxes
	nextRectX := image.Rect(int(nextX), int(g.player.y), int(nextX)+playerWidth, int(g.player.y)+playerHeight)
	nextRectY := image.Rect(int(g.player.x), int(nextY), int(g.player.x)+playerWidth, int(nextY)+playerHeight)

	// Check Horizontal Collision
	collidedX := false
	for _, p := range g.platforms {
		if nextRectX.Overlaps(p) {
			// Collision detected horizontally
			g.player.vx = 0 // Stop horizontal movement
			collidedX = true
			break // Stop checking after first collision for this axis
		}
	}
	// Update horizontal position only if no collision
	if !collidedX {
		g.player.x += g.player.vx
	}

	// Check Vertical Collision
	collidedY := false
	for _, p := range g.platforms {
		if nextRectY.Overlaps(p) {
			// Collision detected vertically
			if g.player.vy > 0 { // Moving down / landing
				// Snap player top to the platform's top
				g.player.y = float64(p.Min.Y - playerHeight)
				g.player.onGround = true // Landed on a platform
			} else if g.player.vy < 0 { // Moving up / hitting ceiling
				// Snap player bottom to the platform's bottom
				g.player.y = float64(p.Max.Y)
			}
			g.player.vy = 0 // Stop vertical movement
			collidedY = true
			break // Stop checking after first collision for this axis
		}
	}
	// Update vertical position only if no collision
	if !collidedY {
		g.player.y += g.player.vy
	}

	// --- Screen Boundaries (Simple) ---
	// Prevent moving off left/right edges
	if g.player.x < 0 {
		g.player.x = 0
	} else if g.player.x+playerWidth > screenWidth {
		g.player.x = screenWidth - playerWidth
	}

	// Optional: Reset if falls off bottom (or handle game over)
	if g.player.y > screenHeight {
		g.player.x = 100
		g.player.y = 100
		g.player.vy = 0
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60 second) after Update.
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen (optional, depending on desired background)
	screen.Fill(color.RGBA{R: 30, G: 30, B: 50, A: 255}) // Dark blue background

	// Draw platforms
	platformColor := color.RGBA{R: 100, G: 100, B: 120, A: 255} // Grey
	for _, p := range g.platforms {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(p.Dx()), float64(p.Dy())) // Scale the 1x1 white image
		op.GeoM.Translate(float64(p.Min.X), float64(p.Min.Y))
		op.ColorScale.ScaleWithColor(platformColor)
		screen.DrawImage(g.whiteImage, op)

		// Debug: Draw platform outlines (optional)
		// ebitenutil.DrawRect(screen, float64(p.Min.X), float64(p.Min.Y), float64(p.Dx()), float64(p.Dy()), color.White)
	}

	// Draw player
	playerColor := color.RGBA{R: 255, G: 80, B: 80, A: 255} // Reddish
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Scale(playerWidth, playerHeight) // Scale the 1x1 white image
	playerOp.GeoM.Translate(g.player.x, g.player.y)
	playerOp.ColorScale.ScaleWithColor(playerColor)
	screen.DrawImage(g.whiteImage, playerOp)

	// Debug Info (Optional)
	// debugMsg := fmt.Sprintf("X: %.1f, Y: %.1f\nVX: %.1f, VY: %.1f\nOnGround: %t",
	// 	g.player.x, g.player.y, g.player.vx, g.player.vy, g.player.onGround)
	// ebitenutil.DebugPrint(screen, debugMsg)
}

// Layout takes the outside size (e.g., window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2) // Make window bigger for visibility
	ebiten.SetWindowTitle("Simple Go Platformer")

	game := NewGame()

	// RunGame starts the game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}