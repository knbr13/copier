package main

import (
	"fmt"
)

var board [3][3]rune
var currentPlayer rune

const PADDING = "   "

func printCell(cellValue string) {
	switch cellValue {
	case "X":
		fmt.Print(" X ")
	case "O":
		fmt.Print(" O ")
	default:
		fmt.Print("   ")
	}
}
