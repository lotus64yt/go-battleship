package game

import (
	"battleship/console"
	"net"
)

const (
	WaitingForPlayer GameState = iota
	PlacingShips
	Playing
	Finished
)

func StartGame(conn net.Conn, creator bool) {
	game := NewGame(creator)

	for game.State == PlacingShips {
		console.Clear()
		game.PrintPlacementBoard()

		input, _ := console.ReadLine()

		game.LastInput = input
		game.LastError = ""

		if err := game.PlaceShip(input); err != nil {
			game.LastError = err.Error()
		}
	}
}
