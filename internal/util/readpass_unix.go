//go:build unix

package util

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() ([]byte, error) {
	input, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, fmt.Errorf("could not read password from stdin: %w", err)
	}

	return input, nil
}
