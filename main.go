package main

import (
	"battleship/utils/conn"
	"fmt"
)

func main() {
	ip, err := conn.GetLocalIp()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	port, err := conn.GetFreePort()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s:%d\n", ip, port)
}
