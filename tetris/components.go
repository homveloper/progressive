package tetris

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// =============================================================================
// Tetris UI Components
// =============================================================================

// TetrisPage represents the main tetris game page
type TetrisPage struct {
	app.Compo
	game *TetrisGame
}

// TetrisBoard renders the game board
type TetrisBoard struct {
	app.Compo
	Board      [][]string
	GhostPiece []Position
	Width      int
	Height     int
}

// TetrisPreview renders the next piece preview
type TetrisPreview struct {
	app.Compo
	Piece *Tetromino
}

// TetrisStats renders game statistics
type TetrisStats struct {
	app.Compo
	Score    int
	Level    int
	Lines    int
	GameOver bool
	Paused   bool
}

// TetrisControls renders game control buttons
type TetrisControls struct {
	app.Compo
	GameOver  bool
	Paused    bool
	OnStart   func()
	OnPause   func()
	OnRestart func()
}

// =============================================================================
// TetrisPage Implementation
// =============================================================================

func (p *TetrisPage) OnMount(ctx app.Context) {
	if p.game == nil {
		p.game = NewTetrisGame()
		p.game.Start() // 게임 자동 시작
		log.Println("Tetris game started!")
	}

	// Start game update loop
	go p.updateLoop(ctx)
}

func (p *TetrisPage) OnDismount() {
	if p.game != nil {
		p.game.Stop()
	}
}

func (p *TetrisPage) updateLoop(ctx app.Context) {
	lastDrop := time.Now()

	for {
		if p.game != nil && p.game.running {
			// Check if it's time for automatic drop
			if time.Since(lastDrop) >= time.Duration(p.game.State.DropInterval)*time.Millisecond {
				if !p.game.State.Paused && !p.game.State.GameOver {
					p.game.SendEvent(GameEvent{Type: EventGameTick})
					lastDrop = time.Now()
					log.Println("Auto drop triggered")
				}
			}

			ctx.Dispatch(func(ctx app.Context) {
				ctx.Update()
			})
		}
		// Update at 60 FPS
		time.Sleep(16 * time.Millisecond) // ~16ms for 60fps
	}
}

func (p *TetrisPage) handleKeyDown(ctx app.Context, e app.Event) {
	if p.game == nil {
		return
	}

	key := e.Get("key").String()

	switch strings.ToLower(key) {
	case "arrowleft", "a":
		p.game.SendEvent(GameEvent{Type: EventMoveLeft})
	case "arrowright", "d":
		p.game.SendEvent(GameEvent{Type: EventMoveRight})
	case "arrowdown", "s":
		p.game.SendEvent(GameEvent{Type: EventMoveDown})
	case "arrowup", "w":
		p.game.SendEvent(GameEvent{Type: EventRotate})
	case " ": // Spacebar for hard drop
		p.game.SendEvent(GameEvent{Type: EventDrop})
		e.PreventDefault()
	case "p":
		p.game.SendEvent(GameEvent{Type: EventPause})
	case "r":
		p.game.SendEvent(GameEvent{Type: EventRestart})
	}
}

func (p *TetrisPage) Render() app.UI {
	if p.game == nil {
		return app.Div().Text("Loading...")
	}

	board := p.game.GetVisibleBoard()
	ghostPiece := p.game.GetGhostPiece()

	return app.Div().Class("tetris-container").
		TabIndex(0).
		OnKeyDown(p.handleKeyDown).
		Body(
			app.H1().Class("tetris-title").Text("🧩 Tetris Game"),

			app.Div().Class("tetris-game").Body(
				// Left panel - Game board
				app.Div().Class("tetris-game-area").Body(
					&TetrisBoard{
						Board:      board,
						GhostPiece: ghostPiece,
						Width:      p.game.State.Board.Width,
						Height:     p.game.State.Board.Height,
					},
				),

				// Right panel - Stats and controls
				app.Div().Class("tetris-sidebar").Body(
					&TetrisStats{
						Score:    p.game.State.Score,
						Level:    p.game.State.Level,
						Lines:    p.game.State.Lines,
						GameOver: p.game.State.GameOver,
						Paused:   p.game.State.Paused,
					},

					&TetrisPreview{
						Piece: p.game.State.NextPiece,
					},

					&TetrisControls{
						GameOver:  p.game.State.GameOver,
						Paused:    p.game.State.Paused,
						OnStart:   func() { p.game.Start() },
						OnPause:   func() { p.game.SendEvent(GameEvent{Type: EventPause}) },
						OnRestart: func() { p.game.SendEvent(GameEvent{Type: EventRestart}) },
					},

					// Controls help
					app.Div().Class("tetris-help").Body(
						app.H4().Text("Controls"),
						app.Div().Class("control-item").Text("← → : Move"),
						app.Div().Class("control-item").Text("↑ : Rotate"),
						app.Div().Class("control-item").Text("↓ : Soft Drop"),
						app.Div().Class("control-item").Text("Space : Hard Drop"),
						app.Div().Class("control-item").Text("P : Pause"),
						app.Div().Class("control-item").Text("R : Restart"),
					),
				),
			),
		)
}

// =============================================================================
// TetrisBoard Implementation
// =============================================================================

func (b *TetrisBoard) Render() app.UI {
	var rows []app.UI

	for y := 0; y < b.Height; y++ {
		var cells []app.UI

		for x := 0; x < b.Width; x++ {
			cellColor := ""
			cellClass := "tetris-cell"

			// Check if this is part of the board
			if y < len(b.Board) && x < len(b.Board[y]) {
				cellColor = b.Board[y][x]
			}

			// Check if this is part of the ghost piece
			if cellColor == "" && b.GhostPiece != nil {
				for _, ghostPos := range b.GhostPiece {
					if ghostPos.X == x && ghostPos.Y == y {
						cellClass += " ghost-piece"
						break
					}
				}
			}

			// Apply color if cell is filled
			if cellColor != "" {
				cellClass += " filled"
			}

			cell := app.Div().Class(cellClass)
			if cellColor != "" {
				cell = cell.Style("background-color", cellColor)
			}

			cells = append(cells, cell)
		}

		rows = append(rows, app.Div().Class("tetris-row").Body(cells...))
	}

	return app.Div().Class("tetris-board").Body(rows...)
}

// =============================================================================
// TetrisPreview Implementation
// =============================================================================

func (p *TetrisPreview) Render() app.UI {
	if p.Piece == nil {
		return app.Div().Class("tetris-preview").Body(
			app.H4().Text("Next Piece"),
			app.Div().Class("preview-area").Text("No piece"),
		)
	}

	// Create a 4x4 grid for the preview
	var rows []app.UI
	for y := 0; y < 4; y++ {
		var cells []app.UI
		for x := 0; x < 4; x++ {
			cellClass := "preview-cell"
			cellColor := ""

			// Check if this position contains a block
			for _, block := range p.Piece.Blocks {
				if block.X+2 == x && block.Y+1 == y { // Center the piece
					cellClass += " filled"
					cellColor = p.Piece.Color
					break
				}
			}

			cell := app.Div().Class(cellClass)
			if cellColor != "" {
				cell = cell.Style("background-color", cellColor)
			}

			cells = append(cells, cell)
		}
		rows = append(rows, app.Div().Class("preview-row").Body(cells...))
	}

	return app.Div().Class("tetris-preview").Body(
		app.H4().Text("Next Piece"),
		app.Div().Class("preview-area").Body(rows...),
	)
}

// =============================================================================
// TetrisStats Implementation
// =============================================================================

func (s *TetrisStats) Render() app.UI {
	statusText := "Playing"
	statusClass := "status-playing"

	if s.GameOver {
		statusText = "Game Over"
		statusClass = "status-gameover"
	} else if s.Paused {
		statusText = "Paused"
		statusClass = "status-paused"
	}

	return app.Div().Class("tetris-stats").Body(
		app.H4().Text("Stats"),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("Status:"),
			app.Span().Class(statusClass).Text(statusText),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("Score:"),
			app.Span().Class("stat-value").Text(fmt.Sprintf("%d", s.Score)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("Level:"),
			app.Span().Class("stat-value").Text(fmt.Sprintf("%d", s.Level)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("Lines:"),
			app.Span().Class("stat-value").Text(fmt.Sprintf("%d", s.Lines)),
		),
	)
}

// =============================================================================
// TetrisControls Implementation
// =============================================================================

func (c *TetrisControls) Render() app.UI {
	var buttons []app.UI

	if c.GameOver {
		buttons = append(buttons,
			app.Button().Class("control-button restart").Text("🔄 Restart").
				OnClick(func(ctx app.Context, e app.Event) {
					if c.OnRestart != nil {
						c.OnRestart()
					}
				}),
		)
	} else {
		if c.Paused {
			buttons = append(buttons,
				app.Button().Class("control-button resume").Text("▶️ Resume").
					OnClick(func(ctx app.Context, e app.Event) {
						if c.OnPause != nil {
							c.OnPause()
						}
					}),
			)
		} else {
			buttons = append(buttons,
				app.Button().Class("control-button pause").Text("⏸️ Pause").
					OnClick(func(ctx app.Context, e app.Event) {
						if c.OnPause != nil {
							c.OnPause()
						}
					}),
			)
		}

		buttons = append(buttons,
			app.Button().Class("control-button restart").Text("🔄 Restart").
				OnClick(func(ctx app.Context, e app.Event) {
					if c.OnRestart != nil {
						c.OnRestart()
					}
				}),
		)
	}

	// Add start button if game is not running
	if c.GameOver {
		buttons = append([]app.UI{
			app.Button().Class("control-button start").Text("🎮 Start Game").
				OnClick(func(ctx app.Context, e app.Event) {
					if c.OnStart != nil {
						c.OnStart()
					}
				}),
		}, buttons...)
	}

	return app.Div().Class("tetris-controls").Body(
		app.H4().Text("Controls"),
		app.Div().Class("control-buttons").Body(buttons...),
	)
}
