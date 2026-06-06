package console

import (
	"battleship/server"
	"fmt"
)

func CreateGame(args ...string) error {
	s := &server.Server{}

	if err := s.Init(); err != nil {
		return err
	}

	fmt.Printf("Room created, connect at :\n   > %s:%d\n", s.Ip, s.Port)

	return s.Start()
}
