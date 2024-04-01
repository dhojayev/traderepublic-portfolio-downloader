//go:build unix

package util

import (
	"syscall"
)

//nolint:gochecknoglobals
var stdinInt = syscall.Stdin
