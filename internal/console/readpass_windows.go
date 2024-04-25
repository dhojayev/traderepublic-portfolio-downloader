//go:build windows

package console

import (
	"syscall"
)

//nolint:gochecknoglobals
var stdinInt = int(syscall.Stdin)
