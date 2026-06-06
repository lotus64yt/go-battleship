package console

import (
	"bufio"
	"os"
	"strings"
)

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(str, "\r\n"), nil
}
