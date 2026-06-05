package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	status bool
}

func (s *Server) Start(port int) error {
	return s.ListenTCP(port)
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
