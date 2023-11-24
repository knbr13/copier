package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gookit/color"
)

var board [3][3]rune
var currentPlayer rune
var winning_cells [3][2]int

const PADDING = "   "

func main() {
	ticker := time.NewTicker(time.Second)
	delay := 3

	for delay > 0 {
		clearConsole()
		fmt.Println(`
	==================================================================
	 _______ _____ _____   _______       _____   _______ ____  ______ 
	|__   __|_   _/ ____| |__   __|/\   / ____| |__   __/ __ \|  ____|
	   | |    | || |         | |  /  \ | |         | | | |  | | |__   
	   | |    | || |         | | / /\ \| |         | | | |  | |  __|  
	   | |   _| || |____     | |/ ____ \ |____     | | | |__| | |____ 
	   |_|  |_____\_____|    |_/_/    \_\_____|    |_|  \____/|______|

	==================================================================
	`)

		fmt.Printf("\n\t\t\tthe game will start in %d seconds\t\t\t\n", delay)
		delay--
		<-ticker.C
	}

	clearConsole()
	initializeBoard()
	currentPlayer = 'X'

	for {
		clearConsole()
		printBoard()

		var row, col int
		fmt.Printf("%sPlayer %c's turn. \n%sEnter your move (row column): ", PADDING, currentPlayer, PADDING)
		fmt.Scanln(&row, &col)

		if isValidMove(row-1, col-1) {
			board[row-1][col-1] = currentPlayer
			clearConsole()
			if checkWin() {
				printBoard()
				fmt.Printf(PADDING+"  player %c wins!\n", currentPlayer)
				if restartPrompt() {
					resetGame()
					continue
				}
				break
			} else if checkDraw() {
				printBoard()
				fmt.Println(PADDING + "  It's a draw!")
				if restartPrompt() {
					resetGame()
					continue
				}
				break
			}
			switchPlayer()
		}
	}
}

func resetGame() {
	board = [3][3]rune{}
	currentPlayer = 'X'
	initializeBoard()
}

func restartPrompt() bool {
	fmt.Println(PADDING + "y for yes, other for no")
	fmt.Print(PADDING + "do you want to restart game? ")

	reader := bufio.NewReader(os.Stdin)
	r, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal("tic-tac-toe: ", err)
	}
	if r == 'y' || r == 'Y' {
		return true
	}
	return false
}

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
	if board[0][0] != ' ' && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		winning_cells[0] = [2]int{0, 0}
		winning_cells[1] = [2]int{1, 1}
		winning_cells[2] = [2]int{2, 2}
		return true
	}
	if board[0][2] != ' ' && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		winning_cells[0] = [2]int{0, 2}
		winning_cells[1] = [2]int{1, 1}
		winning_cells[2] = [2]int{2, 0}
		return true
	}
	return false
}

func checkColumns() bool {
	for i := 0; i < 3; i++ {
		if board[0][i] != ' ' && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
			winning_cells[0] = [2]int{0, i}
			winning_cells[1] = [2]int{1, i}
			winning_cells[2] = [2]int{2, i}
			return true
		}
	}
	return false
}

func checkRows() bool {
	for i := 0; i < 3; i++ {
		if board[i][0] != ' ' && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			winning_cells[0] = [2]int{i, 0}
			winning_cells[1] = [2]int{i, 1}
			winning_cells[2] = [2]int{i, 2}
			return true
		}
	}
	return false
}

func initializeBoard() {
	for i := 0; i < len(winning_cells); i++ {
		for j := 0; j < len(winning_cells[0]); j++ {
			winning_cells[i][j] = -1
		}
	}
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
	fmt.Println()
	fmt.Println(PADDING + "    1   2   3")
	fmt.Println(PADDING + "  +---+---+---+")
	for i := 0; i < 3; i++ {
		fmt.Printf(PADDING+"%d |", i+1)
		for j := 0; j < 3; j++ {
			printCell(i, j)
			fmt.Print("|")
		}
		fmt.Println()
		fmt.Println(PADDING + "  +---+---+---+")
	}
	fmt.Println()
}

func printCell(row, col int) {

	for _, v := range winning_cells {
		if v[0] == row && v[1] == col {
			color.Greenf(" %c ", board[row][col])
			return
		}
	}

	switch board[row][col] {
	case 'X':
		fmt.Print(" X ")
	case 'O':
		fmt.Print(" O ")
	default:
		fmt.Print("   ")
	}
}

func switchPlayer() {
	if currentPlayer == 'X' {
		currentPlayer = 'O'
	} else {
		currentPlayer = 'X'
	}
}
