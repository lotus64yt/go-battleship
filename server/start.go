package server

import (
	"battleship/game"
	"battleship/utils/conn"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	Ip   net.IP
	Port int

	listener net.Listener
	connChan chan net.Conn
}

func (s *Server) Init() error {
	port, err := conn.GetFreePort()
	if err != nil {
		return err
	}

	s.Port = port
	ip, err := conn.GetLocalIp()
	if err != nil {
		return err
	}
	s.Ip = ip

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.listener = l
	s.connChan = make(chan net.Conn, 1)

	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		ws.PayloadType = websocket.TextFrame
		s.connChan <- ws
		var dummy []byte
		for {
			if _, err := ws.Read(dummy); err != nil {
				break
			}
		}
	}))

	go http.Serve(s.listener, mux)

	return nil
}

func (s *Server) Start() error {
	fmt.Println("Waiting for player...")

	conn := <-s.connChan

	fmt.Println("Player connected!")

	game.StartGame(conn, true)

	return nil
}
