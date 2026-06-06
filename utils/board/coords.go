package utilsboard

import (
	"battleship/utils/array"
	"fmt"
	"strings"
)

var numbers = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

func NotationToCoords(notation string) (int, int) {
	if len(notation) > 2 {
		return 0, 0
	}
	args := strings.Split(notation, "")
	line, col := args[0], args[1]

	if !array.Contain(alphabet, line) {
		line = alphabet[len(alphabet)-1]
	}

	if !array.Contain(numbers, col) {
		col = numbers[len(numbers)-1]
	}

	lineIdx := array.IndexOf(alphabet, line)
	colIdx := array.IndexOf(numbers, col)

	return lineIdx, colIdx
}

func CoordsToNotation(x, y int) string {
	if x > len(numbers) {
		x = len(numbers)
	}
	if y > len(alphabet) {
		y = len(alphabet)
	}

	return fmt.Sprintf("%s%s", alphabet[x], numbers[y])
}
