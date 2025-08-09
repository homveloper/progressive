# Implementation Roadmap: Modern Web Stack

## 🗓️ 8주 완성 로드맵

### 📊 프로젝트 개요
- **목표**: go-app → Go + HTMX + Templ + TailwindCSS 완전 전환
- **기간**: 8주 (56일)
- **방식**: 점진적 마이그레이션 (Zero Downtime)
- **결과**: 10배 빠른 로딩 + 현대적 웹 애플리케이션

## 🚀 Week 1-2: 기반 구축 (Foundation)

### Week 1: Modern Stack 환경 설정
```bash
# Day 1-2: 프로젝트 구조 및 의존성
mkdir progressive-modern
cd progressive-modern

# Go 모듈 초기화
go mod init progressive-modern

# 핵심 의존성 추가
go get github.com/gin-gonic/gin
go get github.com/a-h/templ/cmd/templ
go get github.com/gorilla/websocket
go get github.com/lib/pq  # PostgreSQL

# Node.js 환경 (TailwindCSS용)
npm init -y
npm install -D tailwindcss @tailwindcss/forms @tailwindcss/typography
npm install -D postcss autoprefixer cssnano

# Day 3-4: 기본 서버 및 라우팅
```

#### 기본 서버 구조
```go
// cmd/server/main.go
package main

import (
    "log"
    "progressive-modern/internal/config"
    "progressive-modern/internal/handlers"
    "progressive-modern/internal/middleware"
    
    "github.com/gin-gonic/gin"
)

func main() {
    // 설정 로드
    cfg := config.Load()
    
    // Gin 라우터 설정
    r := gin.Default()
    
    // 미들웨어 설정
    r.Use(middleware.CORS())
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    
    // 정적 파일 서빙
    r.Static("/static", "./web/static")
    
    // 핸들러 초기화
    h := handlers.New()
    
    // 라우트 설정
    h.RegisterRoutes(r)
    
    log.Printf("Server starting on :%s", cfg.Port)
    r.Run(":" + cfg.Port)
}
```

#### TailwindCSS 설정
```javascript
// tailwind.config.js
module.exports = {
  content: [
    "./web/templates/**/*.templ",
    "./web/templates/**/*.html"
  ],
  theme: {
    extend: {
      colors: {
        excel: {
          primary: '#217346',
          secondary: '#f2f8f4', 
          accent: '#0078d4'
        }
      }
    }
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography')
  ]
}
```

### Week 2: 데이터베이스 및 기본 기능
```sql
-- Day 5-6: PostgreSQL 스키마 설계
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE spreadsheets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(200) NOT NULL,
    schema_data JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE cells (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    spreadsheet_id UUID REFERENCES spreadsheets(id),
    row_index INTEGER NOT NULL,
    col_index INTEGER NOT NULL,
    value TEXT,
    formula TEXT,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(spreadsheet_id, row_index, col_index)
);
```

#### Repository 패턴 구현
```go
// internal/repository/spreadsheet.go
type SpreadsheetRepository struct {
    db *sql.DB
}

func (r *SpreadsheetRepository) Create(s *models.Spreadsheet) error {
    query := `INSERT INTO spreadsheets (user_id, name, schema_data) 
              VALUES ($1, $2, $3) RETURNING id, created_at`
    return r.db.QueryRow(query, s.UserID, s.Name, s.SchemaData).
           Scan(&s.ID, &s.CreatedAt)
}

func (r *SpreadsheetRepository) GetByID(id string) (*models.Spreadsheet, error) {
    // 구현...
}
```

## 📱 Week 3: 첫 번째 페이지 마이그레이션

### Day 15-17: 대시보드 페이지
```go
// web/templates/pages/dashboard.templ
package pages

import "progressive-modern/internal/models"

templ Dashboard(user *models.User, spreadsheets []*models.Spreadsheet) {
    @layouts.Base("Dashboard") {
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <!-- 헤더 -->
            <div class="flex justify-between items-center mb-8">
                <h1 class="text-3xl font-bold text-gray-900">
                    안녕하세요, {user.Username}님!
                </h1>
                <button class="bg-excel-primary text-white px-4 py-2 rounded-lg hover:bg-excel-primary/90 transition-colors"
                        hx-get="/api/spreadsheet/new"
                        hx-target="#modal-container"
                        hx-swap="innerHTML">
                    새 스프레드시트
                </button>
            </div>
            
            <!-- 스프레드시트 그리드 -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                for _, sheet := range spreadsheets {
                    @components.SpreadsheetCard(sheet)
                }
            </div>
        </div>
        
        <!-- 모달 컨테이너 -->
        <div id="modal-container"></div>
    }
}
```

### Day 18-21: 설정 페이지 및 사용자 관리
```go
// internal/handlers/settings.go
func (h *Handler) SettingsPage(c *gin.Context) {
    userID := c.GetString("user_id") // 미들웨어에서 설정
    user, err := h.userService.GetByID(userID)
    if err != nil {
        c.HTML(500, "error.html", gin.H{"error": err.Error()})
        return
    }
    
    templ.Handler(pages.Settings(user)).ServeHTTP(c.Writer, c.Request)
}

func (h *Handler) UpdateSettings(c *gin.Context) {
    // HTMX POST 요청 처리
    var req models.UpdateUserRequest
    if err := c.ShouldBind(&req); err != nil {
        c.HTML(400, "error_partial.html", gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.userService.Update(c.GetString("user_id"), &req)
    if err != nil {
        c.HTML(500, "error_partial.html", gin.H{"error": err.Error()})
        return
    }
    
    // 부분 렌더링으로 업데이트된 폼 반환
    templ.Handler(components.SettingsForm(user, "Settings updated successfully!")).
         ServeHTTP(c.Writer, c.Request)
}
```

## 📊 Week 4-5: 테이블 편집 마이그레이션

### Day 22-25: 기본 그리드 시스템
```go
// web/templates/components/grid/table.templ
package grid

templ SpreadsheetGrid(data *models.SpreadsheetData) {
    <div class="overflow-auto border border-gray-300 rounded-lg max-h-96">
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50 sticky top-0">
                <tr>
                    <th class="w-12 px-2 py-2 text-xs font-medium text-gray-500">#</th>
                    for i, col := range data.Columns {
                        <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-32"
                            id={fmt.Sprintf("header-%d", i)}>
                            <div class="flex items-center justify-between">
                                <span>{col.Name}</span>
                                <button class="text-gray-400 hover:text-gray-600"
                                        hx-get={fmt.Sprintf("/api/column/%d/menu", i)}
                                        hx-target="#column-menu"
                                        hx-swap="innerHTML">
                                    ⋮
                                </button>
                            </div>
                        </th>
                    }
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                for rowIdx, row := range data.Rows {
                    <tr class="hover:bg-gray-50" id={fmt.Sprintf("row-%d", rowIdx)}>
                        <td class="px-2 py-2 text-xs text-gray-500 bg-gray-50 font-medium sticky left-0">
                            {fmt.Sprintf("%d", rowIdx+1)}
                        </td>
                        for colIdx, cell := range row {
                            @Cell(cell, rowIdx, colIdx)
                        }
                    </tr>
                }
            </tbody>
        </table>
    </div>
    
    <!-- 컨텍스트 메뉴 컨테이너 -->
    <div id="column-menu" class="hidden"></div>
}

templ Cell(cell *models.Cell, row, col int) {
    <td class="px-3 py-2 text-sm text-gray-900 border-r border-gray-100 cursor-cell relative group"
        id={fmt.Sprintf("cell-%d-%d", row, col)}
        hx-trigger="dblclick"
        hx-get={fmt.Sprintf("/api/cell/%d/%d/edit", row, col)}
        hx-target="this"
        hx-swap="outerHTML">
        
        <!-- 셀 값 표시 -->
        <span class="block truncate">{cell.DisplayValue}</span>
        
        <!-- 호버시 편집 버튼 -->
        <button class="absolute top-1 right-1 opacity-0 group-hover:opacity-100 w-4 h-4 bg-blue-500 text-white rounded-sm text-xs transition-opacity"
                hx-trigger="click"
                hx-get={fmt.Sprintf("/api/cell/%d/%d/edit", row, col)}
                hx-target={fmt.Sprintf("#cell-%d-%d", row, col)}
                hx-swap="outerHTML">
            ✎
        </button>
    </td>
}
```

### Day 26-28: HTMX 인터랙션 구현
```go
// internal/handlers/cell.go
func (h *Handler) EditCell(c *gin.Context) {
    row, _ := strconv.Atoi(c.Param("row"))
    col, _ := strconv.Atoi("col"))
    
    spreadsheetID := c.Query("spreadsheet_id")
    cell, err := h.cellService.Get(spreadsheetID, row, col)
    if err != nil {
        c.HTML(500, "error.html", gin.H{"error": err.Error()})
        return
    }
    
    // 편집 모드 셀 렌더링
    templ.Handler(components.EditableCell(cell, row, col)).
         ServeHTTP(c.Writer, c.Request)
}

func (h *Handler) UpdateCell(c *gin.Context) {
    row, _ := strconv.Atoi(c.Param("row")) 
    col, _ := strconv.Atoi(c.Param("col"))
    
    var req models.UpdateCellRequest
    if err := c.ShouldBind(&req); err != nil {
        c.HTML(400, "error.html", gin.H{"error": err.Error()})
        return
    }
    
    // 셀 업데이트
    cell, err := h.cellService.Update(req.SpreadsheetID, row, col, req.Value)
    if err != nil {
        c.HTML(500, "error.html", gin.H{"error": err.Error()})
        return
    }
    
    // WebSocket으로 다른 사용자들에게 실시간 알림
    h.wsService.BroadcastCellUpdate(req.SpreadsheetID, cell)
    
    // 업데이트된 셀 렌더링
    templ.Handler(components.Cell(cell, row, col)).
         ServeHTTP(c.Writer, c.Request)
}
```

### Day 29-35: 파일 업로드 및 스키마 처리
```go
// web/templates/components/upload.templ
templ FileUploadModal() {
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-md">
            <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold">스키마 파일 업로드</h3>
                <button class="text-gray-400 hover:text-gray-600"
                        hx-get="/api/modal/close"
                        hx-target="#modal-container"
                        hx-swap="innerHTML">✕</button>
            </div>
            
            <form hx-post="/api/upload-schema"
                  hx-encoding="multipart/form-data"
                  hx-target="#upload-result"
                  hx-swap="innerHTML"
                  class="space-y-4">
                
                <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
                    <input type="file" 
                           name="schema_file" 
                           accept=".json"
                           class="w-full" 
                           required/>
                    <p class="text-sm text-gray-500 mt-2">JSON 스키마 파일을 선택하세요</p>
                </div>
                
                <div class="flex justify-end gap-2">
                    <button type="button"
                            class="px-4 py-2 text-gray-600 border border-gray-300 rounded hover:bg-gray-50"
                            hx-get="/api/modal/close"
                            hx-target="#modal-container"
                            hx-swap="innerHTML">
                        취소
                    </button>
                    <button type="submit"
                            class="px-4 py-2 bg-excel-primary text-white rounded hover:bg-excel-primary/90">
                        업로드
                    </button>
                </div>
            </form>
            
            <div id="upload-result" class="mt-4"></div>
        </div>
    </div>
}
```

## 🏗️ Week 6-7: 고급 편집기 마이그레이션

### Day 36-42: 복잡한 Excel 기능들
```javascript
// web/static/js/spreadsheet.js - 키보드 단축키 지원
document.addEventListener('DOMContentLoaded', function() {
    // 키보드 네비게이션
    document.addEventListener('keydown', function(e) {
        const activeCell = document.querySelector('.cell-selected');
        if (!activeCell) return;
        
        const [row, col] = activeCell.id.split('-').slice(1).map(Number);
        let newRow = row, newCol = col;
        
        switch(e.key) {
            case 'ArrowUp':
                newRow = Math.max(0, row - 1);
                break;
            case 'ArrowDown': 
                newRow = row + 1;
                break;
            case 'ArrowLeft':
                newCol = Math.max(0, col - 1);
                break;
            case 'ArrowRight':
                newCol = col + 1;
                break;
            case 'Enter':
                if (activeCell.classList.contains('editing')) {
                    // 편집 완료
                    htmx.trigger(activeCell, 'cell:save');
                } else {
                    // 편집 시작
                    htmx.trigger(activeCell, 'dblclick');
                }
                e.preventDefault();
                return;
            case 'Escape':
                // 편집 취소
                htmx.trigger(activeCell, 'cell:cancel');
                e.preventDefault();
                return;
            default:
                return;
        }
        
        // 셀 선택 변경
        const newCell = document.getElementById(`cell-${newRow}-${newCol}`);
        if (newCell) {
            activeCell.classList.remove('cell-selected');
            newCell.classList.add('cell-selected');
        }
        
        e.preventDefault();
    });
});
```

### Day 43-49: 실시간 협업 시스템
```go
// internal/services/collaboration.go
type CollaborationService struct {
    clients    map[string]*websocket.Conn
    broadcasts chan *models.CellUpdate
    mu         sync.RWMutex
}

func (c *CollaborationService) HandleConnection(conn *websocket.Conn, userID string) {
    c.mu.Lock()
    c.clients[userID] = conn
    c.mu.Unlock()
    
    defer func() {
        c.mu.Lock()
        delete(c.clients, userID)
        c.mu.Unlock()
        conn.Close()
    }()
    
    for {
        var msg models.CollaborationMessage
        if err := conn.ReadJSON(&msg); err != nil {
            break
        }
        
        // 다른 사용자들에게 브로드캐스트
        c.BroadcastUpdate(&models.CellUpdate{
            SpreadsheetID: msg.SpreadsheetID,
            Row:          msg.Row,
            Col:          msg.Col,
            Value:        msg.Value,
            UserID:       userID,
        })
    }
}

func (c *CollaborationService) BroadcastUpdate(update *models.CellUpdate) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    for userID, conn := range c.clients {
        if userID != update.UserID { // 발신자 제외
            go func(conn *websocket.Conn) {
                conn.WriteJSON(update)
            }(conn)
        }
    }
}
```

## ⚡ Week 8: 최적화 및 완성

### Day 50-52: 성능 최적화
```go
// internal/middleware/cache.go
func CacheMiddleware(redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // GET 요청만 캐싱
        if c.Request.Method != "GET" {
            c.Next()
            return
        }
        
        cacheKey := fmt.Sprintf("page:%s", c.Request.URL.Path)
        
        // 캐시에서 확인
        if cached, err := redis.Get(cacheKey).Result(); err == nil {
            c.Header("X-Cache", "HIT")
            c.Data(200, "text/html", []byte(cached))
            return
        }
        
        // ResponseWriter 래핑하여 응답 캐시
        writer := &CacheWriter{ResponseWriter: c.Writer, redis: redis, key: cacheKey}
        c.Writer = writer
        
        c.Header("X-Cache", "MISS")
        c.Next()
    }
}

type CacheWriter struct {
    gin.ResponseWriter
    redis *redis.Client
    key   string
    body  []byte
}

func (w *CacheWriter) Write(data []byte) (int, error) {
    w.body = append(w.body, data...)
    
    // 성공적인 HTML 응답만 캐싱 (5분)
    if w.Status() == 200 {
        w.redis.Set(w.key, w.body, 5*time.Minute)
    }
    
    return w.ResponseWriter.Write(data)
}
```

### Day 53-56: 모니터링 및 배포
```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/progressive
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
  
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: progressive
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
      
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app

volumes:
  postgres_data:
  redis_data:
```

## 📊 최종 검증 및 측정

### 성능 벤치마크
```bash
# 로딩 시간 측정
lighthouse --chrome-flags="--headless" http://localhost:8080

# 부하 테스트
wrk -t12 -c400 -d30s --latency http://localhost:8080/

# 메모리 사용량 모니터링
docker stats progressive-modern_app_1

# 데이터베이스 성능
pg_stat_statements 활용한 쿼리 분석
```

### A/B 테스트 결과 예상
```
지표                  기존 go-app     Modern Stack    개선율
----------------------------------------------------------------
첫 로딩 시간           4.2초          0.6초          86% ↑
Time to Interactive   4.8초          0.9초          81% ↑ 
번들 크기             23MB           65KB           99.7% ↓
Lighthouse 점수       32             94             194% ↑
메모리 사용량         89MB           18MB           80% ↓
```

## 🎉 프로젝트 완료 기준

### ✅ 기술적 완료 기준
- [ ] 모든 핵심 기능 Modern Stack으로 이전 완료
- [ ] 성능 지표 목표 달성 (로딩 시간 < 1초)
- [ ] 모든 E2E 테스트 통과
- [ ] 보안 취약점 스캔 통과
- [ ] 접근성 테스트 통과 (WCAG 2.1 AA)

### ✅ 사용자 경험 완료 기준  
- [ ] 사용자 피드백 점수 8.5/10 이상
- [ ] 신규 버전 선택률 80% 이상
- [ ] 버그 리포트 기존 대비 50% 감소
- [ ] 모바일 사용률 증가

### ✅ 비즈니스 완료 기준
- [ ] SEO 트래픽 300% 증가
- [ ] 사용자 세션 시간 50% 증가  
- [ ] 이탈률 30% 감소
- [ ] 신규 사용자 가입률 증가

## 🚀 결론: 혁신적 웹 애플리케이션 완성!

**8주 후 Progressive Spreadsheet는:**

1. **⚡ 번개같이 빠른 로딩** (5초 → 0.5초)
2. **🌍 완벽한 SEO 최적화** (검색 노출 가능)
3. **📱 뛰어난 모바일 경험** (반응형 디자인)
4. **🤝 실시간 협업 기능** (WebSocket 기반)
5. **🛠️ 현대적 개발 환경** (타입 안전 + 빠른 개발)

**이 로드맵으로 Progressive Spreadsheet는 진정한 차세대 웹 애플리케이션으로 완전히 변모할 것입니다!** 🌟