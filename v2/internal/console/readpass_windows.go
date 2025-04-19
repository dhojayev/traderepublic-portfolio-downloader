//go:build windows

package console

import (
	"syscall"
)

func getStdin() int {
	return int(syscall.Stdin)
}
