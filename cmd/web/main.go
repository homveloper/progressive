package main

import (
	"log"
	"net/http"

	"progressive/internal/handlers"
)

func main() {
	// 라우트 설정
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/table/create", handlers.TableCreateHandler)
	http.HandleFunc("/table/", handlers.TableEditorHandler)

	// 정적 파일 서빙
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("🚀 Progressive 앱이 http://localhost:8081 에서 실행 중입니다...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
