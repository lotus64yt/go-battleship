package utilsboard

import (
	"battleship/utils/array"
	"fmt"
	"strings"
)

var numbers = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

func NotationToCoords(notation string) (int, int) {
	if ok, _ := IsValidNotation(notation); !ok {
		return 0, 0
	}

	args := strings.Split(notation, "")
	line, col := args[0], args[1]

	lineIdx := array.IndexOf(alphabet, line)
	colIdx := array.IndexOf(numbers, col)

	return lineIdx, colIdx
}

func CoordsToNotation(x, y int) string {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= len(alphabet) {
		x = len(alphabet) - 1
	}
	if y >= len(numbers) {
		y = len(numbers) - 1
	}

	return fmt.Sprintf("%s%s", alphabet[x], numbers[y])
}

func IsValidNotation(notation string) (bool, error) {
	if len(notation) != 2 {
		return false, fmt.Errorf("notation must have exactly 2 characters")
	}

	args := strings.Split(notation, "")
	line, col := args[0], args[1]

	if !array.Contain(alphabet, line) {
		return false, fmt.Errorf("invalid line: must be A-I")
	}

	if !array.Contain(numbers, col) {
		return false, fmt.Errorf("invalid column: must be 1-9")
	}

	return true, nil
}
