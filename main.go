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
