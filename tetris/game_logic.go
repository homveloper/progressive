package tetris

import (
	"time"
)

// =============================================================================
// Tetris Game Logic
// =============================================================================

// TetrisGame manages the complete tetris game state and logic
type TetrisGame struct {
	State          *GameState
	eventChannel   chan GameEvent
	running        bool
	dropTimer      *time.Timer
}

// GameEvent represents different game events
type GameEvent struct {
	Type EventType
	Data interface{}
}

// EventType represents the type of game event
type EventType int

const (
	EventMoveLeft EventType = iota
	EventMoveRight
	EventMoveDown
	EventRotate
	EventDrop
	EventPause
	EventRestart
	EventGameTick
)

// =============================================================================
// Game Constructor
// =============================================================================

// NewTetrisGame creates a new tetris game instance
func NewTetrisGame() *TetrisGame {
	game := &TetrisGame{
		State:        NewGameState(),
		eventChannel: make(chan GameEvent, 100),
		running:      false,
	}
	
	return game
}

// =============================================================================
// Game Control Methods
// =============================================================================

// Start begins the game loop
func (g *TetrisGame) Start() {
	if g.running {
		return
	}
	
	g.running = true
	g.State.GameOver = false
	g.State.Paused = false
	
	// Start the game loop (drop timer handled by UI loop)
	go g.gameLoop()
}

// Stop ends the game
func (g *TetrisGame) Stop() {
	g.running = false
	if g.dropTimer != nil {
		g.dropTimer.Stop()
	}
}

// Pause toggles the game pause state
func (g *TetrisGame) Pause() {
	g.State.Paused = !g.State.Paused
	if g.State.Paused {
		if g.dropTimer != nil {
			g.dropTimer.Stop()
		}
	} else {
		g.scheduleNextDrop()
	}
}

// Restart resets the game to initial state
func (g *TetrisGame) Restart() {
	g.Stop()
	g.State = NewGameState()
	g.Start()
}

// =============================================================================
// Game Loop
// =============================================================================

// gameLoop runs the main game loop
func (g *TetrisGame) gameLoop() {
	for g.running && !g.State.GameOver {
		select {
		case event := <-g.eventChannel:
			if !g.State.Paused {
				g.handleEvent(event)
			} else if event.Type == EventPause || event.Type == EventRestart {
				g.handleEvent(event)
			}
		default:
			// Small delay to prevent high CPU usage
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// scheduleNextDrop schedules the next automatic drop
func (g *TetrisGame) scheduleNextDrop() {
	if g.dropTimer != nil {
		g.dropTimer.Stop()
	}
	
	if !g.running || g.State.Paused || g.State.GameOver {
		return
	}
	
	g.dropTimer = time.AfterFunc(time.Duration(g.State.DropInterval)*time.Millisecond, func() {
		g.SendEvent(GameEvent{Type: EventGameTick})
		g.scheduleNextDrop()
	})
}

// =============================================================================
// Event Handling
// =============================================================================

// SendEvent sends an event to the game
func (g *TetrisGame) SendEvent(event GameEvent) {
	select {
	case g.eventChannel <- event:
	default:
		// Channel is full, ignore event
	}
}

// handleEvent processes a game event
func (g *TetrisGame) handleEvent(event GameEvent) {
	switch event.Type {
	case EventMoveLeft:
		g.movePiece(-1, 0)
	case EventMoveRight:
		g.movePiece(1, 0)
	case EventMoveDown:
		g.movePiece(0, 1)
	case EventRotate:
		g.rotatePiece()
	case EventDrop:
		g.dropPiece()
	case EventPause:
		g.Pause()
	case EventRestart:
		g.Restart()
	case EventGameTick:
		g.tick()
	}
}

// =============================================================================
// Game Mechanics
// =============================================================================

// movePiece attempts to move the current piece
func (g *TetrisGame) movePiece(dx, dy int) bool {
	if g.State.CurrentPiece == nil {
		return false
	}
	
	newPos := Position{
		X: g.State.CurrentPiece.Position.X + dx,
		Y: g.State.CurrentPiece.Position.Y + dy,
	}
	
	if g.State.IsValidPosition(g.State.CurrentPiece, newPos, g.State.CurrentPiece.Rotation) {
		g.State.CurrentPiece.Position = newPos
		return true
	}
	
	return false
}

// rotatePiece attempts to rotate the current piece
func (g *TetrisGame) rotatePiece() bool {
	if g.State.CurrentPiece == nil {
		return false
	}
	
	newRotation := (g.State.CurrentPiece.Rotation + 1) % 4
	
	if g.State.IsValidPosition(g.State.CurrentPiece, g.State.CurrentPiece.Position, newRotation) {
		g.State.CurrentPiece.Rotation = newRotation
		g.State.CurrentPiece.Blocks = TetrominoShapes[g.State.CurrentPiece.Type][newRotation]
		return true
	}
	
	// Try wall kicks (simple implementation)
	for _, kick := range []Position{{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: -1}} {
		kickPos := Position{
			X: g.State.CurrentPiece.Position.X + kick.X,
			Y: g.State.CurrentPiece.Position.Y + kick.Y,
		}
		
		if g.State.IsValidPosition(g.State.CurrentPiece, kickPos, newRotation) {
			g.State.CurrentPiece.Position = kickPos
			g.State.CurrentPiece.Rotation = newRotation
			g.State.CurrentPiece.Blocks = TetrominoShapes[g.State.CurrentPiece.Type][newRotation]
			return true
		}
	}
	
	return false
}

// dropPiece drops the current piece to the bottom
func (g *TetrisGame) dropPiece() {
	if g.State.CurrentPiece == nil {
		return
	}
	
	dropDistance := 0
	for g.movePiece(0, 1) {
		dropDistance++
	}
	
	// Add score for hard drop
	g.State.Score += dropDistance * 2
	
	// Place the piece immediately
	g.placePieceAndSpawnNext()
}

// tick handles automatic game progression
func (g *TetrisGame) tick() {
	if g.State.GameOver || g.State.Paused {
		return
	}
	
	// Try to move piece down
	if !g.movePiece(0, 1) {
		// Can't move down, place the piece
		g.placePieceAndSpawnNext()
	}
}

// placePieceAndSpawnNext places the current piece and spawns the next one
func (g *TetrisGame) placePieceAndSpawnNext() {
	if g.State.CurrentPiece == nil {
		return
	}
	
	// Place the piece on the board
	g.State.PlacePiece()
	
	// Clear completed lines
	linesCleared := g.State.ClearLines()
	if linesCleared > 0 {
		g.updateScore(linesCleared)
		g.State.Lines += linesCleared
		g.updateLevel()
	}
	
	// Spawn next piece
	g.State.CurrentPiece = g.State.NextPiece
	g.State.NextPiece = RandomTetromino()
	
	// Check game over
	if g.State.CheckGameOver() {
		g.State.GameOver = true
		g.Stop()
	}
}

// updateScore updates the score based on lines cleared
func (g *TetrisGame) updateScore(linesCleared int) {
	// Standard tetris scoring
	scoreMultiplier := []int{0, 40, 100, 300, 1200} // 0, 1, 2, 3, 4 lines
	if linesCleared > 0 && linesCleared <= 4 {
		g.State.Score += scoreMultiplier[linesCleared] * (g.State.Level + 1)
	}
}

// updateLevel updates the level and drop speed
func (g *TetrisGame) updateLevel() {
	newLevel := g.State.Lines/10 + 1
	if newLevel != g.State.Level {
		g.State.Level = newLevel
		// Increase drop speed (decrease interval)
		g.State.DropInterval = max(50, 1000-((g.State.Level-1)*50))
	}
}

// =============================================================================
// Utility Functions
// =============================================================================

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// =============================================================================
// Game State Queries
// =============================================================================

// GetVisibleBoard returns the board with the current piece overlaid
func (g *TetrisGame) GetVisibleBoard() [][]string {
	// Create a copy of the board
	board := make([][]string, g.State.Board.Height)
	for i := range board {
		board[i] = make([]string, g.State.Board.Width)
		copy(board[i], g.State.Board.Grid[i])
	}
	
	// Add current piece to the visible board
	if g.State.CurrentPiece != nil {
		blocks := g.State.CurrentPiece.GetAbsoluteBlocks()
		for _, block := range blocks {
			if block.Y >= 0 && block.Y < g.State.Board.Height && 
			   block.X >= 0 && block.X < g.State.Board.Width {
				board[block.Y][block.X] = g.State.CurrentPiece.Color
			}
		}
	}
	
	return board
}

// GetGhostPiece returns the position where the current piece would land
func (g *TetrisGame) GetGhostPiece() []Position {
	if g.State.CurrentPiece == nil {
		return nil
	}
	
	// Create a copy of the current piece
	ghostPiece := *g.State.CurrentPiece
	
	// Drop it to the bottom
	for g.State.IsValidPosition(&ghostPiece, Position{X: ghostPiece.Position.X, Y: ghostPiece.Position.Y + 1}, ghostPiece.Rotation) {
		ghostPiece.Position.Y++
	}
	
	return ghostPiece.GetAbsoluteBlocks()
}