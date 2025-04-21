// world/platform.go
package world

import "image"

// Platform represents a solid surface in the game world.
// Currently using image.Rectangle directly, but could be a custom struct later.
type Platform = image.Rectangle

// Screen dimensions needed for platform layout (could be moved to config)
const (
	ScreenWidth  = 640
	ScreenHeight = 480
)


// GetInitialPlatforms returns a slice defining the level's platforms
func GetInitialPlatforms() []Platform {
	return []Platform{
		image.Rect(0, ScreenHeight-40, ScreenWidth, ScreenHeight), // Floor
		image.Rect(200, ScreenHeight-100, 300, ScreenHeight-80),   // Mid platform 1
		image.Rect(400, ScreenHeight-180, 500, ScreenHeight-160),  // Mid platform 2
		image.Rect(50, ScreenHeight-260, 150, ScreenHeight-240),   // High platform
	}
}

// You could add functions here later like LoadLevel(filename string) []Platform