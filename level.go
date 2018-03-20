package spritengine

import (
	"image"
	"image/color"
	"image/draw"
)

// Level is a struct that defines a single level of a game
type Level struct {
	BackgroundColour color.RGBA
	Gravity          float64
	GameObjects      []*GameObject
	Game             *Game
}

// Repaint redraws the entire level for a new frame
func (level *Level) Repaint(stage *image.RGBA) {

	// Figure out where all the floor objects are
	level.AssignFloors()

	// Figure out which objects are colliding
	level.AssignCollisions()

	// Paint the background colour
	draw.Draw(stage, stage.Bounds(), &image.Uniform{level.BackgroundColour}, image.ZP, draw.Src)

	// Update each game object
	for _, gameObject := range level.GameObjects {

		gameObject.Level = level
		gameObject.RecalculatePosition(level.Gravity)

		if gameObject.Direction == DirLeft {
			gameObject.Flipped = true
		} else if gameObject.Direction == DirRight {
			gameObject.Flipped = false
		}

		// 0 is at the bottom, so flip the Y axis to paint correctly
		invertedY := level.Game.Height - int(gameObject.Position.Y) - gameObject.Height()

		gameObject.CurrentSprite().AddToCanvas(stage, int(gameObject.Position.X), invertedY, gameObject.Flipped)

	}

}

// AssignFloors iterates through all objects in the level and defines which
// object beneath them (if any) should be considered their 'floor' object,
// setting its top edge as the lowest point that the object can fall
func (level *Level) AssignFloors() {

	floorXCoords := map[int][]*GameObject{}

	// Make a map of each object's possible X positions
	for _, gameObject := range level.GameObjects {

		// Skip non-interactive objects
		if gameObject.Interactive == false {
			continue
		}

		for i := 0; i < gameObject.Width(); i++ {

			xPos := i + int(gameObject.Position.X)
			floorXCoords[xPos] = append(floorXCoords[xPos], gameObject)

		}

	}

	// Find the objects that sit beneath every other object
	for _, gameObject := range level.GameObjects {

		// Skip objects that float
		if gameObject.Mass == 0 {
			continue
		}

		highestFloorObject := float64(0 - gameObject.Height())

		for i := 0; i < gameObject.Width(); i++ {

			xPos := i + int(gameObject.Position.X)

			if floorObjects, ok := floorXCoords[xPos]; ok {

				// Find the one that is highest while still being lower than
				// the object itself
				for _, floorObject := range floorObjects {

					floorObjectTop := (floorObject.Position.Y + float64(floorObject.Height()))

					if floorObjectTop <= gameObject.Position.Y {

						if floorObjectTop > highestFloorObject {
							highestFloorObject = floorObjectTop
						}

					}

				}

			}

		}

		gameObject.FloorY = highestFloorObject

	}

}

// AssignCollisions iterates through all objects in the level and defines which
// objects (if any) intersect them
func (level *Level) AssignCollisions() {

	xCoords := map[int][]*GameObject{}

	// Make a map of each object's possible X positions
	for _, gameObject := range level.GameObjects {

		// Skip non-interactive objects
		if gameObject.Interactive == false {
			continue
		}

		for i := 0; i < gameObject.Width(); i++ {
			xPos := i + int(gameObject.Position.X)
			xCoords[xPos] = append(xCoords[xPos], gameObject)
		}

	}

	// Find the objects that also intersect on the Y axis
	for _, gameObject := range level.GameObjects {

		intersections := map[*GameObject]bool{}
		gameObjectYMin := gameObject.Position.Y
		gameObjectYMax := gameObjectYMin + float64(gameObject.Height())

		for i := 0; i < gameObject.Width(); i++ {

			xPos := i + int(gameObject.Position.X)

			if intersectingObjects, ok := xCoords[xPos]; ok {

				for _, intersectingObject := range intersectingObjects {

					// Ignore the object itself
					if intersectingObject == gameObject {
						continue
					}

					// Skip the object if it has already been stored
					if _, ok := intersections[intersectingObject]; ok {
						continue
					}

					intersectingObjectYMin := intersectingObject.Position.Y
					intersectingObjectYMax := intersectingObjectYMin + float64(intersectingObject.Height())

					if (gameObjectYMin >= intersectingObjectYMax || gameObjectYMax <= intersectingObjectYMin) == false {
						intersections[intersectingObject] = true
					}

				}

			}

		}

		// Let the game object know that there have been collisions
		if len(intersections) > 0 {

			for collidingObject := range intersections {

				gameObject.CollisionHandler(gameObject, Collision{
					GameObject: collidingObject,
					Edge:       EdgeTop, // TODO: Figure out which edge the collision actually occurred on
				})

			}

		}

	}

}
