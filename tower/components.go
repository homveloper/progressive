package tower

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// =============================================================================
// 타워 디펜스 메인 페이지
// =============================================================================

// TowerDefensePage represents the main tower defense game page
type TowerDefensePage struct {
	app.Compo
	game           *TowerDefenseGame
	selectedTower  TowerType
	mousePos       Position
	showTowerRange bool
}

// =============================================================================
// 게임 보드 컴포넌트
// =============================================================================

// GameBoard renders the main game area
type GameBoard struct {
	app.Compo
	State          *GameState
	SelectedTower  TowerType
	MousePos       Position
	ShowTowerRange bool
	OnTowerPlace   func(TowerType, Position)
}

// =============================================================================
// HUD 컴포넌트들
// =============================================================================

// GameHUD renders the game's heads-up display
type GameHUD struct {
	app.Compo
	State       *GameState
	OnStartWave func()
	OnForceWave func() // 개발자 모드: 다음 웨이브 강제 시작
	OnPause     func()
	OnRestart   func()
}

// TowerSelector renders the tower selection panel
type TowerSelector struct {
	app.Compo
	SelectedTower TowerType
	Gold          int
	OnSelect      func(TowerType)
}

// GameStats renders game statistics
type GameStats struct {
	app.Compo
	State *GameState
}

// =============================================================================
// TowerDefensePage 구현
// =============================================================================

func (p *TowerDefensePage) OnMount(ctx app.Context) {
	if p.game == nil {
		p.game = NewTowerDefenseGame()
		p.game.Start()
	}

	p.selectedTower = TowerArcher

	// Start game update loop
	go p.updateLoop(ctx)
}

func (p *TowerDefensePage) OnDismount() {
	if p.game != nil {
		p.game.Stop()
	}
}

func (p *TowerDefensePage) updateLoop(ctx app.Context) {
	for {
		if p.game != nil && p.game.running {
			ctx.Dispatch(func(ctx app.Context) {
				ctx.Update()
			})
		}
		time.Sleep(16 * time.Millisecond) // ~60fps
	}
}

func (p *TowerDefensePage) handleMouseMove(ctx app.Context, e app.Event) {
	rect := e.Get("currentTarget").Call("getBoundingClientRect")
	p.mousePos.X = e.Get("clientX").Float() - rect.Get("left").Float()
	p.mousePos.Y = e.Get("clientY").Float() - rect.Get("top").Float()

	// Show tower range when hovering over game board
	p.showTowerRange = true
}

func (p *TowerDefensePage) handleClick(ctx app.Context, e app.Event) {
	if p.game.CanPlaceTower(p.selectedTower, p.mousePos) {
		p.game.SendEvent(GameEvent{
			Type: EventPlaceTower,
			Data: PlaceTowerData{
				TowerType: p.selectedTower,
				Position:  p.mousePos,
			},
		})
	}
}

func (p *TowerDefensePage) Render() app.UI {
	if p.game == nil {
		return app.Div().Text("Loading...")
	}

	state := p.game.GetGameState()

	return app.Div().Class("tower-defense-container").Body(
		app.H1().Class("game-title").Text("🏰 Tower Defense"),

		app.Div().Class("game-layout").Body(
			// 게임 보드
			app.Div().Class("game-board-container").Body(
				&GameBoard{
					State:          state,
					SelectedTower:  p.selectedTower,
					MousePos:       p.mousePos,
					ShowTowerRange: p.showTowerRange,
					OnTowerPlace:   p.handleTowerPlace,
				},
			),

			// 사이드바
			app.Div().Class("game-sidebar").Body(
				&GameStats{
					State: state,
				},

				&TowerSelector{
					SelectedTower: p.selectedTower,
					Gold:          state.Gold,
					OnSelect:      p.handleTowerSelect,
				},

				&GameHUD{
					State:       state,
					OnStartWave: p.handleStartWave,
					OnForceWave: p.handleForceWave,
					OnPause:     p.handlePause,
					OnRestart:   p.handleRestart,
				},

				// 컨트롤 도움말
				app.Div().Class("controls-help").Body(
					app.H4().Text("게임 방법"),
					app.Div().Class("help-item").Text("🏰 타워를 선택하고 맵에 클릭하여 배치"),
					app.Div().Class("help-item").Text("🔄 적들이 정사각형을 돌며 무한 순환"),
					app.Div().Class("help-item").Text("💰 적을 처치하면 골드를 획득"),
					app.Div().Class("help-item").Text("⚠️ 적 수가 최대치 초과시 게임오버"),
					app.Div().Class("help-item").Text("🌊 1분마다 웨이브 자동 시작"),
				),
			),
		),
	)
}

func (p *TowerDefensePage) handleTowerPlace(towerType TowerType, pos Position) {
	p.game.SendEvent(GameEvent{
		Type: EventPlaceTower,
		Data: PlaceTowerData{
			TowerType: towerType,
			Position:  pos,
		},
	})
}

func (p *TowerDefensePage) handleTowerSelect(towerType TowerType) {
	p.selectedTower = towerType
}

func (p *TowerDefensePage) handleStartWave() {
	// 웨이브는 이제 자동으로 시작됨 - 더 이상 필요없음
}

func (p *TowerDefensePage) handleForceWave() {
	// 개발자 모드: 다음 웨이브 강제 시작
	p.game.SendEvent(GameEvent{Type: EventForceNextWave})
}

func (p *TowerDefensePage) handlePause() {
	if p.game.State.Paused {
		p.game.SendEvent(GameEvent{Type: EventResumeGame})
	} else {
		p.game.SendEvent(GameEvent{Type: EventPauseGame})
	}
}

func (p *TowerDefensePage) handleRestart() {
	p.game.SendEvent(GameEvent{Type: EventRestartGame})
}

// =============================================================================
// GameBoard 구현
// =============================================================================

func (b *GameBoard) Render() app.UI {
	// HTML div 요소들을 구성
	var elements []app.UI

	// 경로 렌더링
	elements = append(elements, b.renderPath())

	// 타워들 렌더링
	for _, tower := range b.State.Towers {
		elements = append(elements, b.renderTower(tower))
	}

	// 적들 렌더링
	for _, enemy := range b.State.Enemies {
		if enemy.IsAlive {
			elements = append(elements, b.renderEnemy(enemy))
		}
	}

	// 발사체들 렌더링
	for _, projectile := range b.State.Projectiles {
		if projectile.IsActive {
			elements = append(elements, b.renderProjectile(projectile))
		}
	}

	// 마우스 위치에 배치할 타워 미리보기
	if b.ShowTowerRange {
		elements = append(elements, b.renderTowerPreview())
	}

	return app.Div().Class("game-board").
		OnMouseMove(b.handleMouseMove).
		OnClick(b.handleClick).
		Style("width", fmt.Sprintf("%dpx", b.State.Map.Size.Width)).
		Style("height", fmt.Sprintf("%dpx", b.State.Map.Size.Height)).
		Body(elements...)
}

func (b *GameBoard) renderPath() app.UI {
	if len(b.State.Map.Path) < 2 {
		return app.Text("")
	}

	var pathElements []app.UI

	// 경로를 작은 사각형들로 렌더링
	for i := 0; i < len(b.State.Map.Path)-1; i++ {
		start := b.State.Map.Path[i]
		end := b.State.Map.Path[i+1]

		// 두 점 사이의 거리와 방향 계산
		dx := end.X - start.X
		dy := end.Y - start.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// 작은 조각들로 경로 렌더링
		segments := int(distance / 10) // 10픽셀마다 하나씩
		for j := 0; j <= segments; j++ {
			t := float64(j) / float64(segments)
			x := start.X + dx*t
			y := start.Y + dy*t

			pathElements = append(pathElements,
				app.Div().
					Class("path-segment").
					Style("left", fmt.Sprintf("%.0fpx", x-10)).
					Style("top", fmt.Sprintf("%.0fpx", y-10)),
			)
		}
	}

	return app.Div().Body(pathElements...)
}

func (b *GameBoard) renderTower(tower *Tower) app.UI {
	stats := TowerStats[tower.Type]

	return app.Div().
		Class("tower").
		Style("left", fmt.Sprintf("%.0fpx", tower.Position.X-15)).
		Style("top", fmt.Sprintf("%.0fpx", tower.Position.Y-15)).
		Style("background-color", stats.Color).
		Body(
			app.Span().Class("tower-level").Text(strconv.Itoa(tower.Level)),
		)
}

func (b *GameBoard) renderEnemy(enemy *Enemy) app.UI {
	stats := EnemyStats[enemy.Type]
	healthPercent := enemy.Health / enemy.MaxHealth

	return app.Div().
		Class("enemy").
		Style("left", fmt.Sprintf("%.1fpx", enemy.Position.X-8)).     // 소수점 1자리 정밀도
		Style("top", fmt.Sprintf("%.1fpx", enemy.Position.Y-8)).      // 소수점 1자리 정밀도
		Style("background-color", stats.Color).
		Body(
			// 체력바
			app.Div().Class("health-bar").
				Style("width", fmt.Sprintf("%.1f%%", healthPercent*100)).
				Style("background-color", b.getHealthBarColor(healthPercent)),
		)
}

func (b *GameBoard) renderProjectile(projectile *Projectile) app.UI {
	return app.Div().
		Class("projectile").
		Style("left", fmt.Sprintf("%.1fpx", projectile.Position.X-3)).  // 소수점 1자리 정밀도
		Style("top", fmt.Sprintf("%.1fpx", projectile.Position.Y-3)).   // 소수점 1자리 정밀도
		Style("background-color", projectile.Color)
}

func (b *GameBoard) renderTowerPreview() app.UI {
	stats := TowerStats[b.SelectedTower]

	return app.Div().
		Class("tower-preview").
		Style("left", fmt.Sprintf("%.0fpx", b.MousePos.X-15)).
		Style("top", fmt.Sprintf("%.0fpx", b.MousePos.Y-15)).
		Style("border-color", stats.Color).
		Body(
			// 범위 표시
			app.Div().Class("tower-range-preview").
				Style("width", fmt.Sprintf("%.0fpx", stats.Range*2)).
				Style("height", fmt.Sprintf("%.0fpx", stats.Range*2)).
				Style("left", fmt.Sprintf("%.0fpx", 15-stats.Range)).
				Style("top", fmt.Sprintf("%.0fpx", 15-stats.Range)),
		)
}

func (b *GameBoard) getHealthBarColor(healthPercent float64) string {
	if healthPercent > 0.6 {
		return "#00ff00" // 녹색
	} else if healthPercent > 0.3 {
		return "#ffff00" // 노란색
	} else {
		return "#ff0000" // 빨간색
	}
}

func (b *GameBoard) handleMouseMove(ctx app.Context, e app.Event) {
	if b.OnTowerPlace != nil {
		// 마우스 좌표 업데이트는 부모에서 처리
	}
}

func (b *GameBoard) handleClick(ctx app.Context, e app.Event) {
	if b.OnTowerPlace != nil {
		rect := e.Get("currentTarget").Call("getBoundingClientRect")
		x := e.Get("clientX").Float() - rect.Get("left").Float()
		y := e.Get("clientY").Float() - rect.Get("top").Float()

		b.OnTowerPlace(b.SelectedTower, Position{X: x, Y: y})
	}
}

// =============================================================================
// GameStats 구현
// =============================================================================

func (s *GameStats) Render() app.UI {
	return app.Div().Class("game-stats").Body(
		app.H3().Text("📊 게임 상태"),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("💰 골드:"),
			app.Span().Class("stat-value").Text(strconv.Itoa(s.State.Gold)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("📊 최대 적:"),
			app.Span().Class("stat-value").Text(strconv.Itoa(s.State.MaxEnemies)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("🏆 점수:"),
			app.Span().Class("stat-value").Text(strconv.Itoa(s.State.Score)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("🌊 웨이브:"),
			app.Span().Class("stat-value").Text(strconv.Itoa(s.State.CurrentWave)),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("👹 적:"),
			app.Span().Class("stat-value").Text(strconv.Itoa(len(s.State.Enemies))),
		),

		app.Div().Class("stat-item").Body(
			app.Span().Class("stat-label").Text("⏰ 다음 웨이브:"),
			app.Span().Class("stat-value").Text(fmt.Sprintf("%.0f초", s.State.WaveTimeRemaining)),
		),
	)
}

// =============================================================================
// TowerSelector 구현
// =============================================================================

func (t *TowerSelector) Render() app.UI {
	return app.Div().Class("tower-selector").Body(
		app.H3().Text("🏰 타워 선택"),

		app.Div().Class("tower-buttons").Body(
			t.renderTowerButton(TowerArcher),
			t.renderTowerButton(TowerMage),
			t.renderTowerButton(TowerCannon),
		),
	)
}

func (t *TowerSelector) renderTowerButton(towerType TowerType) app.UI {
	stats := TowerStats[towerType]
	isSelected := t.SelectedTower == towerType
	canAfford := t.Gold >= stats.Cost

	buttonClass := "tower-button"
	if isSelected {
		buttonClass += " selected"
	}
	if !canAfford {
		buttonClass += " disabled"
	}

	return app.Button().
		Class(buttonClass).
		Disabled(!canAfford).
		OnClick(func(ctx app.Context, e app.Event) {
			if t.OnSelect != nil && canAfford {
				t.OnSelect(towerType)
			}
		}).
		Body(
			app.Div().Class("tower-info").Body(
				app.Div().Class("tower-name").Text(stats.Name),
				app.Div().Class("tower-cost").Text(fmt.Sprintf("💰 %d", stats.Cost)),
				app.Div().Class("tower-stats").Body(
					app.Span().Text(fmt.Sprintf("⚔️ %.0f", stats.Damage)),
					app.Span().Text(fmt.Sprintf("🎯 %.0f", stats.Range)),
					app.Span().Text(fmt.Sprintf("⚡ %.1f/s", stats.FireRate)),
				),
			),
		)
}

// =============================================================================
// GameHUD 구현
// =============================================================================

func (h *GameHUD) Render() app.UI {
	return app.Div().Class("game-hud").Body(
		app.H3().Text("🎮 게임 컨트롤"),

		app.Div().Class("control-buttons").Body(
			// 자동 웨이브 정보
			app.Div().Class("wave-info").Body(
				app.P().Text("🌊 웨이브는 1분마다 자동 시작"),
				app.P().Text(fmt.Sprintf("다음 웨이브: %.0f초 후", h.State.WaveTimeRemaining)),
			),

			// 개발자 모드: 다음 웨이브 강제 시작 버튼
			app.Button().
				Class("control-button dev-button").
				OnClick(func(ctx app.Context, e app.Event) {
					if h.OnForceWave != nil {
						h.OnForceWave()
					}
				}).
				Text("🚀 개발자: 다음 웨이브 시작"),

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
		),

		// 게임 오버 메시지
		app.If(h.State.GameOver,
			func() app.UI {
				return app.Div().Class("game-over").Body(
					app.H2().Text("💀 게임 오버!"),
					app.P().Text(fmt.Sprintf("최종 점수: %d", h.State.Score)),
				)
			},
		),
	)
}
