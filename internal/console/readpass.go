package console

import (
	"fmt"

	"golang.org/x/term"
)

func ReadPassword(name string) ([]byte, error) {
	fmt.Printf("Enter %s: \n", name)

	input, err := term.ReadPassword(stdinInt)
	if err != nil {
		return nil, fmt.Errorf("could not read %s from stdin: %w", name, err)
	}

	return input, nil
}
