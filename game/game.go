package game

import "net"

type GameState int

const (
	WaitingForPlayer GameState = iota
	PlacingShips
	Playing
	Finished
)

func StartGame(conn net.Conn) {

}
