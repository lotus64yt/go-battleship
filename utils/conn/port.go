package conn

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const portCacheFile = ".battleship_port"

func isPortFree(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	l.Close()
	return true
}

func GetFreePort() (int, error) {
	data, err := os.ReadFile(portCacheFile)
	if err == nil {
		portStr := strings.TrimSpace(string(data))
		port, err := strconv.Atoi(portStr)
		if err == nil && port > 0 && port <= 65535 {
			if isPortFree(port) {
				return port, nil
			}
		}
	}

	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			port := l.Addr().(*net.TCPAddr).Port
			l.Close()
			os.WriteFile(portCacheFile, []byte(strconv.Itoa(port)), 0644)
			return port, nil
		}
	}
	return 0, err
}
