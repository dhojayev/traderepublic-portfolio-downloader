//go:build unix

package console

import (
	"syscall"
)

func getStdin() int {
	return syscall.Stdin
}
