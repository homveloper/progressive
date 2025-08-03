package tetris

import (
	"math/rand"
	"time"
)

// =============================================================================
// Tetris Game Models
// =============================================================================

// TetrominoType represents different types of tetris blocks
type TetrominoType int

const (
	TypeI TetrominoType = iota // I-piece (line)
	TypeO                      // O-piece (square)
	TypeT                      // T-piece
	TypeS                      // S-piece
	TypeZ                      // Z-piece
	TypeJ                      // J-piece
	TypeL                      // L-piece
)

// Position represents a coordinate on the game board
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Tetromino represents a falling tetris piece
type Tetromino struct {
	Type     TetrominoType `json:"type"`
	Position Position      `json:"position"`
	Blocks   []Position    `json:"blocks"`   // Relative positions of blocks
	Rotation int           `json:"rotation"` // 0, 1, 2, 3 for different rotations
	Color    string        `json:"color"`
}

// GameBoard represents the tetris playing field
type GameBoard struct {
	Width  int        `json:"width"`
	Height int        `json:"height"`
	Grid   [][]string `json:"grid"` // Color strings for each cell ("" = empty)
}

// GameState represents the current state of the tetris game
type GameState struct {
	Board        GameBoard  `json:"board"`
	CurrentPiece *Tetromino `json:"currentPiece,omitempty"`
	NextPiece    *Tetromino `json:"nextPiece,omitempty"`
	Score        int        `json:"score"`
	Level        int        `json:"level"`
	Lines        int        `json:"lines"`
	GameOver     bool       `json:"gameOver"`
	Paused       bool       `json:"paused"`
	DropInterval int        `json:"dropInterval"` // milliseconds
	LastDrop     time.Time  `json:"lastDrop"`
}

// =============================================================================
// Tetromino Definitions
// =============================================================================

// TetrominoShapes defines the block patterns for each tetromino type
var TetrominoShapes = map[TetrominoType][4][]Position{
	TypeI: {
		// Rotation 0 (horizontal)
		{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
		// Rotation 1 (vertical)
		{{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
		// Rotation 2 (horizontal)
		{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
		// Rotation 3 (vertical)
		{{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
	},
	TypeO: {
		// Square - same for all rotations
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
	},
	TypeT: {
		// T-piece rotations
		{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}},  // Up
		{{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}},  // Right
		{{X: 0, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}, // Down
		{{X: 0, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}}, // Left
	},
	TypeS: {
		// S-piece rotations
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}},
		{{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}},
		{{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
	},
	TypeZ: {
		// Z-piece rotations
		{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
		{{X: 1, Y: -1}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}},
		{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
		{{X: 1, Y: -1}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}},
	},
	TypeJ: {
		// J-piece rotations
		{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: -1, Y: 1}},
		{{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
		{{X: 1, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}},
		{{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1}},
	},
	TypeL: {
		// L-piece rotations
		{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
		{{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 1}},
		{{X: -1, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}},
		{{X: 1, Y: -1}, {X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1}},
	},
}

// TetrominoColors defines the color for each tetromino type
var TetrominoColors = map[TetrominoType]string{
	TypeI: "#00ffff", // Cyan
	TypeO: "#ffff00", // Yellow
	TypeT: "#800080", // Purple
	TypeS: "#00ff00", // Green
	TypeZ: "#ff0000", // Red
	TypeJ: "#0000ff", // Blue
	TypeL: "#ffa500", // Orange
}

// =============================================================================
// Constructor Functions
// =============================================================================

// NewGameBoard creates a new empty game board
func NewGameBoard(width, height int) GameBoard {
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
	}

	return GameBoard{
		Width:  width,
		Height: height,
		Grid:   grid,
	}
}

// NewTetromino creates a new tetromino of the specified type
func NewTetromino(pieceType TetrominoType) *Tetromino {
	return &Tetromino{
		Type:     pieceType,
		Position: Position{X: 5, Y: 0},          // Start at top center
		Blocks:   TetrominoShapes[pieceType][0], // Start with rotation 0
		Rotation: 0,
		Color:    TetrominoColors[pieceType],
	}
}

// RandomTetromino creates a random tetromino
func RandomTetromino() *Tetromino {
	types := []TetrominoType{TypeI, TypeO, TypeT, TypeS, TypeZ, TypeJ, TypeL}
	randomType := types[rand.Intn(len(types))]
	return NewTetromino(randomType)
}

// NewGameState creates a new game state with default settings
func NewGameState() *GameState {
	board := NewGameBoard(10, 20) // Standard tetris dimensions

	return &GameState{
		Board:        board,
		CurrentPiece: RandomTetromino(),
		NextPiece:    RandomTetromino(),
		Score:        0,
		Level:        1,
		Lines:        0,
		GameOver:     false,
		Paused:       false,
		DropInterval: 1000, // 1 second
		LastDrop:     time.Now(),
	}
}

// =============================================================================
// Helper Methods
// =============================================================================

// GetAbsoluteBlocks returns the absolute positions of all blocks in the tetromino
func (t *Tetromino) GetAbsoluteBlocks() []Position {
	var blocks []Position
	for _, block := range t.Blocks {
		blocks = append(blocks, Position{
			X: t.Position.X + block.X,
			Y: t.Position.Y + block.Y,
		})
	}
	return blocks
}

// IsValidPosition checks if the tetromino can be placed at the given position
func (g *GameState) IsValidPosition(tetromino *Tetromino, newPos Position, newRotation int) bool {
	// Get the shape for the new rotation
	shape := TetrominoShapes[tetromino.Type][newRotation]

	// Check each block of the tetromino
	for _, block := range shape {
		x := newPos.X + block.X
		y := newPos.Y + block.Y

		// Check bounds
		if x < 0 || x >= g.Board.Width || y >= g.Board.Height {
			return false
		}

		// Don't check negative Y (above board)
		if y < 0 {
			continue
		}

		// Check collision with existing blocks
		if g.Board.Grid[y][x] != "" {
			return false
		}
	}

	return true
}

// PlacePiece places the current piece on the board permanently
func (g *GameState) PlacePiece() {
	if g.CurrentPiece == nil {
		return
	}

	blocks := g.CurrentPiece.GetAbsoluteBlocks()
	for _, block := range blocks {
		if block.Y >= 0 && block.Y < g.Board.Height &&
			block.X >= 0 && block.X < g.Board.Width {
			g.Board.Grid[block.Y][block.X] = g.CurrentPiece.Color
		}
	}
}

// ClearLines removes completed lines and returns the number of lines cleared
func (g *GameState) ClearLines() int {
	linesCleared := 0

	for y := g.Board.Height - 1; y >= 0; y-- {
		// Check if line is complete
		complete := true
		for x := 0; x < g.Board.Width; x++ {
			if g.Board.Grid[y][x] == "" {
				complete = false
				break
			}
		}

		if complete {
			// Remove the line
			copy(g.Board.Grid[1:y+1], g.Board.Grid[0:y])
			// Clear the top line
			for x := 0; x < g.Board.Width; x++ {
				g.Board.Grid[0][x] = ""
			}
			linesCleared++
			y++ // Check the same line again since we moved everything down
		}
	}

	return linesCleared
}

// CheckGameOver checks if the game is over
func (g *GameState) CheckGameOver() bool {
	if g.CurrentPiece == nil {
		return false
	}

	// Check if the current piece can be placed at its starting position
	return !g.IsValidPosition(g.CurrentPiece, g.CurrentPiece.Position, g.CurrentPiece.Rotation)
}
