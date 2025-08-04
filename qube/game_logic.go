package qube

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// 메인 게임 클래스
// =============================================================================

// QubeGame represents the main 3D Qube game
type QubeGame struct {
	State     *GameState    `json:"state"`
	Config    *GameConfig   `json:"config"`
	running   bool
	eventChan chan GameEvent
	mutex     sync.RWMutex
	renderer  *ThreeJSRenderer
}

// ThreeJSRenderer handles Three.js bridge communication
type ThreeJSRenderer struct {
	initialized bool
	canvasID    string
}

// NewQubeGame creates a new 3D Qube game instance
func NewQubeGame() *QubeGame {
	game := &QubeGame{
		State:     NewGameState(),
		Config:    DefaultGameConfig(),
		running:   false,
		eventChan: make(chan GameEvent, 100),
		renderer:  &ThreeJSRenderer{canvasID: "qube-canvas"},
	}

	return game
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	state := &GameState{
		Players:    make([]*Player3D, 0),
		Objects:    make([]*GameObject3D, 0),
		Score:      0,
		Time:       0,
		Paused:     false,
		GameOver:   false,
		LastUpdate: time.Now(),
	}

	// Initialize camera
	state.Camera = &Camera3D{
		Position: Vector3{X: 0, Y: 5, Z: 10},
		Target:   Vector3{X: 0, Y: 0, Z: 0},
		Up:       Vector3{X: 0, Y: 1, Z: 0},
		FOV:      75,
		Near:     0.1,
		Far:      1000,
	}

	// Initialize lighting
	state.Lighting = &LightingState{
		AmbientColor:     "#404040",
		AmbientIntensity: 0.4,
	}
	state.Lighting.DirectionalLight.Direction = Vector3{X: 1, Y: 1, Z: 1}.Normalize()
	state.Lighting.DirectionalLight.Color = "#ffffff"
	state.Lighting.DirectionalLight.Intensity = 0.8

	// Initialize physics world
	state.Physics = &PhysicsWorld{
		Bodies:     make([]*RigidBody, 0),
		Gravity:    Vector3{X: 0, Y: -9.8, Z: 0},
		TimeStep:   1.0 / 60.0,
		Collisions: make([]*CollisionInfo, 0),
	}

	return state
}

// =============================================================================
// 게임 라이프사이클
// =============================================================================

// Start starts the game
func (g *QubeGame) Start() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.running {
		return
	}

	g.running = true
	g.initializeGame()

	// Start game loop
	go g.gameLoop()
	go g.eventLoop()
}

// Stop stops the game
func (g *QubeGame) Stop() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.running = false
}

// initializeGame initializes the game objects
func (g *QubeGame) initializeGame() {
	// Create player
	player := &Player3D{
		ID:        "player1",
		Position:  Vector3{X: 0, Y: 1, Z: 0},
		Rotation:  Vector3{X: 0, Y: 0, Z: 0},
		Velocity:  Vector3{X: 0, Y: 0, Z: 0},
		Health:    100,
		MaxHealth: 100,
		Speed:     5.0,
		IsActive:  true,
	}
	g.State.Players = append(g.State.Players, player)

	// Create ground plane
	ground := &GameObject3D{
		ID:       "ground",
		Position: Vector3{X: 0, Y: -0.5, Z: 0},
		Rotation: Vector3{X: 0, Y: 0, Z: 0},
		Scale:    Vector3{X: 20, Y: 1, Z: 20},
		Model:    "box",
		Color:    "#808080",
		IsActive: true,
	}
	g.State.Objects = append(g.State.Objects, ground)

	// Create some collectible cubes
	for i := 0; i < 10; i++ {
		cube := &GameObject3D{
			ID:       fmt.Sprintf("cube_%d", i),
			Position: Vector3{
				X: rand.Float64()*16 - 8,  // -8 to 8
				Y: rand.Float64()*3 + 1,   // 1 to 4
				Z: rand.Float64()*16 - 8,  // -8 to 8
			},
			Rotation: Vector3{
				X: rand.Float64() * 2 * math.Pi,
				Y: rand.Float64() * 2 * math.Pi,
				Z: rand.Float64() * 2 * math.Pi,
			},
			Scale:    Vector3{X: 0.5, Y: 0.5, Z: 0.5},
			Model:    "box",
			Color:    g.getRandomColor(),
			IsActive: true,
		}
		g.State.Objects = append(g.State.Objects, cube)
	}

	// Create physics bodies
	g.initializePhysics()
}

// initializePhysics sets up physics bodies
func (g *QubeGame) initializePhysics() {
	// Player physics body
	if len(g.State.Players) > 0 {
		player := g.State.Players[0]
		playerBody := &RigidBody{
			ID:       player.ID,
			Position: player.Position,
			Velocity: player.Velocity,
			Mass:     1.0,
			Gravity:  true,
			Shape:    "box",
			Size:     Vector3{X: 1, Y: 2, Z: 1},
		}
		g.State.Physics.Bodies = append(g.State.Physics.Bodies, playerBody)
	}

	// Ground physics body
	groundBody := &RigidBody{
		ID:       "ground",
		Position: Vector3{X: 0, Y: -0.5, Z: 0},
		Velocity: Vector3{X: 0, Y: 0, Z: 0},
		Mass:     0, // Static body
		Gravity:  false,
		Shape:    "box",
		Size:     Vector3{X: 20, Y: 1, Z: 20},
	}
	g.State.Physics.Bodies = append(g.State.Physics.Bodies, groundBody)
}

// =============================================================================
// 게임 루프
// =============================================================================

// gameLoop runs the main game update loop
func (g *QubeGame) gameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // 60 FPS
	defer ticker.Stop()

	for g.running {
		select {
		case <-ticker.C:
			if !g.State.Paused && !g.State.GameOver {
				g.update()
			}
		}
	}
}

// eventLoop processes game events
func (g *QubeGame) eventLoop() {
	for g.running {
		select {
		case event := <-g.eventChan:
			g.handleEvent(event)
		}
	}
}

// update updates the game state
func (g *QubeGame) update() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now()
	deltaTime := now.Sub(g.State.LastUpdate).Seconds()
	g.State.LastUpdate = now
	g.State.Time += deltaTime

	// Update physics
	g.updatePhysics(deltaTime)

	// Update game objects
	g.updateObjects(deltaTime)

	// Update camera
	g.updateCamera(deltaTime)

	// Check collisions
	g.checkCollisions()
}

// updatePhysics updates the physics simulation
func (g *QubeGame) updatePhysics(deltaTime float64) {
	physics := g.State.Physics

	for _, body := range physics.Bodies {
		if body.Mass == 0 { // Static body
			continue
		}

		// Apply gravity
		if body.Gravity {
			body.Velocity = body.Velocity.Add(physics.Gravity.Scale(deltaTime))
		}

		// Update position
		body.Position = body.Position.Add(body.Velocity.Scale(deltaTime))

		// Ground collision (simple)
		if body.Position.Y < 0 {
			body.Position.Y = 0
			if body.Velocity.Y < 0 {
				body.Velocity.Y = 0
			}
		}

		// Update corresponding game object
		g.updateObjectFromPhysics(body)
	}
}

// updateObjectFromPhysics updates game objects from physics bodies
func (g *QubeGame) updateObjectFromPhysics(body *RigidBody) {
	// Update player
	for _, player := range g.State.Players {
		if player.ID == body.ID {
			player.Position = body.Position
			player.Velocity = body.Velocity
			return
		}
	}

	// Update objects
	for _, obj := range g.State.Objects {
		if obj.ID == body.ID {
			obj.Position = body.Position
			return
		}
	}
}

// updateObjects updates game objects
func (g *QubeGame) updateObjects(deltaTime float64) {
	// Rotate collectible cubes
	for _, obj := range g.State.Objects {
		if obj.ID != "ground" && obj.ID != "player1" {
			obj.Rotation.Y += deltaTime * 2.0 // 2 radians per second
			obj.Rotation.X += deltaTime * 1.5
		}
	}
}

// updateCamera updates camera position
func (g *QubeGame) updateCamera(deltaTime float64) {
	if g.Config.CameraFollow && len(g.State.Players) > 0 {
		player := g.State.Players[0]
		targetPos := player.Position.Add(Vector3{X: 0, Y: 5, Z: 10})
		
		// Smooth camera following
		lerpFactor := deltaTime * 2.0
		g.State.Camera.Position = Lerp(g.State.Camera.Position, targetPos, lerpFactor)
		g.State.Camera.Target = Lerp(g.State.Camera.Target, player.Position, lerpFactor)
	}
}

// checkCollisions checks for collisions
func (g *QubeGame) checkCollisions() {
	if len(g.State.Players) == 0 {
		return
	}

	player := g.State.Players[0]
	
	// Check collision with collectible objects
	for i, obj := range g.State.Objects {
		if !obj.IsActive || obj.ID == "ground" {
			continue
		}

		distance := Distance(player.Position, obj.Position)
		if distance < 1.5 { // Collision threshold
			// Collect object
			obj.IsActive = false
			g.State.Score += 10
			
			// Remove from objects list
			g.State.Objects = append(g.State.Objects[:i], g.State.Objects[i+1:]...)
			break
		}
	}

	// Check if all objects collected
	activeObjects := 0
	for _, obj := range g.State.Objects {
		if obj.IsActive && obj.ID != "ground" {
			activeObjects++
		}
	}

	if activeObjects == 0 {
		g.spawnMoreObjects()
	}
}

// spawnMoreObjects spawns new collectible objects
func (g *QubeGame) spawnMoreObjects() {
	for i := 0; i < 5; i++ {
		cube := &GameObject3D{
			ID:       fmt.Sprintf("cube_%d_%d", int(g.State.Time), i),
			Position: Vector3{
				X: rand.Float64()*16 - 8,
				Y: rand.Float64()*3 + 1,
				Z: rand.Float64()*16 - 8,
			},
			Rotation: Vector3{
				X: rand.Float64() * 2 * math.Pi,
				Y: rand.Float64() * 2 * math.Pi,
				Z: rand.Float64() * 2 * math.Pi,
			},
			Scale:    Vector3{X: 0.5, Y: 0.5, Z: 0.5},
			Model:    "box",
			Color:    g.getRandomColor(),
			IsActive: true,
		}
		g.State.Objects = append(g.State.Objects, cube)
	}
}

// =============================================================================
// 이벤트 처리
// =============================================================================

// SendEvent sends a game event
func (g *QubeGame) SendEvent(event GameEvent) {
	select {
	case g.eventChan <- event:
	default:
		// Event queue full, skip
	}
}

// handleEvent processes a game event
func (g *QubeGame) handleEvent(event GameEvent) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	switch event.Type {
	case EventPlayerMove:
		g.handlePlayerMove(event.Data.(PlayerMoveData))
	case EventPlayerJump:
		g.handlePlayerJump()
	case EventPauseGame:
		g.State.Paused = true
	case EventResumeGame:
		g.State.Paused = false
	case EventRestartGame:
		g.restartGame()
	case EventCameraMove:
		g.handleCameraMove(event.Data.(CameraMoveData))
	}
}

// handlePlayerMove handles player movement
func (g *QubeGame) handlePlayerMove(data PlayerMoveData) {
	if len(g.State.Players) == 0 {
		return
	}

	player := g.State.Players[0]
	
	// Find player physics body
	for _, body := range g.State.Physics.Bodies {
		if body.ID == player.ID {
			// Apply movement force
			force := data.Direction.Normalize().Scale(data.Speed)
			body.Velocity.X = force.X
			body.Velocity.Z = force.Z
			break
		}
	}
}

// handlePlayerJump handles player jumping
func (g *QubeGame) handlePlayerJump() {
	if len(g.State.Players) == 0 {
		return
	}

	player := g.State.Players[0]
	
	// Find player physics body
	for _, body := range g.State.Physics.Bodies {
		if body.ID == player.ID {
			// Only jump if on ground
			if body.Position.Y <= 0.1 {
				body.Velocity.Y = g.Config.JumpForce
			}
			break
		}
	}
}

// handleCameraMove handles camera movement
func (g *QubeGame) handleCameraMove(data CameraMoveData) {
	g.State.Camera.Position = data.Position
	g.State.Camera.Target = data.Target
}

// restartGame restarts the game
func (g *QubeGame) restartGame() {
	g.State = NewGameState()
	g.initializeGame()
}

// =============================================================================
// 게임 상태 접근자
// =============================================================================

// GetGameState returns a copy of the current game state
func (g *QubeGame) GetGameState() *GameState {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	
	// Return a copy to prevent concurrent access issues
	stateCopy := *g.State
	return &stateCopy
}

// IsRunning returns whether the game is running
func (g *QubeGame) IsRunning() bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.running
}

// =============================================================================
// 헬퍼 함수들
// =============================================================================

// getRandomColor returns a random hex color
func (g *QubeGame) getRandomColor() string {
	colors := []string{
		"#ff6b6b", "#4ecdc4", "#45b7d1", "#96ceb4", "#feca57",
		"#ff9ff3", "#54a0ff", "#5f27cd", "#00d2d3", "#ff9f43",
	}
	return colors[rand.Intn(len(colors))]
}