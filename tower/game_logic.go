package tower

import (
	"math"
	"time"
)

// =============================================================================
// 타워 디펜스 게임 엔진
// =============================================================================

// TowerDefenseGame manages the complete game logic
type TowerDefenseGame struct {
	State          *GameState
	eventChannel   chan GameEvent
	running        bool
	currentWave    Wave
}

// GameEvent represents different game events
type GameEvent struct {
	Type EventType
	Data interface{}
}

// EventType represents the type of game event
type EventType int

const (
	EventPlaceTower EventType = iota
	EventUpgradeTower
	EventSellTower
	EventStartWave
	EventForceNextWave // 개발자 모드: 다음 웨이브 강제 시작
	EventPauseGame
	EventResumeGame
	EventRestartGame
	EventGameTick
)

// PlaceTowerData contains data for placing a tower
type PlaceTowerData struct {
	TowerType TowerType `json:"towerType"`
	Position  Position  `json:"position"`
}

// =============================================================================
// 생성자
// =============================================================================

// NewTowerDefenseGame creates a new tower defense game
func NewTowerDefenseGame() *TowerDefenseGame {
	return &TowerDefenseGame{
		State:        NewGameState(),
		eventChannel: make(chan GameEvent, 100),
		running:      false,
		currentWave:  GetWaveDefinition(1),
	}
}

// =============================================================================
// 게임 제어
// =============================================================================

// Start begins the game
func (g *TowerDefenseGame) Start() {
	if g.running {
		return
	}
	
	g.running = true
	g.State.GameOver = false
	g.State.Paused = false
	
	// Start the game loop
	go g.gameLoop()
}

// Stop ends the game
func (g *TowerDefenseGame) Stop() {
	g.running = false
}

// Pause toggles the game pause state
func (g *TowerDefenseGame) Pause() {
	g.State.Paused = !g.State.Paused
}

// Restart resets the game to initial state
func (g *TowerDefenseGame) Restart() {
	g.Stop()
	g.State = NewGameState()
	g.currentWave = GetWaveDefinition(1)
	g.Start()
}

// SendEvent sends an event to the game
func (g *TowerDefenseGame) SendEvent(event GameEvent) {
	select {
	case g.eventChannel <- event:
	default:
		// Channel is full, ignore event
	}
}

// =============================================================================
// 게임 루프
// =============================================================================

// gameLoop runs the main game loop
func (g *TowerDefenseGame) gameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()
	
	for g.running && !g.State.GameOver {
		select {
		case event := <-g.eventChannel:
			if !g.State.Paused || event.Type == EventPauseGame || event.Type == EventResumeGame || event.Type == EventRestartGame {
				g.handleEvent(event)
			}
		case <-ticker.C:
			if !g.State.Paused {
				g.update()
			}
		}
	}
}

// update updates the game state
func (g *TowerDefenseGame) update() {
	now := time.Now()
	deltaTime := now.Sub(g.State.LastUpdate).Seconds()
	g.State.LastUpdate = now
	
	// Update wave spawning
	g.updateWaveSpawning(deltaTime)
	
	// Update enemies
	g.updateEnemies(deltaTime)
	
	// Update towers (includes pending damage processing)
	g.updateTowers(deltaTime)
	
	// Update projectiles (visual effects only)
	g.updateProjectiles(deltaTime)
	
	// Check for game over conditions
	g.checkGameOver()
}

// =============================================================================
// 이벤트 처리
// =============================================================================

// handleEvent processes a game event
func (g *TowerDefenseGame) handleEvent(event GameEvent) {
	switch event.Type {
	case EventPlaceTower:
		if data, ok := event.Data.(PlaceTowerData); ok {
			g.placeTower(data.TowerType, data.Position)
		}
	case EventStartWave:
		// 웨이브는 이제 자동으로 시작됨 - 더 이상 필요없음
	case EventForceNextWave:
		// 개발자 모드: 다음 웨이브 즉시 시작
		g.forceNextWave()
	case EventPauseGame:
		g.State.Paused = true
	case EventResumeGame:
		g.State.Paused = false
	case EventRestartGame:
		g.Restart()
	case EventGameTick:
		// Handled in main update loop
	}
}

// =============================================================================
// 타워 관리
// =============================================================================

// placeTower places a new tower at the given position
func (g *TowerDefenseGame) placeTower(towerType TowerType, pos Position) bool {
	stats := TowerStats[towerType]
	
	// Check if player has enough gold
	if g.State.Gold < stats.Cost {
		return false
	}
	
	// Check if position is valid (not on path, not too close to other towers)
	if !g.isValidTowerPosition(pos) {
		return false
	}
	
	// Create and place tower
	tower := NewTower(g.State.NextTowerID, towerType, pos)
	g.State.Towers = append(g.State.Towers, tower)
	g.State.NextTowerID++
	
	// Deduct gold
	g.State.Gold -= stats.Cost
	
	return true
}

// isValidTowerPosition checks if a position is valid for placing a tower
func (g *TowerDefenseGame) isValidTowerPosition(pos Position) bool {
	// Check if too close to path
	for _, pathPoint := range g.State.Map.Path {
		if Distance(pos, pathPoint) < 30 {
			return false
		}
	}
	
	// Check if too close to other towers
	for _, tower := range g.State.Towers {
		if Distance(pos, tower.Position) < 50 {
			return false
		}
	}
	
	// Check if within map bounds
	if pos.X < 25 || pos.X > float64(g.State.Map.Size.Width-25) ||
	   pos.Y < 25 || pos.Y > float64(g.State.Map.Size.Height-25) {
		return false
	}
	
	return true
}

// =============================================================================
// 웨이브 관리
// =============================================================================

// startNextWave starts the next wave automatically
func (g *TowerDefenseGame) startNextWave() {
	now := time.Now()
	
	// 새로운 웨이브 시작
	g.currentWave = GetWaveDefinition(g.State.CurrentWave)
	g.State.MaxEnemies = GetMaxEnemiesForWave(g.State.CurrentWave)
	g.State.EnemiesSpawned = 0
	g.State.EnemiesKilled = 0
	g.State.NextSpawnTime = now
	g.State.NextWaveTime = now.Add(60 * time.Second) // 다음 웨이브는 1분 후
	g.State.WaveTimeRemaining = 60.0
	
	g.State.CurrentWave++
}

// forceNextWave 개발자 모드: 다음 웨이브를 즉시 시작
func (g *TowerDefenseGame) forceNextWave() {
	now := time.Now()
	
	// 현재 웨이브 즉시 완료하고 다음 웨이브 시작
	g.currentWave = GetWaveDefinition(g.State.CurrentWave)
	g.State.MaxEnemies = GetMaxEnemiesForWave(g.State.CurrentWave)
	g.State.EnemiesSpawned = 0
	g.State.EnemiesKilled = 0
	g.State.NextSpawnTime = now
	g.State.NextWaveTime = now.Add(60 * time.Second) // 다음 웨이브는 1분 후
	g.State.WaveTimeRemaining = 60.0
	
	g.State.CurrentWave++
}

// updateWaveSpawning handles automatic wave progression and enemy spawning
func (g *TowerDefenseGame) updateWaveSpawning(deltaTime float64) {
	now := time.Now()
	
	// 웨이브 시간 업데이트
	g.State.WaveTimeRemaining = g.State.NextWaveTime.Sub(now).Seconds()
	
	// 새로운 웨이브 시작 체크
	if now.After(g.State.NextWaveTime) {
		g.startNextWave()
	}
	
	// 적 수가 최대치를 초과하면 게임오버
	if len(g.State.Enemies) > g.State.MaxEnemies {
		g.State.GameOver = true
		g.running = false
		return
	}
	
	// 적 스폰 (현재 웨이브에서 아직 스폰할 적이 남아있는 경우)
	if len(g.currentWave.Enemies) > 0 && now.After(g.State.NextSpawnTime) {
		for _, spawn := range g.currentWave.Enemies {
			if g.State.EnemiesSpawned < spawn.Count {
				// 적 수 한도 내에서만 스폰
				if len(g.State.Enemies) < g.State.MaxEnemies {
					enemy := NewEnemy(g.State.NextEnemyID, spawn.Type)
					g.State.Enemies = append(g.State.Enemies, enemy)
					g.State.NextEnemyID++
					g.State.EnemiesSpawned++
					
					// 다음 스폰 시간 설정
					g.State.NextSpawnTime = now.Add(time.Duration(spawn.Delay * float64(time.Second)))
					break
				}
				break
			}
		}
	}
}

// No longer needed - waves progress automatically every minute

// =============================================================================
// 적 업데이트
// =============================================================================

// updateEnemies updates all enemies
func (g *TowerDefenseGame) updateEnemies(deltaTime float64) {
	for _, enemy := range g.State.Enemies {
		if !enemy.IsAlive {
			continue
		}
		
		// 정사각형 경로를 따라 이동 (무한 루프)
		g.moveEnemyAlongSquarePath(enemy, deltaTime)
	}
	
	// 죽은 적들을 주기적으로 정리
	g.cleanupDeadEnemies()
}

// moveEnemyAlongSquarePath moves an enemy along the square loop path
func (g *TowerDefenseGame) moveEnemyAlongSquarePath(enemy *Enemy, deltaTime float64) {
	// 정사각형 경로 좌표 (시계방향)
	corners := []Position{
		{X: 100, Y: 100}, // 0: 좌상단
		{X: 500, Y: 100}, // 1: 우상단
		{X: 500, Y: 300}, // 2: 우하단
		{X: 100, Y: 300}, // 3: 좌하단
	}
	
	// 현재 변의 길이 계산 (정확한 거리)
	startCorner := corners[enemy.Corner]
	endCorner := corners[(enemy.Corner+1)%4]
	
	dx := endCorner.X - startCorner.X
	dy := endCorner.Y - startCorner.Y
	sideLength := math.Sqrt(dx*dx + dy*dy)
	
	// 이동 거리 계산
	moveDistance := enemy.Speed * deltaTime
	enemy.PathProgress += moveDistance / sideLength
	
	// 다음 꼭짓점에 도달한 경우 (부드러운 전환)
	if enemy.PathProgress >= 1.0 {
		// 정확히 꼭짓점에 위치시킨 후 다음 변으로 전환
		enemy.PathProgress = 0.0  // 정확히 0으로 리셋
		enemy.Corner = (enemy.Corner + 1) % 4
		
		// 한 바퀴 완주 체크
		if enemy.Corner == 0 {
			enemy.LapCount++
		}
		
		// 새로운 변의 시작점과 끝점 다시 계산
		startCorner = corners[enemy.Corner]
		endCorner = corners[(enemy.Corner+1)%4]
	}
	
	// PathProgress를 0.0-1.0 범위로 제한
	if enemy.PathProgress < 0.0 {
		enemy.PathProgress = 0.0
	}
	if enemy.PathProgress > 1.0 {
		enemy.PathProgress = 1.0
	}
	
	// 위치 계산 (선형 보간)
	enemy.Position.X = startCorner.X + (endCorner.X-startCorner.X)*enemy.PathProgress
	enemy.Position.Y = startCorner.Y + (endCorner.Y-startCorner.Y)*enemy.PathProgress
}

// =============================================================================
// 타워 업데이트
// =============================================================================

// updateTowers updates all towers
func (g *TowerDefenseGame) updateTowers(deltaTime float64) {
	for _, tower := range g.State.Towers {
		g.updateTower(tower, deltaTime)
		g.processPendingShots(tower)
	}
}

// updateTower updates a single tower
func (g *TowerDefenseGame) updateTower(tower *Tower, deltaTime float64) {
	// Find target
	target := g.findNearestEnemy(tower)
	tower.Target = target
	
	// Shoot if target is in range and fire rate allows
	if target != nil && time.Since(tower.LastShot).Seconds() >= 1.0/tower.FireRate {
		g.shootAtTarget(tower, target)
		tower.LastShot = time.Now()
	}
}

// findNearestEnemy finds the nearest enemy within tower's range
func (g *TowerDefenseGame) findNearestEnemy(tower *Tower) *Enemy {
	var nearestEnemy *Enemy
	nearestDistance := tower.Range
	
	for _, enemy := range g.State.Enemies {
		if !enemy.IsAlive {
			continue
		}
		
		distance := Distance(tower.Position, enemy.Position)
		if distance <= tower.Range && distance < nearestDistance {
			nearestEnemy = enemy
			nearestDistance = distance
		}
	}
	
	return nearestEnemy
}

// shootAtTarget creates both delayed damage and visual projectile
func (g *TowerDefenseGame) shootAtTarget(tower *Tower, target *Enemy) {
	stats := TowerStats[tower.Type]
	now := time.Now()
	
	// 거리에 따른 발사체 도달 시간 계산
	distance := Distance(tower.Position, target.Position)
	projectileSpeed := 200.0 // pixels per second
	travelTime := distance / projectileSpeed
	
	// 확정적 데미지를 예약 (일정 시간 후 적용)
	pendingShot := PendingShot{
		Target:  target,
		Damage:  tower.Damage,
		HitTime: now.Add(time.Duration(travelTime * float64(time.Second))),
		ShotID:  g.State.NextProjectileID,
	}
	tower.PendingShots = append(tower.PendingShots, pendingShot)
	
	// 시각적 발사체 생성 (연출용)
	projectile := NewProjectile(
		g.State.NextProjectileID,
		tower.Position,
		target,
		stats.Color,
		travelTime, // 발사체 생존 시간
	)
	
	g.State.Projectiles = append(g.State.Projectiles, projectile)
	g.State.NextProjectileID++
}

// processPendingShots handles delayed damage application
func (g *TowerDefenseGame) processPendingShots(tower *Tower) {
	now := time.Now()
	remainingShots := make([]PendingShot, 0)
	
	for _, shot := range tower.PendingShots {
		if now.After(shot.HitTime) {
			// 타겟이 아직 살아있으면 데미지 적용
			if shot.Target != nil && shot.Target.IsAlive {
				// 추가 검증: 타겟이 여전히 사거리 내에 있는지 확인
				if Distance(tower.Position, shot.Target.Position) <= tower.Range {
					g.damageEnemy(shot.Target, shot.Damage)
				}
			}
		} else {
			// 아직 시간이 안된 샷은 보존
			remainingShots = append(remainingShots, shot)
		}
	}
	
	tower.PendingShots = remainingShots
}

// =============================================================================
// 발사체 업데이트
// =============================================================================

// updateProjectiles updates all projectiles (visual effects only)
func (g *TowerDefenseGame) updateProjectiles(deltaTime float64) {
	for _, projectile := range g.State.Projectiles {
		if !projectile.IsActive {
			continue
		}
		
		g.updateVisualProjectile(projectile, deltaTime)
	}
	
	// Clean up inactive projectiles
	g.cleanupProjectiles()
}

// updateVisualProjectile updates a projectile for visual effects only
func (g *TowerDefenseGame) updateVisualProjectile(projectile *Projectile, deltaTime float64) {
	now := time.Now()
	
	// 수명이 다하면 비활성화
	if now.Sub(projectile.StartTime).Seconds() >= projectile.LifeTime {
		projectile.IsActive = false
		return
	}
	
	// 타겟이 없거나 죽었으면 비활성화
	if projectile.Target == nil || !projectile.Target.IsAlive {
		projectile.IsActive = false
		return
	}
	
	// 타겟을 향해 이동 (순수 시각적 효과)
	target := projectile.Target
	dx := target.Position.X - projectile.Position.X
	dy := target.Position.Y - projectile.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	if distance > 5 {
		// 타겟을 향해 이동
		moveDistance := projectile.Speed * deltaTime
		if distance > 0 { // 0으로 나누기 방지
			projectile.Position.X += (dx / distance) * moveDistance
			projectile.Position.Y += (dy / distance) * moveDistance
		}
	} else {
		// 타겟에 도달하면 비활성화
		projectile.IsActive = false
	}
}

// damageEnemy applies damage to an enemy
func (g *TowerDefenseGame) damageEnemy(enemy *Enemy, damage float64) {
	enemy.Health -= damage
	
	if enemy.Health <= 0 {
		enemy.IsAlive = false
		g.State.Gold += enemy.Reward
		g.State.Score += enemy.Reward * 5
		g.State.EnemiesKilled++
		
		// 죽은 적을 타겟으로 하는 모든 발사체를 비활성화 (시각적 효과)
		for _, projectile := range g.State.Projectiles {
			if projectile.Target == enemy {
				projectile.IsActive = false
			}
		}
		
		// 죽은 적을 타겟으로 하는 모든 예정된 샷을 제거
		for _, tower := range g.State.Towers {
			remainingShots := make([]PendingShot, 0)
			for _, shot := range tower.PendingShots {
				if shot.Target != enemy {
					remainingShots = append(remainingShots, shot)
				}
			}
			tower.PendingShots = remainingShots
		}
	}
}

// =============================================================================
// 정리 함수들
// =============================================================================

// cleanupDeadEnemies removes dead enemies from the game state
func (g *TowerDefenseGame) cleanupDeadEnemies() {
	aliveEnemies := make([]*Enemy, 0)
	for _, enemy := range g.State.Enemies {
		if enemy.IsAlive {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}
	g.State.Enemies = aliveEnemies
}

// cleanupProjectiles removes inactive projectiles
func (g *TowerDefenseGame) cleanupProjectiles() {
	activeProjectiles := make([]*Projectile, 0)
	for _, projectile := range g.State.Projectiles {
		if projectile.IsActive {
			activeProjectiles = append(activeProjectiles, projectile)
		}
	}
	g.State.Projectiles = activeProjectiles
}

// =============================================================================
// 게임 상태 체크
// =============================================================================

// checkGameOver checks for game over conditions
func (g *TowerDefenseGame) checkGameOver() {
	// 게임오버 조건: 적 수가 최대치를 초과
	if len(g.State.Enemies) > g.State.MaxEnemies {
		g.State.GameOver = true
		g.running = false
	}
}

// =============================================================================
// 공개 접근자 함수들
// =============================================================================

// GetGameState returns a copy of the current game state
func (g *TowerDefenseGame) GetGameState() *GameState {
	return g.State
}

// CanPlaceTower checks if a tower can be placed at the given position
func (g *TowerDefenseGame) CanPlaceTower(towerType TowerType, pos Position) bool {
	stats := TowerStats[towerType]
	return g.State.Gold >= stats.Cost && g.isValidTowerPosition(pos)
}

// GetTowerCost returns the cost of a tower type
func (g *TowerDefenseGame) GetTowerCost(towerType TowerType) int {
	return TowerStats[towerType].Cost
}