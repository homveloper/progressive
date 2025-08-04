package tower

import (
	"math"
	"time"
)

// =============================================================================
// 게임 기본 타입들
// =============================================================================

// Position represents a coordinate on the game map
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Size represents dimensions
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// =============================================================================
// 타워 시스템
// =============================================================================

// TowerType represents different types of towers
type TowerType int

const (
	TowerArcher TowerType = iota // 궁수 타워
	TowerMage                    // 마법사 타워
	TowerCannon                  // 대포 타워
)

// Tower represents a defensive tower
type Tower struct {
	ID           int       `json:"id"`
	Type         TowerType `json:"type"`
	Position     Position  `json:"position"`
	Level        int       `json:"level"`
	Damage       float64   `json:"damage"`
	Range        float64   `json:"range"`
	FireRate     float64   `json:"fireRate"` // shots per second
	Cost         int       `json:"cost"`
	LastShot     time.Time `json:"lastShot"`
	Target       *Enemy    `json:"target,omitempty"`
	PendingShots []PendingShot `json:"pendingShots,omitempty"` // 예정된 데미지
}

// PendingShot represents a delayed damage application
type PendingShot struct {
	Target    *Enemy    `json:"target"`
	Damage    float64   `json:"damage"`
	HitTime   time.Time `json:"hitTime"`   // 데미지가 적용될 시간
	ShotID    int       `json:"shotID"`    // 발사 ID
}

// TowerStats defines base stats for each tower type
var TowerStats = map[TowerType]struct {
	Name     string
	Damage   float64
	Range    float64
	FireRate float64
	Cost     int
	Color    string
}{
	TowerArcher: {"궁수 타워", 25.0, 80.0, 1.5, 50, "#8B4513"},
	TowerMage:   {"마법사 타워", 40.0, 60.0, 1.0, 75, "#4B0082"},
	TowerCannon: {"대포 타워", 80.0, 100.0, 0.5, 100, "#2F4F4F"},
}

// =============================================================================
// 적 시스템
// =============================================================================

// EnemyType represents different types of enemies
type EnemyType int

const (
	EnemyGoblin EnemyType = iota // 고블린
	EnemyOrc                     // 오크
	EnemyTroll                   // 트롤
)

// Enemy represents a hostile unit that moves in a square loop
type Enemy struct {
	ID          int       `json:"id"`
	Type        EnemyType `json:"type"`
	Position    Position  `json:"position"`
	Health      float64   `json:"health"`
	MaxHealth   float64   `json:"maxHealth"`
	Speed       float64   `json:"speed"`
	Reward      int       `json:"reward"`
	PathProgress float64  `json:"pathProgress"` // 0.0-4.0 representing progress around square
	Corner      int       `json:"corner"`      // 0-3 representing current corner
	IsAlive     bool      `json:"isAlive"`
	LapCount    int       `json:"lapCount"`    // Number of completed laps
}

// EnemyStats defines base stats for each enemy type
var EnemyStats = map[EnemyType]struct {
	Name      string
	Health    float64
	Speed     float64 // pixels per second along the square path
	Reward    int
	Color     string
}{
	EnemyGoblin: {"고블린", 50.0, 60.0, 10, "#228B22"},  // 빠름
	EnemyOrc:    {"오크", 100.0, 40.0, 20, "#8B0000"},   // 보통
	EnemyTroll:  {"트롤", 200.0, 25.0, 40, "#2F2F2F"},   // 느림
}

// =============================================================================
// 발사체 시스템
// =============================================================================

// Projectile represents a visual projectile effect (no actual damage logic)
type Projectile struct {
	ID         int      `json:"id"`
	Position   Position `json:"position"`
	Target     *Enemy   `json:"target"`
	Speed      float64  `json:"speed"`
	Color      string   `json:"color"`
	IsActive   bool     `json:"isActive"`
	StartTime  time.Time `json:"startTime"`  // 발사 시작 시간
	LifeTime   float64   `json:"lifeTime"`   // 생존 시간 (초)
}

// =============================================================================
// 맵 시스템
// =============================================================================

// GameMap represents the game map
type GameMap struct {
	Size Size       `json:"size"`
	Path []Position `json:"path"`
}

// =============================================================================
// 웨이브 시스템
// =============================================================================

// Wave represents a wave of enemies
type Wave struct {
	Number    int         `json:"number"`
	Enemies   []EnemySpawn `json:"enemies"`
	Completed bool        `json:"completed"`
}

// EnemySpawn defines when and what type of enemy to spawn
type EnemySpawn struct {
	Type  EnemyType `json:"type"`
	Count int       `json:"count"`
	Delay float64   `json:"delay"` // seconds between spawns
}

// =============================================================================
// 게임 상태
// =============================================================================

// GameState represents the complete game state
type GameState struct {
	// 맵 정보
	Map GameMap `json:"map"`
	
	// 게임 오브젝트들
	Towers      []*Tower      `json:"towers"`
	Enemies     []*Enemy      `json:"enemies"`
	Projectiles []*Projectile `json:"projectiles"`
	
	// 게임 상태
	CurrentWave    int       `json:"currentWave"`
	NextWaveTime   time.Time `json:"nextWaveTime"`   // 다음 웨이브 시작 시간
	MaxEnemies     int       `json:"maxEnemies"`     // 현재 웨이브 최대 적 수
	EnemiesSpawned int       `json:"enemiesSpawned"` // 현재 웨이브에서 스폰된 적 수
	EnemiesKilled  int       `json:"enemiesKilled"`  // 현재 웨이브에서 죽인 적 수
	
	// 플레이어 상태
	Gold  int `json:"gold"`
	Score int `json:"score"`
	
	// 게임 메타
	GameOver bool `json:"gameOver"`
	Paused   bool `json:"paused"`
	
	// 내부 카운터들
	NextTowerID      int `json:"nextTowerID"`
	NextEnemyID      int `json:"nextEnemyID"`
	NextProjectileID int `json:"nextProjectileID"`
	
	// 타이밍
	LastUpdate       time.Time `json:"lastUpdate"`
	NextSpawnTime    time.Time `json:"nextSpawnTime"`
	WaveTimeRemaining float64  `json:"waveTimeRemaining"` // 다음 웨이브까지 남은 시간(초)
}

// =============================================================================
// 생성자 함수들
// =============================================================================

// NewGameState creates a new game state
func NewGameState() *GameState {
	now := time.Now()
	return &GameState{
		Map:         CreateDefaultMap(),
		Towers:      make([]*Tower, 0),
		Enemies:     make([]*Enemy, 0),
		Projectiles: make([]*Projectile, 0),
		
		CurrentWave:      1,
		NextWaveTime:     now.Add(60 * time.Second), // 1분 후 첫 웨이브
		MaxEnemies:       GetMaxEnemiesForWave(1),
		EnemiesSpawned:   0,
		EnemiesKilled:    0,
		
		Gold:  100,
		Score: 0,
		
		GameOver: false,
		Paused:   false,
		
		NextTowerID:      1,
		NextEnemyID:      1,
		NextProjectileID: 1,
		
		LastUpdate:       now,
		NextSpawnTime:    now,
		WaveTimeRemaining: 60.0,
	}
}

// CreateDefaultMap creates the default square loop map
func CreateDefaultMap() GameMap {
	// 정사각형 경로 정의 (시계방향)
	path := []Position{
		{X: 100, Y: 100}, // 좌상단 꼭짓점 (스폰 지점)
		{X: 500, Y: 100}, // 우상단 꼭짓점
		{X: 500, Y: 300}, // 우하단 꼭짓점
		{X: 100, Y: 300}, // 좌하단 꼭짓점
		{X: 100, Y: 100}, // 다시 시작점으로 (무한 루프)
	}
	
	return GameMap{
		Size: Size{Width: 600, Height: 400},
		Path: path,
	}
}

// NewTower creates a new tower
func NewTower(id int, towerType TowerType, pos Position) *Tower {
	stats := TowerStats[towerType]
	
	return &Tower{
		ID:           id,
		Type:         towerType,
		Position:     pos,
		Level:        1,
		Damage:       stats.Damage,
		Range:        stats.Range,
		FireRate:     stats.FireRate,
		Cost:         stats.Cost,
		LastShot:     time.Now(),
		PendingShots: make([]PendingShot, 0),
	}
}

// NewEnemy creates a new enemy at the spawn corner
func NewEnemy(id int, enemyType EnemyType) *Enemy {
	stats := EnemyStats[enemyType]
	
	return &Enemy{
		ID:          id,
		Type:        enemyType,
		Position:    Position{X: 100, Y: 100}, // 좌상단 꼭짓점에서 시작
		Health:      stats.Health,
		MaxHealth:   stats.Health,
		Speed:       stats.Speed,
		Reward:      stats.Reward,
		PathProgress: 0.0,
		Corner:      0,
		IsAlive:     true,
		LapCount:    0,
	}
}

// NewProjectile creates a new visual projectile (no damage logic)
func NewProjectile(id int, pos Position, target *Enemy, color string, lifeTime float64) *Projectile {
	return &Projectile{
		ID:        id,
		Position:  pos,
		Target:    target,
		Speed:     200.0, // pixels per second
		Color:     color,
		IsActive:  true,
		StartTime: time.Now(),
		LifeTime:  lifeTime, // 시각적 효과 지속 시간
	}
}

// =============================================================================
// 유틸리티 함수들
// =============================================================================

// Distance calculates distance between two positions
func Distance(p1, p2 Position) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// GetMaxEnemiesForWave returns the maximum number of enemies allowed for a wave
func GetMaxEnemiesForWave(waveNumber int) int {
	// 웨이브가 올라갈수록 더 많은 적이 허용됨
	baseMax := 10
	increment := 3
	return baseMax + (waveNumber-1)*increment
}

// GetWaveDefinition returns the enemy composition for a given wave
func GetWaveDefinition(waveNumber int) Wave {
	// 웨이브별 스폰할 적 구성 (더 이상 한번에 다 스폰하지 않음)
	baseGoblins := 3 + waveNumber
	baseOrcs := int(math.Max(0, float64(waveNumber-2)))
	baseTrolls := int(math.Max(0, float64(waveNumber-5)))
	
	enemies := []EnemySpawn{}
	
	if baseGoblins > 0 {
		enemies = append(enemies, EnemySpawn{
			Type:  EnemyGoblin,
			Count: baseGoblins,
			Delay: 2.0, // 2초마다 스폰
		})
	}
	
	if baseOrcs > 0 {
		enemies = append(enemies, EnemySpawn{
			Type:  EnemyOrc,
			Count: baseOrcs,
			Delay: 3.0, // 3초마다 스폰
		})
	}
	
	if baseTrolls > 0 {
		enemies = append(enemies, EnemySpawn{
			Type:  EnemyTroll,
			Count: baseTrolls,
			Delay: 4.0, // 4초마다 스폰
		})
	}
	
	return Wave{
		Number:    waveNumber,
		Enemies:   enemies,
		Completed: false,
	}
}