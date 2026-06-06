package app

import "os"

func Exit(args ...string) error {
	os.Exit(1)
	return nil
}
