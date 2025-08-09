package handlers

import (
	"net/http"
	"progressive/internal/models"
	"progressive/internal/pages"
	"strconv"
	"strings"
)

type TicTacToeHandler struct {
	game *models.TicTacToeGame
}

func NewTicTacToeHandler() *TicTacToeHandler {
	return &TicTacToeHandler{
		game: models.NewTicTacToeGame(),
	}
}

func (h *TicTacToeHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := pages.TicTacToeData{
		Board:         h.game.Board,
		CurrentPlayer: h.game.CurrentPlayer,
		Winner:        h.game.Winner,
		IsGameOver:    h.game.IsGameOver,
	}

	component := pages.TicTacToe(data)
	component.Render(r.Context(), w)
}

func (h *TicTacToeHandler) MoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 위치 추출 /tictactoe/move/0, /tictactoe/move/1, ...
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid position", http.StatusBadRequest)
		return
	}

	position, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid position", http.StatusBadRequest)
		return
	}

	// 게임에서 움직임 실행
	if !h.game.MakeMove(position) {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	// 게임 보드만 반환 (HTMX가 #game-container를 업데이트)
	h.renderGameContainer(w, r)
}

func (h *TicTacToeHandler) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.game.Reset()
	h.renderGameContainer(w, r)
}

func (h *TicTacToeHandler) renderGameContainer(w http.ResponseWriter, r *http.Request) {
	data := pages.TicTacToeData{
		Board:         h.game.Board,
		CurrentPlayer: h.game.CurrentPlayer,
		Winner:        h.game.Winner,
		IsGameOver:    h.game.IsGameOver,
	}

	// 게임 컨테이너 부분만 렌더링
	w.Header().Set("Content-Type", "text/html")

	// 게임 상태 표시
	if data.IsGameOver {
		if data.Winner != "" {
			w.Write([]byte(`<div class="text-center mb-6">
				<div class="text-2xl font-bold text-green-600">🎉 ` + data.Winner + `이(가) 승리했습니다!</div>
				<button hx-post="/tictactoe/reset" hx-target="#game-container" 
					class="mt-4 px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
					다시 시작
				</button>
			</div>`))
		} else {
			w.Write([]byte(`<div class="text-center mb-6">
				<div class="text-2xl font-bold text-yellow-600">🤝 무승부입니다!</div>
				<button hx-post="/tictactoe/reset" hx-target="#game-container" 
					class="mt-4 px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
					다시 시작
				</button>
			</div>`))
		}
	} else {
		w.Write([]byte(`<div class="text-center mb-6">
			<div class="text-xl font-semibold text-blue-600">현재 차례: ` + data.CurrentPlayer + `</div>
		</div>`))
	}

	// 게임 보드
	w.Write([]byte(`<div class="grid grid-cols-3 gap-2 max-w-xs mx-auto">`))

	for i := 0; i < 9; i++ {
		cellValue := data.Board[i]
		cellClass := "w-20 h-20 text-3xl font-bold rounded-lg border-2 transition-all"

		if cellValue == "" && !data.IsGameOver {
			cellClass += " bg-gray-100 border-gray-300 hover:bg-gray-200 hover:shadow-md cursor-pointer"
		} else if cellValue == "X" {
			cellClass += " bg-blue-100 border-blue-300 text-blue-600"
		} else if cellValue == "O" {
			cellClass += " bg-red-100 border-red-300 text-red-600"
		}

		if cellValue != "" || data.IsGameOver {
			cellClass += " cursor-not-allowed opacity-50"
		}

		w.Write([]byte(`<button `))
		if cellValue == "" && !data.IsGameOver {
			w.Write([]byte(`hx-post="/tictactoe/move/` + strconv.Itoa(i) + `" hx-target="#game-container" `))
		}
		w.Write([]byte(`class="` + cellClass + `"`))
		if cellValue != "" || data.IsGameOver {
			w.Write([]byte(` disabled`))
		}
		w.Write([]byte(`>` + cellValue + `</button>`))
	}

	w.Write([]byte(`</div>`))
}
