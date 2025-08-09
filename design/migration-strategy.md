# Migration Strategy: go-app → Modern Stack

## 🎯 전환 전략 설계

### 핵심 원칙
1. **Zero Downtime**: 서비스 중단 없는 점진적 전환
2. **Data Preservation**: 기존 사용자 데이터 완전 보존  
3. **Progressive Migration**: 페이지별 단계적 마이그레이션
4. **Risk Mitigation**: 롤백 가능한 안전한 전환
5. **User Experience**: 사용자 혼란 최소화

## 🏗️ 마이그레이션 아키텍처 설계

### 1. 하이브리드 라우팅 시스템
```go
// main.go - 하이브리드 라우터 설계
package main

import (
    "progressive/internal/handlers"    // 새로운 Modern Stack 핸들러
    "progressive/spreadsheet"          // 기존 go-app 컴포넌트
    
    "github.com/gin-gonic/gin"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
    // 새로운 Gin 라우터 (Modern Stack)
    r := gin.Default()
    
    // 정적 파일 서빙
    r.Static("/static", "./web/static")
    
    // === MODERN STACK ROUTES ===
    modernHandler := handlers.NewModernHandler()
    
    // Phase 1: 정적/단순 페이지들 먼저 전환
    r.GET("/", modernHandler.LandingPage)           // ✅ 새 버전
    r.GET("/dashboard", modernHandler.Dashboard)    // ✅ 새 버전
    r.GET("/settings", modernHandler.Settings)     // ✅ 새 버전
    
    // Phase 2: 중간 복잡도 페이지들
    r.GET("/tableedit", modernHandler.TableEdit)   // ✅ 새 버전 (우선순위)
    
    // Phase 3: 복잡한 페이지들 (나중에 전환)
    // r.GET("/editor", modernHandler.Editor)      // 🟡 아직 go-app
    
    // API 엔드포인트 (HTMX용)
    api := r.Group("/api")
    {
        api.POST("/spreadsheet", modernHandler.CreateSpreadsheet)
        api.GET("/spreadsheet/:id/grid", modernHandler.GetGrid)
        api.PUT("/cell/:id", modernHandler.UpdateCell)
        api.POST("/upload-schema", modernHandler.UploadSchema)
    }
    
    // === LEGACY GO-APP ROUTES ===
    // 아직 마이그레이션되지 않은 페이지들
    app.Route("/legacy/editor", func() app.Composer { 
        return &editor.Editor{} 
    })
    app.Route("/editor", func() app.Composer { 
        return &editor.Editor{}  // 임시로 기존 버전 유지
    })
    
    // go-app WebAssembly 핸들러 (기존 페이지용)
    r.Any("/app.wasm", gin.WrapH(&app.Handler{
        Name: "Progressive Spreadsheet Legacy",
        Styles: []string{"/static/css/tailwind.css"},
    }))
    
    // 하이브리드 서버 실행
    go app.RunWhenOnBrowser() // 클라이언트에서 go-app 실행
    
    r.Run(":8080") // 서버에서 Modern Stack 실행
}
```

### 2. 데이터 마이그레이션 시스템
```go
// internal/migration/data_migration.go
package migration

type DataMigrationService struct {
    legacyStorage *storage.LocalStorage
    modernDB      *sql.DB
}

// 기존 LocalStorage 데이터를 현대적 DB로 마이그레이션
func (d *DataMigrationService) MigrateUserData(userID string) error {
    // 1. 기존 LocalStorage에서 데이터 추출
    legacyData, err := d.legacyStorage.GetAllUserData(userID)
    if err != nil {
        return err
    }
    
    // 2. 데이터 구조 변환
    modernData := d.transformLegacyData(legacyData)
    
    // 3. 새로운 DB에 저장
    return d.modernDB.SaveUserData(userID, modernData)
}

// 양방향 동기화 (전환 기간 동안)
func (d *DataMigrationService) SyncBidirectional(userID string) error {
    // LocalStorage ↔ Modern DB 양방향 동기화
    // 사용자가 어느 버전을 사용하든 데이터 일관성 보장
    return nil
}
```

### 3. 점진적 UI 전환 시스템
```html
<!-- 사용자에게 새 버전 체험 옵션 제공 -->
<div class="migration-banner bg-blue-50 border border-blue-200 p-4 rounded-lg mb-4">
    <div class="flex items-center justify-between">
        <div>
            <h3 class="font-semibold text-blue-800">🚀 새로운 버전 체험해보기</h3>
            <p class="text-blue-600 text-sm">더 빠르고 향상된 스프레드시트 편집기를 만나보세요!</p>
        </div>
        <div class="flex gap-2">
            <button onclick="tryNewVersion()" 
                    class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
                새 버전 체험
            </button>
            <button onclick="stayLegacy()" 
                    class="text-blue-500 px-4 py-2 rounded border border-blue-300 hover:bg-blue-50">
                기존 버전 사용
            </button>
        </div>
    </div>
</div>

<script>
function tryNewVersion() {
    // 사용자 선택을 기록하고 새 버전으로 리다이렉트
    localStorage.setItem('preferred_version', 'modern');
    window.location.href = '/tableedit?version=modern';
}

function stayLegacy() {
    localStorage.setItem('preferred_version', 'legacy'); 
    // 기존 go-app 버전 계속 사용
}
</script>
```

## 📅 상세 마이그레이션 로드맵

### Phase 1: 기반 구축 (Week 1-2)
```
🎯 목표: Modern Stack 기본 인프라 구축

📋 할 일:
✅ Gin HTTP 서버 셋업
✅ Templ 템플릿 시스템 도입  
✅ TailwindCSS 빌드 시스템
✅ 기본 라우팅 및 미들웨어
✅ 데이터베이스 설계 (SQLite → PostgreSQL)
✅ 하이브리드 라우팅 시스템 구현

🏗️ 기술 작업:
- package.json, tailwind.config.js 설정
- internal/ 디렉토리 구조 생성
- web/templates/ Templ 템플릿 구조
- Docker 컨테이너 설정
- CI/CD 파이프라인 준비
```

### Phase 2: 정적 페이지 마이그레이션 (Week 3)
```
🎯 목표: 복잡하지 않은 페이지들 먼저 전환

📋 할 일:
✅ 랜딩 페이지 (/dashboard) 
✅ 설정 페이지 (/settings)
✅ 사용자 프로필 페이지
✅ 정적 콘텐츠 페이지들

🏗️ 구현:
// web/templates/pages/dashboard.templ
package pages

templ Dashboard() {
    @layouts.Base("Dashboard") {
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <h1 class="text-3xl font-bold text-gray-900 mb-8">
                Progressive Spreadsheet Dashboard
            </h1>
            
            <!-- 대시보드 컨텐츠 -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                @components.DashboardCard("Recent Files", "12 files")
                @components.DashboardCard("Storage Used", "2.3 GB")
                @components.DashboardCard("Collaborators", "8 users")
            </div>
        </div>
    }
}
```

### Phase 3: 테이블 편집 페이지 마이그레이션 (Week 4-5)
```
🎯 목표: 핵심 기능인 테이블 편집을 Modern Stack으로 전환

📋 할 일:
✅ 기본 그리드 렌더링 (서버사이드)
✅ HTMX 기반 셀 편집 인터랙션
✅ 파일 업로드 기능 (HTMX)
✅ 스키마 처리 및 검증
✅ 데이터 내보내기 기능
✅ 실시간 업데이트 (WebSocket)

🏗️ 핵심 구현:
// internal/handlers/tableedit.go
func (h *Handler) TableEditPage(c *gin.Context) {
    // 서버에서 초기 데이터 로드
    spreadsheetID := c.Query("id")
    data, err := h.service.GetSpreadsheetData(spreadsheetID)
    
    // Templ 템플릿 렌더링
    templ.Handler(pages.TableEdit(data)).ServeHTTP(c.Writer, c.Request)
}

func (h *Handler) UpdateCell(c *gin.Context) {
    // HTMX 요청 처리
    cellID := c.Param("id")
    newValue := c.PostForm("value")
    
    // 셀 업데이트
    cell, err := h.service.UpdateCell(cellID, newValue)
    
    // 부분 렌더링으로 업데이트된 셀만 반환
    templ.Handler(components.Cell(cell)).ServeHTTP(c.Writer, c.Request)
}
```

### Phase 4: 고급 편집기 마이그레이션 (Week 6-7)
```
🎯 목표: 가장 복잡한 Excel-like 편집기를 Modern Stack으로 완전 전환

📋 할 일:
✅ 복잡한 그리드 인터랙션 (드래그, 다중 선택)
✅ 수식 계산 엔진 (서버사이드)
✅ 차트 및 그래프 생성
✅ 협업 기능 (실시간 편집)
✅ 무한 스크롤 가상화
✅ 키보드 단축키 지원

🏗️ 고급 구현:
// 실시간 협업 WebSocket
func (h *Handler) HandleWebSocket(c *gin.Context) {
    conn, _ := websocket.Upgrade(c.Writer, c.Request, nil)
    userID := c.Query("user_id")
    
    h.collabService.AddUser(userID, conn)
    
    for {
        var msg CollaborationMessage
        if err := conn.ReadJSON(&msg); err != nil {
            break
        }
        
        // 다른 사용자들에게 변경사항 브로드캐스트
        h.collabService.BroadcastChange(msg)
    }
}
```

### Phase 5: 최적화 및 완성 (Week 8)
```
🎯 목표: 성능 최적화 및 기존 시스템 완전 교체

📋 할 일:
✅ 성능 벤치마크 및 최적화
✅ 캐싱 시스템 (Redis) 도입
✅ CDN 및 정적 자산 최적화
✅ 모니터링 및 로깅 시스템
✅ 기존 go-app 코드 제거
✅ 사용자 피드백 수집 및 개선

🏗️ 최적화:
// Redis 캐싱 
func (s *SpreadsheetService) GetGridWithCache(id string) (*Grid, error) {
    // 캐시에서 먼저 확인
    if cached := s.cache.Get("grid:" + id); cached != nil {
        return cached.(*Grid), nil
    }
    
    // 캐시 미스시 DB에서 로드 후 캐싱
    grid, err := s.repo.GetGrid(id)
    if err == nil {
        s.cache.Set("grid:"+id, grid, time.Hour)
    }
    
    return grid, err
}
```

## 🎛️ 위험 관리 및 롤백 전략

### 1. A/B 테스팅 시스템
```go
// Feature Flag 기반 A/B 테스팅
type FeatureFlag struct {
    Name    string
    Enabled bool
    Rollout float64 // 0.0 ~ 1.0
}

func (h *Handler) shouldUseModernStack(userID string) bool {
    flag := h.featureService.GetFlag("modern_stack_rollout")
    
    // 사용자 해시 기반 일관된 경험 제공
    hash := h.hashUserID(userID)
    return hash < flag.Rollout
}

func (h *Handler) TableEditRouter(c *gin.Context) {
    userID := c.GetHeader("User-ID")
    
    if h.shouldUseModernStack(userID) {
        // Modern Stack 버전
        h.ModernTableEdit(c)
    } else {
        // Legacy go-app 버전으로 리다이렉트
        c.Redirect(302, "/legacy/tableedit")
    }
}
```

### 2. 데이터 백업 및 롤백
```go
// 마이그레이션 중 데이터 백업
type BackupService struct {
    legacyDB  *sql.DB
    modernDB  *sql.DB
    backupDir string
}

func (b *BackupService) BackupBeforeMigration(userID string) error {
    // 1. 기존 데이터 전체 백업
    backup, err := b.legacyDB.DumpUserData(userID)
    if err != nil {
        return err
    }
    
    // 2. 백업 파일 저장 (타임스탬프 포함)
    filename := fmt.Sprintf("%s/user_%s_%d.backup", 
        b.backupDir, userID, time.Now().Unix())
    
    return ioutil.WriteFile(filename, backup, 0644)
}

func (b *BackupService) Rollback(userID string, backupTimestamp int64) error {
    // 백업 파일에서 데이터 복원
    filename := fmt.Sprintf("%s/user_%s_%d.backup", 
        b.backupDir, userID, backupTimestamp)
        
    backup, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    
    return b.legacyDB.RestoreUserData(userID, backup)
}
```

### 3. 헬스 체크 및 모니터링
```go
// 시스템 헬스 모니터링
func (h *Handler) HealthCheck(c *gin.Context) {
    health := struct {
        Status       string                 `json:"status"`
        Version      string                 `json:"version"`
        Dependencies map[string]interface{} `json:"dependencies"`
        Metrics      map[string]float64     `json:"metrics"`
    }{
        Status:  "healthy",
        Version: "modern-stack-v1.0",
        Dependencies: map[string]interface{}{
            "database": h.checkDBHealth(),
            "redis":    h.checkRedisHealth(), 
            "websocket": h.checkWSHealth(),
        },
        Metrics: map[string]float64{
            "response_time_avg": h.metrics.GetAvgResponseTime(),
            "error_rate":        h.metrics.GetErrorRate(),
            "concurrent_users":  h.metrics.GetConcurrentUsers(),
        },
    }
    
    c.JSON(200, health)
}
```

## 📊 성공 지표 및 측정

### 1. 기술 지표
```
성능 지표:
- 초기 로딩 시간: 5초 → 0.5초 (90% 개선)
- Time to Interactive: 5초 → 0.8초 (84% 개선) 
- 번들 크기: 23MB → 50KB (99.8% 감소)
- 메모리 사용량: 100MB → 20MB (80% 감소)

SEO 지표:
- Google PageSpeed Score: 30 → 95
- Lighthouse Performance: 40 → 95
- 검색 노출 가능 여부: 불가능 → 가능

개발자 경험:
- 디버깅 용이성: 어려움 → 쉬움
- 빌드 시간: 30초 → 5초
- Hot Reload: 제한적 → 즉시
```

### 2. 사용자 지표
```
사용성 지표:
- 이탈률 (Bounce Rate): 현재 → 목표 30% 감소
- 세션 시간: 현재 → 목표 50% 증가
- 기능 사용률: 그리드 편집 빈도 측정

만족도 지표:
- 사용자 피드백 점수: NPS 측정
- 새 버전 선택률: A/B 테스트 결과
- 버그 리포트 수: 감소 목표
```

## 🚀 최종 실행 계획

### 즉시 시작 가능한 작업들
```
Week 1 (기반 구축):
□ gin HTTP 서버 설정
□ templ 템플릿 시스템 도입
□ tailwind.config.js 설정
□ Docker 컨테이너 준비
□ CI/CD 파이프라인 설정

Week 2 (데이터 계층):
□ PostgreSQL 스키마 설계
□ Repository 패턴 구현
□ 데이터 마이그레이션 도구 개발
□ API 엔드포인트 기본 구조
□ 하이브리드 라우팅 시스템
```

## 🎯 결론

**이 마이그레이션 전략의 핵심 장점:**

1. **Zero Risk**: 기존 시스템 유지하면서 점진적 전환
2. **User Choice**: 사용자가 원하는 버전 선택 가능  
3. **Data Safety**: 완전한 데이터 보존 및 백업
4. **Performance**: 10배 향상된 성능 기대
5. **Future Ready**: 현대적 웹 표준 기반 확장성

**이 전략으로 Progressive Spreadsheet는 안전하고 확실하게 차세대 웹 애플리케이션으로 진화할 수 있습니다!** 🚀