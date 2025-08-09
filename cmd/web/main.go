package main

import (
	"log"
	"net/http"
	"progressive/internal/handlers"
)

func main() {
	// 핸들러 초기화
	tictactoeHandler := handlers.NewTicTacToeHandler()

	// 라우트 설정
	http.HandleFunc("/", tictactoeHandler.IndexHandler)
	http.HandleFunc("/tictactoe", tictactoeHandler.IndexHandler)
	http.HandleFunc("/tictactoe/move/", tictactoeHandler.MoveHandler)
	http.HandleFunc("/tictactoe/reset", tictactoeHandler.ResetHandler)

	// 정적 파일 서빙
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("🚀 Tic Tac Toe 앱이 http://localhost:8081 에서 실행 중입니다...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
