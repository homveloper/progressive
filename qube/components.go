package qube

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// =============================================================================
// 메인 3D Qube 게임 페이지
// =============================================================================

// QubePage represents the main 3D Qube game page
type QubePage struct {
	app.Compo
	game         *QubeGame
	keys         map[string]bool
	mouseX       float64
	mouseY       float64
	mouseLocked  bool
	mounted      bool
}

// =============================================================================
// 게임 컴포넌트들
// =============================================================================

// GameCanvas3D renders the 3D canvas area
type GameCanvas3D struct {
	app.Compo
	State        *GameState
}

// GameHUD3D renders the 3D game's heads-up display
type GameHUD3D struct {
	app.Compo
	State       *GameState
	OnPause     func()
	OnRestart   func()
	OnToggleDebug func()
}

// GameStats3D renders 3D game statistics
type GameStats3D struct {
	app.Compo
	State *GameState
}

// GameControls3D renders control instructions
type GameControls3D struct {
	app.Compo
}

// =============================================================================
// QubePage 구현
// =============================================================================

func (p *QubePage) OnMount(ctx app.Context) {
	p.mounted = true
	
	if p.game == nil {
		p.game = NewQubeGame()
		p.game.Start()
	}

	p.keys = make(map[string]bool)

	// Start game update loop with proper lifecycle
	go p.updateLoop(ctx)

	// Initialize Three.js after component mounts - delayed to ensure DOM is ready
	ctx.Async(func() {
		time.Sleep(100 * time.Millisecond) // Wait for DOM
		if p.mounted {
			p.initThreeJS()
		}
	})
}

func (p *QubePage) OnDismount() {
	p.mounted = false
	
	if p.game != nil {
		p.game.Stop()
	}
	
	// Cleanup Three.js if available
	window := app.Window()
	if window.Truthy() && window.Get("disposeThreeJS").Truthy() {
		window.Call("disposeThreeJS")
	}
}

func (p *QubePage) updateLoop(ctx app.Context) {
	ticker := time.NewTicker(16 * time.Millisecond) // 60 FPS
	defer ticker.Stop()

	for p.mounted {
		if p.game != nil && p.game.IsRunning() {
			// Process input
			p.processInput()
			
			// Update Three.js scene with current game state
			p.updateThreeJSGameState()
			
			// Update UI using ctx.Dispatch properly
			ctx.Dispatch(func(ctx app.Context) {
				ctx.Update()
			})
		}
		
		select {
		case <-ticker.C:
			// Continue to next iteration
		case <-time.After(16 * time.Millisecond):
			// Fallback timing
		}
		
		// Exit condition for goroutine cleanup
		if !p.mounted {
			break
		}
	}
}

func (p *QubePage) updateThreeJSGameState() {
	if p.game == nil {
		return
	}

	state := p.game.GetGameState()
	if state == nil {
		return
	}

	// Safely send game state to Three.js bridge using JSON
	window := app.Window()
	if window.Truthy() && window.Get("updateGameState").Truthy() {
		// Create a simplified state structure for JavaScript
		gameStateData := struct {
			Players  []*Player3D     `json:"players"`
			Objects  []*GameObject3D `json:"objects"`
			Camera   *Camera3D       `json:"camera"`
			Lighting *LightingState  `json:"lighting"`
			Time     float64         `json:"time"`
		}{
			Players:  state.Players,
			Objects:  state.Objects,
			Camera:   state.Camera,
			Lighting: state.Lighting,
			Time:     state.Time,
		}

		// Convert to JSON string
		jsonData, err := json.Marshal(gameStateData)
		if err != nil {
			return // Silently fail to avoid spam
		}

		// Send as JSON string to avoid syscall/js issues
		window.Call("updateGameState", string(jsonData))
	}
}

func (p *QubePage) initThreeJS() {
	// Check if window and functions are available before calling
	window := app.Window()
	if !window.Truthy() {
		return
	}
	
	// Check if Three.js is loaded
	if !window.Get("THREE").Truthy() {
		fmt.Println("Three.js not loaded yet, retrying...")
		time.AfterFunc(100*time.Millisecond, p.initThreeJS)
		return
	}
	
	// Check if our bridge function exists
	if !window.Get("initThreeJSBridge").Truthy() {
		fmt.Println("ThreeJS bridge not loaded yet, retrying...")  
		time.AfterFunc(100*time.Millisecond, p.initThreeJS)
		return
	}
	
	// Initialize Three.js through JavaScript bridge
	result := window.Call("initThreeJSBridge", "qube-canvas")
	if result.Bool() {
		fmt.Println("Three.js initialized successfully")
	} else {
		fmt.Println("Failed to initialize Three.js")
	}
}

func (p *QubePage) processInput() {
	if p.game == nil || len(p.game.State.Players) == 0 {
		return
	}

	var direction Vector3
	var moving bool

	// WASD movement
	if p.keys["KeyW"] || p.keys["ArrowUp"] {
		direction.Z -= 1
		moving = true
	}
	if p.keys["KeyS"] || p.keys["ArrowDown"] {
		direction.Z += 1
		moving = true
	}
	if p.keys["KeyA"] || p.keys["ArrowLeft"] {
		direction.X -= 1
		moving = true
	}
	if p.keys["KeyD"] || p.keys["ArrowRight"] {
		direction.X += 1
		moving = true
	}

	if moving {
		p.game.SendEvent(GameEvent{
			Type: EventPlayerMove,
			Data: PlayerMoveData{
				PlayerID:  "player1",
				Direction: direction,
				Speed:     5.0,
			},
		})
	}

	// Space for jump
	if p.keys["Space"] {
		p.game.SendEvent(GameEvent{Type: EventPlayerJump})
		p.keys["Space"] = false // Prevent continuous jumping
	}
}

func (p *QubePage) handleKeyDown(ctx app.Context, e app.Event) {
	key := e.Get("code").String()
	p.keys[key] = true

	// Prevent default for game keys
	gameKeys := []string{"KeyW", "KeyA", "KeyS", "KeyD", "Space", "ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"}
	for _, gameKey := range gameKeys {
		if key == gameKey {
			e.PreventDefault()
			break
		}
	}
}

func (p *QubePage) handleKeyUp(ctx app.Context, e app.Event) {
	key := e.Get("code").String()
	p.keys[key] = false
}

func (p *QubePage) handleMouseMove(ctx app.Context, e app.Event) {
	if !p.mouseLocked || p.game == nil {
		return
	}

	deltaX := e.Get("movementX").Float()
	deltaY := e.Get("movementY").Float()

	// Update camera rotation
	p.mouseX += deltaX * 0.002
	p.mouseY += deltaY * 0.002

	// Clamp vertical rotation
	if p.mouseY > 1.5 {
		p.mouseY = 1.5
	}
	if p.mouseY < -1.5 {
		p.mouseY = -1.5
	}

	// Apply camera rotation to game state
	if len(p.game.State.Players) > 0 {
		player := p.game.State.Players[0]
		
		// Calculate camera position based on mouse rotation
		distance := 10.0
		height := 5.0
		
		// Calculate new camera position around player
		x := player.Position.X + distance * math.Cos(p.mouseX) * math.Cos(p.mouseY)
		y := player.Position.Y + height + distance * math.Sin(p.mouseY)
		z := player.Position.Z + distance * math.Sin(p.mouseX) * math.Cos(p.mouseY)
		
		// Send camera move event
		p.game.SendEvent(GameEvent{
			Type: EventCameraMove,
			Data: CameraMoveData{
				Position: Vector3{X: x, Y: y, Z: z},
				Target:   player.Position,
			},
		})
	}
}

func (p *QubePage) handleCanvasClick(ctx app.Context, e app.Event) {
	// Request pointer lock for FPS-style camera control
	window := app.Window()
	if !window.Truthy() {
		return
	}
	
	document := window.Get("document")
	if !document.Truthy() {
		return
	}
	
	canvas := document.Call("getElementById", "qube-canvas")
	if !canvas.Truthy() {
		return
	}
	
	// Only request pointer lock if not already locked
	pointerLockElement := document.Get("pointerLockElement")
	if !pointerLockElement.Truthy() || pointerLockElement != canvas {
		// Use a timeout to handle SecurityError
		ctx.Async(func() {
			if window.Truthy() && canvas.Truthy() && canvas.Get("requestPointerLock").Truthy() {
				canvas.Call("requestPointerLock")
			}
		})
	}
	
	// Set up pointer lock change listener
	p.setupPointerLockListener()
}

func (p *QubePage) setupPointerLockListener() {
	window := app.Window()
	if !window.Truthy() {
		return
	}
	
	document := window.Get("document")
	if !document.Truthy() {
		return
	}
	
	// Listen for pointer lock changes
	document.Call("addEventListener", "pointerlockchange", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		pointerLockElement := document.Get("pointerLockElement")
		canvas := document.Call("getElementById", "qube-canvas")
		
		if pointerLockElement.Truthy() && pointerLockElement.Equal(canvas) {
			p.mouseLocked = true
		} else {
			p.mouseLocked = false
		}
		return nil
	}))
}

func (p *QubePage) Render() app.UI {
	if p.game == nil {
		return app.Div().Text("Loading 3D Engine...")
	}

	state := p.game.GetGameState()

	return app.Div().Class("qube-game-container").
		TabIndex(0). // Allow focus for keyboard events
		OnKeyDown(p.handleKeyDown).
		OnKeyUp(p.handleKeyUp).
		OnMouseMove(p.handleMouseMove). // Add mouse move handler
		Body(
			app.H1().Class("game-title").Text("🎮 Qube 3D Game"),

			app.Div().Class("game-layout-3d").Body(
				// 3D Canvas Area
				app.Div().Class("canvas-container").Body(
					&GameCanvas3D{
						State: state,
					},
				),

				// Side Panel
				app.Div().Class("side-panel-3d").Body(
					&GameStats3D{
						State: state,
					},

					&GameHUD3D{
						State:         state,
						OnPause:       p.handlePause,
						OnRestart:     p.handleRestart,
						OnToggleDebug: p.handleToggleDebug,
					},

					&GameControls3D{},
				),
			),
		)
}

func (p *QubePage) handlePause() {
	if p.game.State.Paused {
		p.game.SendEvent(GameEvent{Type: EventResumeGame})
	} else {
		p.game.SendEvent(GameEvent{Type: EventPauseGame})
	}
}

func (p *QubePage) handleRestart() {
	p.game.SendEvent(GameEvent{Type: EventRestartGame})
}

func (p *QubePage) handleToggleDebug() {
	p.game.Config.ShowDebug = !p.game.Config.ShowDebug
}

// =============================================================================
// GameCanvas3D 구현
// =============================================================================

func (c *GameCanvas3D) OnMount(ctx app.Context) {
	// Initial setup - Three.js initialization happens in parent component
}

func (c *GameCanvas3D) updateThreeJSScene() {
	if c.State == nil {
		return
	}

	// Safely send game state to Three.js bridge using JSON
	window := app.Window()
	if window.Truthy() && window.Get("updateGameState").Truthy() {
		// Create a simplified state structure for JavaScript
		gameStateData := struct {
			Players  []*Player3D     `json:"players"`
			Objects  []*GameObject3D `json:"objects"`
			Camera   *Camera3D       `json:"camera"`
			Lighting *LightingState  `json:"lighting"`
			Time     float64         `json:"time"`
		}{
			Players:  c.State.Players,
			Objects:  c.State.Objects,
			Camera:   c.State.Camera,
			Lighting: c.State.Lighting,
			Time:     c.State.Time,
		}

		// Convert to JSON string
		jsonData, err := json.Marshal(gameStateData)
		if err != nil {
			fmt.Printf("Error marshaling game state: %v\n", err)
			return
		}

		// Send as JSON string to avoid syscall/js issues
		window.Call("updateGameState", string(jsonData))
	}
}

func (c *GameCanvas3D) Render() app.UI {
	// Don't call updateThreeJSScene on every render - it's handled in OnMount
	
	return app.Div().Class("canvas-wrapper").Body(
		app.Canvas().
			ID("qube-canvas").
			Class("game-canvas-3d").
			OnClick(func(ctx app.Context, e app.Event) {
				// Safely request pointer lock
				window := app.Window()
				if window.Truthy() {
					document := window.Get("document")
					if document.Truthy() {
						canvas := document.Call("getElementById", "qube-canvas")
						if canvas.Truthy() && canvas.Get("requestPointerLock").Truthy() {
							canvas.Call("requestPointerLock")
						}
					}
				}
			}),
		
		// Loading message
		app.Div().Class("canvas-loading").
			Style("display", "none").
			Text("Loading 3D Scene..."),
	)
}

// =============================================================================
// GameStats3D 구현
// =============================================================================

func (s *GameStats3D) Render() app.UI {
	activeObjects := 0
	for _, obj := range s.State.Objects {
		if obj.IsActive && obj.ID != "ground" {
			activeObjects++
		}
	}

	var playerPos Vector3
	if len(s.State.Players) > 0 {
		playerPos = s.State.Players[0].Position
	}

	return app.Div().Class("game-stats-3d").Body(
		app.H3().Text("📊 게임 정보"),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("🏆 점수:"),
			app.Span().Class("stat-value").Text(strconv.Itoa(s.State.Score)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("⏰ 시간:"),
			app.Span().Class("stat-value").Text(fmt.Sprintf("%.1f초", s.State.Time)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("🎯 남은 큐브:"),
			app.Span().Class("stat-value").Text(strconv.Itoa(activeObjects)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("📍 위치:"),
			app.Span().Class("stat-value").Text(fmt.Sprintf("(%.1f, %.1f, %.1f)", 
				playerPos.X, playerPos.Y, playerPos.Z)),
		),

		app.If(s.State.Paused,
			func() app.UI {
				return app.Div().Class("pause-indicator").Text("⏸️ 일시정지")
			},
		),
	)
}

// =============================================================================
// GameHUD3D 구현
// =============================================================================

func (h *GameHUD3D) Render() app.UI {
	return app.Div().Class("game-hud-3d").Body(
		app.H3().Text("🎮 컨트롤"),

		app.Div().Class("control-buttons-3d").Body(
			// 일시정지 버튼
			app.Button().
				Class("control-button pause-button").
				OnClick(func(ctx app.Context, e app.Event) {
					if h.OnPause != nil {
						h.OnPause()
					}
				}).
				Body(
					app.If(h.State.Paused,
						func() app.UI { return app.Text("▶️ 재개") },
					).Else(
						func() app.UI { return app.Text("⏸️ 일시정지") },
					),
				),

			// 재시작 버튼
			app.Button().
				Class("control-button restart-button").
				OnClick(func(ctx app.Context, e app.Event) {
					if h.OnRestart != nil {
						h.OnRestart()
					}
				}).
				Text("🔄 재시작"),

			// 디버그 토글 버튼
			app.Button().
				Class("control-button debug-button").
				OnClick(func(ctx app.Context, e app.Event) {
					if h.OnToggleDebug != nil {
						h.OnToggleDebug()
					}
				}).
				Text("🔍 디버그 모드"),
		),

		// 게임 오버 메시지
		app.If(h.State.GameOver,
			func() app.UI {
				return app.Div().Class("game-over-3d").Body(
					app.H2().Text("🎮 게임 종료!"),
					app.P().Text(fmt.Sprintf("최종 점수: %d", h.State.Score)),
				)
			},
		),
	)
}

// =============================================================================
// GameControls3D 구현
// =============================================================================

func (c *GameControls3D) Render() app.UI {
	return app.Div().Class("game-controls-3d").Body(
		app.H4().Text("🎮 조작법"),
		
		app.Div().Class("control-item").Text("WASD / 방향키: 이동"),
		app.Div().Class("control-item").Text("Space: 점프"),
		app.Div().Class("control-item").Text("마우스: 카메라 회전"),
		app.Div().Class("control-item").Text("클릭: 포인터 잠금"),
		app.Div().Class("control-item").Text("ESC: 포인터 잠금 해제"),
		
		app.Div().Class("help-separator"),
		
		app.H4().Text("🎯 목표"),
		app.Div().Class("control-item").Text("💎 모든 컬러 큐브를 수집하세요!"),
		app.Div().Class("control-item").Text("🏃‍♂️ 새로운 큐브가 계속 생성됩니다"),
	)
}