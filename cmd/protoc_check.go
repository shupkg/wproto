package cmd

import (
	"os"
)

func InProtoc() bool {
	stat, err := os.Stdin.Stat()
	return len(os.Args) == 1 && err == nil && (stat.Mode()&os.ModeCharDevice) == 0
}
