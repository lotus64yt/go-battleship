package server

import (
	"battleship/game"
	"battleship/utils/conn"
	"fmt"
	"io"
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

func (s *Server) ListenTCP(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn) {
			buf := make([]byte, 1024)
			n, _ := c.Read(buf)
			fmt.Println(string(buf[:n]))
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}
