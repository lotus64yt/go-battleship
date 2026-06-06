package game

import (
	"battleship/console"
	utilsboard "battleship/utils/board"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	WaitingForPlayer GameState = iota
	PlacingShips
	Playing
	Finished
)

func StartGame(conn net.Conn, creator bool) {
	game := NewGame(creator)
	incoming := make(chan Packet)
	decoder := json.NewDecoder(conn)

	go func() {
		for {
			var packet Packet
			if err := decoder.Decode(&packet); err != nil {
				close(incoming)
				return
			}

			incoming <- packet
		}
	}()

	inputs := make(chan string)

	go func() {
		for {
			input, err := console.ReadLine()
			if err != nil {
				close(inputs)
				return
			}

			inputs <- input
		}
	}()

	for game.State == PlacingShips {
		console.Clear()
		game.PrintPlacementBoard()

		select {
		case packet, ok := <-incoming:
			if !ok {
				return
			}

			switch packet.Type {
			case "ready":
				game.OtherPlayerState = Playing
			}

		case input, ok := <-inputs:
			if !ok {
				return
			}

			game.LastInput = input
			game.LastError = ""

			if err := game.PlaceShip(input); err != nil {
				game.LastError = err.Error()
			}

			if game.IsAllShipsPlaced() {
				game.State = Playing
			}
		}
	}

	SendReadyPacket(conn)
	game.State = Playing

	fmt.Println("You are ready ! Waiting for your oponent...")

	for game.OtherPlayerState != Playing {
		packet := <-incoming

		switch packet.Type {
		case "ready":
			game.OtherPlayerState = Playing
		}
	}

	fmt.Println("Both players are ready!")

	for game.State == Playing && game.OtherPlayerState == Playing {
		console.Clear()
		game.PrintGameBoards()
		time.Sleep(time.Second * 1)

		if !game.PlayerTurn {
			fmt.Println("Waiting for opponent to fire...")
			select {
			case packet, ok := <-incoming:
				if !ok {
					return
				}

				switch packet.Type {
				case "attack":
					var data AttackPacket
					_ = json.Unmarshal(packet.Data, &data)

					hit, sunk, gameOver := game.Fire(data.Notation)

					resp := AttackResultPacket{
						Hit:      hit,
						Sunk:     sunk,
						GameOver: gameOver,
					}
					respData, _ := json.Marshal(resp)

					json.NewEncoder(conn).Encode(Packet{
						Type: "attack_result",
						Data: respData,
					})

					if gameOver {
						game.State = Finished
						fmt.Println("You lost!")
						return
					}
				}
			}
		} else {
			select {
			case input, ok := <-inputs:
				if !ok {
					return
				}

				game.LastInput = input
				game.LastError = ""

				hit, sunk, gameOver, err := SendAttackPacket(conn, input, incoming)
				if err != nil {
					game.LastError = err.Error()
					continue
				}

				game.PlayerTurn = false

				x, y := utilsboard.NotationToCoords(input)
				if hit {
					game.EnemyBoard.Cells[x][y] = Hit
				} else {
					game.EnemyBoard.Cells[x][y] = Miss
				}

				if gameOver {
					game.State = Finished
					fmt.Println("You win !")
					return
				} else if sunk {
					game.LastError = "You sunk a ship!"
				}
			}
		}
	}
}
