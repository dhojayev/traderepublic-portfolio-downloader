//go:build windows

package util

import (
	"syscall"
)

//nolint:gochecknoglobals
var stdinInt = int(syscall.Stdin)
