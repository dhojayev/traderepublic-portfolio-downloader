package api

import (
	"fmt"
	"os"
)

const (
	TokenNameSession TokenName = "session"
	TokenNameRefresh TokenName = "refresh"

	filePermissions = 0o600
)

type TokenName string

type Token struct {
	name  TokenName
	value string
}

func (t Token) Name() string {
	return string(t.name)
}

func (t Token) Value() string {
	return t.value
}

func NewToken(name TokenName, value string) Token {
	return Token{
		name:  name,
		value: value,
	}
}

func NewTokenFromFile(name TokenName) (Token, error) {
	token := Token{
		name: name,
	}
	filepath := "." + string(name)

	contentBytes, err := os.ReadFile(filepath)
	if err != nil {
		return token, fmt.Errorf("could not read token file '%s': %w", filepath, err)
	}

	token.value = string(contentBytes)

	return token, nil
}

func (t Token) WriteToFile() error {
	filepath := "." + string(t.name)
	if err := os.WriteFile(filepath, []byte(t.value), filePermissions); err != nil {
		return fmt.Errorf("could not write token file '%s': %w", filepath, err)
	}

	return nil
}
