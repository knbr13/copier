package main

import (
	"fmt"
)

var board [3][3]rune
var currentPlayer rune

const PADDING = "   "

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
