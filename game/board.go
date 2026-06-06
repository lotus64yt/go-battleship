package game

import (
	utilsboard "battleship/utils/board"
	"fmt"
	"strings"
)

type Cell int

const (
	Empty Cell = iota
	Ship
	Hit
	Miss
)

type Board struct {
	Cells [9][9]Cell
}

type GameState int

type Game struct {
	Board Board

	ShipsPlaced map[int]int
	playerTurn  bool

	State     GameState
	LastInput string
	LastError string
}

func NewGame(creator bool) *Game {
	return &Game{
		ShipsPlaced: make(map[int]int),

		playerTurn: creator,
		State:      PlacingShips,
	}
}

func (b *Board) BoardString() []string {
	var result []string
	var line string

	pad := func(s string) string {
		return fmt.Sprintf("%-20s", s)
	}

	numbers := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for _, n := range numbers {
		line = fmt.Sprintf("%s%s ", line, n)
	}
	result = append(result, pad(line))

	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	for i, row := range b.Cells {
		line = fmt.Sprintf("%s ", alphabet[i])
		for _, cell := range row {
			switch cell {
			case Empty:
				line = fmt.Sprintf("%s. ", line)
			case Ship:
				line = fmt.Sprintf("%sS ", line)
			case Hit:
				line = fmt.Sprintf("%sX ", line)
			case Miss:
				line = fmt.Sprintf("%so ", line)
			default:
				line = fmt.Sprintf("%s? ", line)
			}
		}
		result = append(result, pad(line))
	}

	return result
}

func (g *Game) PlaceShip(content string) error {
	content = strings.ReplaceAll(content, ">", "")
	content = strings.ReplaceAll(content, "\n", "")
	args := strings.Split(content, " ")
	startNotation, endNotation := args[0], args[1]
	startX, startY := utilsboard.NotationToCoords(startNotation)
	endX, endY := utilsboard.NotationToCoords(endNotation)

	rows := len(g.Board.Cells)
	cols := len(g.Board.Cells[0])

	if startX < 0 || startX >= rows || startY < 0 || startY >= cols ||
		endX < 0 || endX >= rows || endY < 0 || endY >= cols {
		return fmt.Errorf("invalid placement")
	}

	if startX != endX && startY != endY {
		return fmt.Errorf("must be horizontal or vertical \n See '%s'", content)
	}

	if startX == endX {
		if startY > endY {
			startY, endY = endY, startY
		}
	} else {
		if startX > endX {
			startX, endX = endX, startX
		}
	}

	length := 0
	if startX == endX {
		length = endY - startY + 1
	} else {
		length = endX - startX + 1
	}

	if length < 2 || length > 5 {
		return fmt.Errorf("invalid ship size")
	}

	limits := map[int]int{
		2: 1,
		3: 2,
		4: 1,
		5: 1,
	}

	if g.ShipsPlaced[length] >= limits[length] {
		return fmt.Errorf("too many ships of size %d", length)
	}

	if startX == endX {
		for y := startY; y <= endY; y++ {
			if g.Board.Cells[startX][y] != Empty {
				return fmt.Errorf("collision detected")
			}
		}
		for y := startY; y <= endY; y++ {
			g.Board.Cells[startX][y] = Ship
		}
	} else {
		for x := startX; x <= endX; x++ {
			if g.Board.Cells[x][startY] != Empty {
				return fmt.Errorf("collision detected")
			}
		}
		for x := startX; x <= endX; x++ {
			g.Board.Cells[x][startY] = Ship
		}
	}

	g.ShipsPlaced[length]++

	return nil
}

func (g *Game) IsSinked(x, y int) bool {
	rows := len(g.Board.Cells)
	cols := len(g.Board.Cells[0])

	isOccupied := func(cx, cy int) bool {
		return cx >= 0 && cx < rows && cy >= 0 && cy < cols &&
			(g.Board.Cells[cx][cy] == Ship || g.Board.Cells[cx][cy] == Hit)
	}

	horizontal := (y > 0 && isOccupied(x, y-1)) || (y < cols-1 && isOccupied(x, y+1))
	if horizontal {
		for cy := y - 1; cy >= 0 && isOccupied(x, cy); cy-- {
			if g.Board.Cells[x][cy] == Ship {
				return false
			}
		}
		for cy := y + 1; cy < cols && isOccupied(x, cy); cy++ {
			if g.Board.Cells[x][cy] == Ship {
				return false
			}
		}
		return true
	}

	vertical := (x > 0 && isOccupied(x-1, y)) || (x < rows-1 && isOccupied(x+1, y))
	if vertical {
		for cx := x - 1; cx >= 0 && isOccupied(cx, y); cx-- {
			if g.Board.Cells[cx][y] == Ship {
				return false
			}
		}
		for cx := x + 1; cx < rows && isOccupied(cx, y); cx++ {
			if g.Board.Cells[cx][y] == Ship {
				return false
			}
		}
		return true
	}

	return true
}

func (g *Game) Fire(notation string) (bool, bool) {
	x, y := utilsboard.NotationToCoords(notation)

	cell := g.Board.Cells[x][y]
	g.playerTurn = !g.playerTurn

	if cell == Ship {
		g.Board.Cells[x][y] = Hit
		return true, g.IsSinked(x, y)
	}

	g.Board.Cells[x][y] = Miss
	return false, false
}

func (g *Game) PrintPlacementBoard() {
	fmt.Println("\n")
	check := func(v bool) string {
		if v {
			return "X"
		}
		return " "
	}

	board := g.Board.BoardString()

	fmt.Printf("Place your ships (Other player: Placing)\n")
	fmt.Println()

	remaingShips := []string{
		"",
		"Ships remaining:",
		fmt.Sprintf("[%s] Carrier    (5)", check(g.ShipsPlaced[5] >= 1)),
		fmt.Sprintf("[%s] BattleShip (4)", check(g.ShipsPlaced[4] >= 1)),
		fmt.Sprintf("[%s] Crusier    (3)", check(g.ShipsPlaced[3] >= 1)),
		fmt.Sprintf("[%s] Submarine  (3)", check(g.ShipsPlaced[3] >= 2)),
		fmt.Sprintf("[%s] Destroyer  (2)", check(g.ShipsPlaced[2] >= 1)),
	}

	for i, l := range board {
		var remain string = ""
		if i <= len(remaingShips)-1 {
			remain = remaingShips[i]
		}

		fmt.Printf("%s  %s\n", l, remain)
	}

	if g.LastError != "" {
		fmt.Printf("\033[31m> %s\033[0m\n", g.LastInput)
		fmt.Printf("\033[31m%s\033[0m\n", g.LastError)
	} else {
		fmt.Print("> ")
	}
}
