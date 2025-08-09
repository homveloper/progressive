package models

type TicTacToeGame struct {
	Board         [9]string
	CurrentPlayer string
	Winner        string
	IsGameOver    bool
	MoveCount     int
}

func NewTicTacToeGame() *TicTacToeGame {
	return &TicTacToeGame{
		Board:         [9]string{},
		CurrentPlayer: "X",
		Winner:        "",
		IsGameOver:    false,
		MoveCount:     0,
	}
}

func (g *TicTacToeGame) MakeMove(position int) bool {
	// 유효하지 않은 위치이거나 이미 게임이 끝났거나 해당 위치에 이미 표시가 있는 경우
	if position < 0 || position > 8 || g.IsGameOver || g.Board[position] != "" {
		return false
	}

	// 움직임 실행
	g.Board[position] = g.CurrentPlayer
	g.MoveCount++

	// 승리 조건 확인
	if g.checkWinner() {
		g.Winner = g.CurrentPlayer
		g.IsGameOver = true
	} else if g.MoveCount == 9 {
		// 무승부 (모든 칸이 찬 경우)
		g.IsGameOver = true
	} else {
		// 플레이어 교체
		if g.CurrentPlayer == "X" {
			g.CurrentPlayer = "O"
		} else {
			g.CurrentPlayer = "X"
		}
	}

	return true
}

func (g *TicTacToeGame) checkWinner() bool {
	// 승리 조건: 가로, 세로, 대각선
	winConditions := [][]int{
		// 가로
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		// 세로
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		// 대각선
		{0, 4, 8}, {2, 4, 6},
	}

	for _, condition := range winConditions {
		a, b, c := condition[0], condition[1], condition[2]
		if g.Board[a] != "" && g.Board[a] == g.Board[b] && g.Board[b] == g.Board[c] {
			return true
		}
	}

	return false
}

func (g *TicTacToeGame) Reset() {
	g.Board = [9]string{}
	g.CurrentPlayer = "X"
	g.Winner = ""
	g.IsGameOver = false
	g.MoveCount = 0
}