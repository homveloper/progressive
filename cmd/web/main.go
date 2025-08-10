package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"progressive/internal/handlers"
	"progressive/internal/infrastructure"
	"progressive/internal/middleware"
)

func main() {
	// 임베디드 PostgreSQL 인스턴스 생성 (자동 포트 발견 기능 사용)
	embeddedDB, err := infrastructure.NewEmbeddedDB(
		infrastructure.WithAutoPortDiscovery(10), // 최대 10개 포트 시도
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer embeddedDB.Close()

	// 사용된 설정 정보 출력
	config := embeddedDB.GetConfig()
	log.Printf("📊 PostgreSQL running on %s:%d", config.Host, config.Port)

	// sqlx.DB 인스턴스
	db := embeddedDB.DB

	// 마이그레이션 실행
	if err := infrastructure.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 핸들러에 DB 의존성 주입 (템플릿 초기화는 핸들러 생성 시 자동으로 실행됨)
	h := handlers.NewHandlers(db)

	// 라우트 설정을 위한 ServeMux 생성
	mux := http.NewServeMux()

	// 페이지 라우트 설정 (GET 요청으로 HTML 페이지 렌더링)
	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("/dashboard", h.DashboardHandler)
	mux.HandleFunc("/table/create", h.TableCreatePageHandler)
	mux.HandleFunc("/table/", h.TableEditorPageHandler)
	mux.HandleFunc("/fakeit", h.FakeitPageHandler)

	// API 라우트 설정 (JSON 데이터 처리)
	mux.HandleFunc("/api/templates", h.TemplatesAPIHandler)
	mux.HandleFunc("/api/tables", h.TablesAPIHandler)
	mux.HandleFunc("/api/table/", func(w http.ResponseWriter, r *http.Request) {
		// Route to specific table API handlers based on URL pattern
		path := r.URL.Path
		if strings.Contains(path, "/export") {
			h.Table.API.ExportHandler(w, r)
		} else if strings.Contains(path, "/import") {
			h.Table.API.ImportHandler(w, r)
		} else if strings.Contains(path, "/record") {
			h.Table.API.RecordHandler(w, r)
		} else {
			h.Table.API.DataHandler(w, r)
		}
	})
	mux.HandleFunc("/api/table/create", h.TableCreateAPIHandler)
	mux.HandleFunc("/api/fakeit/generate", h.FakeitGenerateAPIHandler)

	// 정적 파일 서빙
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// 미들웨어 체인 적용 (에러 핸들링 -> 로깅 순서)
	errorHandledMux := middleware.ErrorHandlingMiddleware(mux)
	loggedMux := middleware.LoggingMiddleware(errorHandledMux)

	// Graceful shutdown 설정
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("🛑 Shutting down server...")
		embeddedDB.Close()
		os.Exit(0)
	}()

	log.Println("🚀 Progressive 앱이 http://localhost:8081 에서 실행 중입니다...")
	log.Fatal(http.ListenAndServe(":8081", loggedMux))
}
