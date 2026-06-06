package console

import (
	"bufio"
	"os"
)

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}
