package spritengine

// Game is a struct that defines a game and the window that contains it
type Game struct {
	Title           string
	Width           int
	Height          int
	ScaleFactor     int
	TargetFrameRate int
	FramePainter    FramePainter
	KeyListener     KeyListener
	Levels          []*Level
	CurrentLevelID  int
}

// CreateGame sets up a game and its window
func CreateGame(title string, width int, height int, scaleFactor int, targetFrameRate int, framePainter FramePainter, keyListener KeyListener, levels []*Level) *Game {

	game := Game{
		Title:           title,
		Width:           width,
		Height:          height,
		ScaleFactor:     scaleFactor,
		TargetFrameRate: targetFrameRate,
		FramePainter:    framePainter,
		KeyListener:     keyListener,
		Levels:          levels,
		CurrentLevelID:  0,
	}

	for _, level := range levels {
		level.Game = &game
	}

	createWindow(&game)

	return &game

}

// CurrentLevel gets the current Level object
func (game *Game) CurrentLevel() *Level {

	return game.Levels[game.CurrentLevelID]

}
