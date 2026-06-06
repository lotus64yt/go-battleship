package server

import (
	"battleship/game"
	"battleship/utils/conn"
	"fmt"
	"net"
)

type Server struct {
	Ip   net.IP
	Port int

	listener net.Listener
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

	return nil
}

func (s *Server) Start() error {
	fmt.Println("Waiting for player...")

	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}

	fmt.Println("Player connected!")

	game.StartGame(conn)

	return nil
}
