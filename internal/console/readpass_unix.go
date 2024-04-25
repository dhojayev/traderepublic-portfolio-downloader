//go:build unix

package console

import (
	"syscall"
)

//nolint:gochecknoglobals
var stdinInt = syscall.Stdin
