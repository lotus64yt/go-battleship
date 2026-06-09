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

	for game.State == PlacingShips {
	drain:
		for {
			select {
			case packet, ok := <-incoming:
				if !ok {
					fmt.Println("\nOpponent disconnected!")
					time.Sleep(2 * time.Second)
					return
				}
				if packet.Type == "ready" {
					game.OtherPlayerState = Playing
				}
			default:
				break drain
			}
		}

		console.Clear()
		game.PrintPlacementBoard()

		input, err := console.ReadLine()
		if err != nil {
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

	console.Clear()
	game.PrintPlacementBoard()

	SendReadyPacket(conn)
	game.State = Playing

	fmt.Println("You are ready ! Waiting for your oponent...")

	for game.OtherPlayerState != Playing {
		packet, ok := <-incoming
		if !ok {
			fmt.Println("\nOpponent disconnected!")
			time.Sleep(2 * time.Second)
			return
		}

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
			packet, ok := <-incoming
			if !ok {
				fmt.Println("\nOpponent disconnected!")
				time.Sleep(2 * time.Second)
				return
			}

			if packet.Type == "attack" {
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
					console.Clear()
					game.PrintGameBoards()
					fmt.Println("You lost!")
					time.Sleep(3 * time.Second)
					return
				}
			}
		} else {
			input, err := console.ReadLine()
			if err != nil {
				return
			}

			game.LastInput = input
			game.LastError = ""

			hit, sunk, gameOver, err := SendAttackPacket(conn, input, incoming)
			if err != nil {
				if err.Error() == "connection closed" {
					fmt.Println("\nOpponent disconnected!")
					time.Sleep(2 * time.Second)
					return
				}
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
				console.Clear()
				game.PrintGameBoards()
				fmt.Println("You win !")
				time.Sleep(3 * time.Second)
				return
			} else if sunk {
				game.LastError = "You sunk a ship!"
			}
		}
	}
}
