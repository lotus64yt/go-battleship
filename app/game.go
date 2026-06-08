package app

import (
	"battleship/server"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"battleship/game"
)

func CreateGame() error {
	s := &server.Server{}

	if err := s.Init(); err != nil {
		return err
	}

	fmt.Printf("Room created, connect with :\n   > join %s:%d\n", s.Ip, s.Port)

	return s.Start()
}

func JoinGame(args ...string) error {
	if len(args) == 0 {
		fmt.Printf("Please provide ip:port\n")
		return nil
	}

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", args[0])
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	game.StartGame(conn, false)

	return nil
}
