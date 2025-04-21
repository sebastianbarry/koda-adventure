// player/player.go
package player

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Player-specific constants (exported if needed by other packages, like game for drawing)
const (
	Width   = 16
	Height  = 32
	moveSpeed = 3
	jumpPower = 7   // Initial upward velocity on jump
	gravity   = 0.4 // Acceleration due to gravity per frame
	maxFallSpeed = 8   // Terminal velocity
)

// Player represents the character controlled by the user
type Player struct {
	X, Y   float64 // Position (top-left corner) - Exported
	Vx, Vy float64 // Velocity - Exported
	OnGround bool    // Is the player currently standing on something? - Exported
}

// NewPlayer creates a new player instance at a default position
func NewPlayer() *Player {
	return &Player{X: 100, Y: 100} // Initial position
}

// Rect returns the player's current bounding box
func (p *Player) Rect() image.Rectangle {
	return image.Rect(int(p.X), int(p.Y), int(p.X)+Width, int(p.Y)+Height)
}

// NextRectX calculates the bounding box based on potential next horizontal position
func (p *Player) NextRectX() image.Rectangle {
	nextX := p.X + p.Vx
	return image.Rect(int(nextX), int(p.Y), int(nextX)+Width, int(p.Y)+Height)
}

// NextRectY calculates the bounding box based on potential next vertical position
func (p *Player) NextRectY() image.Rectangle {
	nextY := p.Y + p.Vy
	return image.Rect(int(p.X), int(nextY), int(p.X)+Width, int(p.Y)+Height)
}

// HandleInput updates player velocity based on keyboard input
func (p *Player) HandleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Vx = -moveSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Vx = moveSpeed
	} else {
		p.Vx = 0 // Stop horizontal movement if no key is pressed
	}

	// Jumping
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && p.OnGround {
		p.Vy = -jumpPower
		p.OnGround = false // Player is now airborne
	}
}

// ApplyGravity updates the vertical velocity
func (p *Player) ApplyGravity() {
	p.Vy += gravity
	// Clamp fall speed
	if p.Vy > maxFallSpeed {
		p.Vy = maxFallSpeed
	}
}

// UpdatePosition actually changes the player's X and Y coordinates based on velocity
func (p *Player) UpdatePosition() {
	p.X += p.Vx
	p.Y += p.Vy
}

// ResetPosition puts the player back at a starting point (e.g., after falling)
func (p *Player) ResetPosition() {
	p.X = 100
	p.Y = 100
	p.Vy = 0
	p.Vx = 0
	p.OnGround = false
}

// HandleCollisionX stops horizontal movement on collision
func (p *Player) HandleCollisionX() {
	p.Vx = 0
}

// HandleCollisionY handles vertical collision response
func (p *Player) HandleCollisionY(platformY float64) {
	// Determine if landing or hitting ceiling based on current velocity
	if p.Vy > 0 { // Moving down / landing
		p.Y = platformY - Height // Snap player top to the platform's top
		p.OnGround = true        // Landed on a platform
	} else if p.Vy < 0 { // Moving up / hitting ceiling
		p.Y = platformY // Snap player bottom to the platform's bottom (platformY is p.Max.Y)
	}
	p.Vy = 0 // Stop vertical movement
}