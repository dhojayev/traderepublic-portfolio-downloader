package util

import (
	"fmt"

	"golang.org/x/term"
)

func ReadPassword() ([]byte, error) {
	input, err := term.ReadPassword(stdinInt)
	if err != nil {
		return nil, fmt.Errorf("could not read password from stdin: %w", err)
	}

	return input, nil
}
