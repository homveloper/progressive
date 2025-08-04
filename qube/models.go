package qube

import (
	"math"
	"time"
)

// =============================================================================
// 3D 수학 구조체들
// =============================================================================

// Vector3 represents a 3D vector
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Add adds two vectors
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// Subtract subtracts two vectors
func (v Vector3) Subtract(other Vector3) Vector3 {
	return Vector3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// Scale scales a vector by a scalar
func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

// Length returns the length of the vector
func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a normalized vector
func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length == 0 {
		return Vector3{}
	}
	return Vector3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

// =============================================================================
// 게임 오브젝트들
// =============================================================================

// GameObject3D represents a 3D game object
type GameObject3D struct {
	ID       string  `json:"id"`
	Position Vector3 `json:"position"`
	Rotation Vector3 `json:"rotation"`
	Scale    Vector3 `json:"scale"`
	Model    string  `json:"model"`
	Texture  string  `json:"texture"`
	Color    string  `json:"color"`
	IsActive bool    `json:"isActive"`
}

// Player3D represents the player in 3D space
type Player3D struct {
	ID        string  `json:"id"`
	Position  Vector3 `json:"position"`
	Rotation  Vector3 `json:"rotation"`
	Velocity  Vector3 `json:"velocity"`
	Health    float64 `json:"health"`
	MaxHealth float64 `json:"maxHealth"`
	Speed     float64 `json:"speed"`
	IsActive  bool    `json:"isActive"`
}

// Camera3D represents the camera settings
type Camera3D struct {
	Position Vector3 `json:"position"`
	Target   Vector3 `json:"target"`
	Up       Vector3 `json:"up"`
	FOV      float64 `json:"fov"`
	Near     float64 `json:"near"`
	Far      float64 `json:"far"`
}

// LightingState represents the lighting configuration
type LightingState struct {
	AmbientColor     string  `json:"ambientColor"`
	AmbientIntensity float64 `json:"ambientIntensity"`
	DirectionalLight struct {
		Direction Vector3 `json:"direction"`
		Color     string  `json:"color"`
		Intensity float64 `json:"intensity"`
	} `json:"directionalLight"`
}

// =============================================================================
// 물리 시스템
// =============================================================================

// RigidBody represents a physics body
type RigidBody struct {
	ID       string  `json:"id"`
	Position Vector3 `json:"position"`
	Velocity Vector3 `json:"velocity"`
	Mass     float64 `json:"mass"`
	Gravity  bool    `json:"gravity"`
	Shape    string  `json:"shape"` // "box", "sphere", "plane"
	Size     Vector3 `json:"size"`
}

// CollisionInfo represents collision data
type CollisionInfo struct {
	ObjectA   string  `json:"objectA"`
	ObjectB   string  `json:"objectB"`
	Point     Vector3 `json:"point"`
	Normal    Vector3 `json:"normal"`
	Penetration float64 `json:"penetration"`
}

// PhysicsWorld manages physics simulation
type PhysicsWorld struct {
	Bodies      []*RigidBody      `json:"bodies"`
	Gravity     Vector3           `json:"gravity"`
	TimeStep    float64           `json:"timeStep"`
	Collisions  []*CollisionInfo  `json:"collisions"`
}

// =============================================================================
// 게임 상태
// =============================================================================

// GameState represents the complete game state
type GameState struct {
	Players     []*Player3D       `json:"players"`
	Objects     []*GameObject3D   `json:"objects"`
	Camera      *Camera3D         `json:"camera"`
	Lighting    *LightingState    `json:"lighting"`
	Physics     *PhysicsWorld     `json:"physics"`
	Score       int               `json:"score"`
	Time        float64           `json:"time"`
	Paused      bool              `json:"paused"`
	GameOver    bool              `json:"gameOver"`
	LastUpdate  time.Time         `json:"-"`
}

// =============================================================================
// 게임 이벤트
// =============================================================================

// GameEventType represents different types of game events
type GameEventType string

const (
	EventPlayerMove     GameEventType = "player_move"
	EventPlayerJump     GameEventType = "player_jump"
	EventObjectCollect  GameEventType = "object_collect"
	EventPauseGame      GameEventType = "pause_game"
	EventResumeGame     GameEventType = "resume_game"
	EventRestartGame    GameEventType = "restart_game"
	EventCameraMove     GameEventType = "camera_move"
)

// GameEvent represents a game event
type GameEvent struct {
	Type GameEventType `json:"type"`
	Data interface{}   `json:"data"`
}

// PlayerMoveData represents player movement data
type PlayerMoveData struct {
	PlayerID  string  `json:"playerId"`
	Direction Vector3 `json:"direction"`
	Speed     float64 `json:"speed"`
}

// CameraMoveData represents camera movement data
type CameraMoveData struct {
	Position Vector3 `json:"position"`
	Target   Vector3 `json:"target"`
}

// =============================================================================
// 게임 설정
// =============================================================================

// GameConfig holds game configuration
type GameConfig struct {
	WorldSize    Vector3 `json:"worldSize"`
	Gravity      Vector3 `json:"gravity"`
	PlayerSpeed  float64 `json:"playerSpeed"`
	JumpForce    float64 `json:"jumpForce"`
	CameraFollow bool    `json:"cameraFollow"`
	ShowDebug    bool    `json:"showDebug"`
}

// DefaultGameConfig returns default game configuration
func DefaultGameConfig() *GameConfig {
	return &GameConfig{
		WorldSize:    Vector3{X: 20, Y: 20, Z: 20},
		Gravity:      Vector3{X: 0, Y: -9.8, Z: 0},
		PlayerSpeed:  5.0,
		JumpForce:    10.0,
		CameraFollow: true,
		ShowDebug:    false,
	}
}

// =============================================================================
// 헬퍼 함수들
// =============================================================================

// Distance calculates the distance between two vectors
func Distance(a, b Vector3) float64 {
	return a.Subtract(b).Length()
}

// Lerp performs linear interpolation between two vectors
func Lerp(a, b Vector3, t float64) Vector3 {
	return a.Scale(1 - t).Add(b.Scale(t))
}

// Clamp clamps a value between min and max
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}