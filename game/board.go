package game

import (
	utilsboard "battleship/utils/board"
	"fmt"
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

type Game struct {
	Board Board

	ShipsPlaced map[int]int
	playerTurn  bool
}

func NewGame() *Game {
	return &Game{
		ShipsPlaced: make(map[int]int),
	}
}

func (b *Board) PrintBoard() {
	numbers := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for _, n := range numbers {
		fmt.Printf("%s ", n)
	}
	fmt.Println()

	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	for i, row := range b.Cells {
		fmt.Printf("%s ", alphabet[i])
		for _, cell := range row {
			switch cell {
			case Empty:
				fmt.Print(". ")
			case Ship:
				fmt.Print("S ")
			case Hit:
				fmt.Print("X ")
			case Miss:
				fmt.Print("o ")
			default:
				fmt.Print("? ")
			}
		}
		fmt.Println()
	}
}

func (g *Game) PlaceShip(startNotation, endNotation string) error {
	startX, startY := utilsboard.NotationToCoords(startNotation)
	endX, endY := utilsboard.NotationToCoords(endNotation)

	rows := len(g.Board.Cells)
	cols := len(g.Board.Cells[0])

	if startX < 0 || startX >= rows || startY < 0 || startY >= cols ||
		endX < 0 || endX >= rows || endY < 0 || endY >= cols {
		return fmt.Errorf("invalid placement")
	}

	if startX != endX && startY != endY {
		return fmt.Errorf("must be horizontal or vertical")
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
