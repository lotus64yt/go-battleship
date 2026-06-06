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

type PlacedShip struct {
	StartX, StartY int
	EndX, EndY     int
}

type GameState int

type Game struct {
	Board      Board
	EnemyBoard Board

	ShipsPlaced map[int]int
	Ships       []PlacedShip
	PlayerTurn  bool

	State            GameState
	OtherPlayerState GameState
	LastInput        string
	LastError        string
}

func NewGame(creator bool) *Game {
	return &Game{
		ShipsPlaced: make(map[int]int),

		PlayerTurn:       creator,
		State:            PlacingShips,
		OtherPlayerState: PlacingShips,
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
				line = fmt.Sprintf("%s\033[34m.\033[0m ", line)
			case Ship:
				line = fmt.Sprintf("%s\033[90mS\033[0m ", line)
			case Hit:
				line = fmt.Sprintf("%s\033[31mX\033[0m ", line)
			case Miss:
				line = fmt.Sprintf("%s\033[37mo\033[0m ", line)
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
	content = strings.TrimSpace(content)
	args := strings.Fields(content)
	if len(args) < 2 {
		return fmt.Errorf("invalid placement")
	}
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
	g.Ships = append(g.Ships, PlacedShip{startX, startY, endX, endY})

	return nil
}

func (g *Game) IsSinked(x, y int) bool {
	for _, s := range g.Ships {
		minX, maxX := s.StartX, s.EndX
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		minY, maxY := s.StartY, s.EndY
		if minY > maxY {
			minY, maxY = maxY, minY
		}

		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			for cx := minX; cx <= maxX; cx++ {
				for cy := minY; cy <= maxY; cy++ {
					if g.Board.Cells[cx][cy] == Ship {
						return false
					}
				}
			}
			return true
		}
	}
	return false
}

func (g *Game) IsGameOver() bool {
	for i := 0; i < len(g.Board.Cells); i++ {
		for j := 0; j < len(g.Board.Cells[0]); j++ {
			if g.Board.Cells[i][j] == Ship {
				return false
			}
		}
	}
	return true
}

func (g *Game) Fire(notation string) (bool, bool, bool) {
	x, y := utilsboard.NotationToCoords(notation)

	cell := g.Board.Cells[x][y]
	g.PlayerTurn = !g.PlayerTurn

	if cell == Ship {
		g.Board.Cells[x][y] = Hit
		sunk := g.IsSinked(x, y)

		if sunk && g.IsGameOver() {
			g.State = Finished
			g.OtherPlayerState = Finished
			return true, sunk, true
		}

		return true, sunk, false
	}

	g.Board.Cells[x][y] = Miss
	return false, false, false
}

func (g *Game) PrintPlacementBoard() {
	fmt.Println()
	fmt.Println()
	check := func(v bool) string {
		if v {
			return "X"
		}
		return " "
	}

	board := g.Board.BoardString()
	otherPlayerState := "Placing"
	if g.OtherPlayerState == Playing {
		otherPlayerState = "Ready!"
	}

	fmt.Printf("Place your ships (Other player: %s)\n", otherPlayerState)
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

func (g *Game) PrintGameBoards() {
	fmt.Println()
	fmt.Println()

	turn := "opponent"
	if g.PlayerTurn {
		turn = "your"
	}

	fmt.Printf("Its %s turn!\n", turn)

	fmt.Printf("%-20s     %-20s\n\n", fmt.Sprintf("%15s", "Your Board"), fmt.Sprintf("%15s", "Opponent Board"))

	enemy := g.EnemyBoard.BoardString()

	for i, l := range g.Board.BoardString() {
		fmt.Printf("%s   %s\n", l, enemy[i])
	}

	if g.LastError != "" {
		fmt.Printf("\033[31m> %s\033[0m\n", g.LastInput)
		fmt.Printf("\033[31m%s\033[0m\n", g.LastError)
	} else {
		fmt.Print("> ")
	}
}

func (g *Game) IsAllShipsPlaced() bool {
	return g.ShipsPlaced[5] >= 1 && g.ShipsPlaced[4] >= 1 && g.ShipsPlaced[3] >= 1 && g.ShipsPlaced[3] >= 2 && g.ShipsPlaced[2] >= 1
}
