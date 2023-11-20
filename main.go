package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var board [3][3]rune
var currentPlayer rune

const PADDING = "   "

func checkWin() bool {
	return checkRows() || checkColumns() || checkDiagonals()
}

func checkDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == ' ' {
				return false
			}
		}
	}
	return true
}

func checkDiagonals() bool {
	return (board[0][0] != ' ' && board[0][0] == board[1][1] && board[1][1] == board[2][2]) ||
		(board[0][2] != ' ' && board[0][2] == board[1][1] && board[1][1] == board[2][0])
}

func checkColumns() bool {
	for i := 0; i < 3; i++ {
		if board[0][i] != ' ' && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
			return true
		}
	}
	return false
}

func checkRows() bool {
	for i := 0; i < 3; i++ {
		if board[i][0] != ' ' && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return true
		}
	}
	return false
}

func initializeBoard() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = ' '
		}
	}
}

func isValidMove(row, col int) bool {
	return row >= 0 && row < 3 && col >= 0 && col < 3 && board[row][col] == ' '
}

func clearConsole() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported platform. Cannot clear console.")
	}
}

func printBoard() {
	fmt.Println(PADDING + "   0   1   2")
	fmt.Println(PADDING + "  +---+---+---+")
	for i := 0; i < 3; i++ {
		fmt.Printf(PADDING+"%d |", i)
		for j := 0; j < 3; j++ {
			printCell(board[i][j])
			fmt.Print("|")
		}
		fmt.Println()
		fmt.Println(PADDING + "  +---+---+---+")
	}
	fmt.Println()
}

func printCell(cellValue rune) {
	switch cellValue {
	case 'X':
		fmt.Print(" X ")
	case 'O':
		fmt.Print(" O ")
	default:
		fmt.Print("   ")
	}
}
