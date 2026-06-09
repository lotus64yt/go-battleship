package app

import (
	"battleship/server"
	"fmt"
	"log"
	"strings"

	"battleship/game"
	"golang.org/x/net/websocket"
)

func CreateGame() error {
	s := &server.Server{}

	if err := s.Init(); err != nil {
		return err
	}

	fmt.Printf("Room created!\nShare this address with your friend so they can join:\n   -> %s:%d\n\n", s.Ip, s.Port)

	return s.Start()
}

func JoinGame(args ...string) error {
	if len(args) == 0 {
		fmt.Printf("Please provide ip:port or URL\n")
		return nil
	}

	urlStr := args[0]
	if !strings.HasPrefix(urlStr, "ws://") && !strings.HasPrefix(urlStr, "wss://") && !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "ws://" + urlStr
	} else if strings.HasPrefix(urlStr, "http://") {
		urlStr = strings.Replace(urlStr, "http://", "ws://", 1)
	} else if strings.HasPrefix(urlStr, "https://") {
		urlStr = strings.Replace(urlStr, "https://", "wss://", 1)
	}

	ws, err := websocket.Dial(urlStr, "", "http://localhost/")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	ws.PayloadType = websocket.TextFrame
	defer ws.Close()

	game.StartGame(ws, false)

	return nil
}
