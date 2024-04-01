package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	TokenNameSession TokenName = "session"
	TokenNameRefresh TokenName = "refresh"

	cookieNamePrefix = "tr_"
	cookieValueLen   = 514

	headerSetCookie = "Set-Cookie"
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

func NewTokenFromHeader(name TokenName, header http.Header) (Token, error) {
	token := Token{name: name}

	cookieHeader := header.Values(headerSetCookie)
	if len(cookieHeader) == 0 {
		return token, fmt.Errorf("could not find '%s' in header", headerSetCookie)
	}

	found := false

	for _, v := range cookieHeader {
		if !strings.Contains(v, string(name)) {
			continue
		}

		found = true
		cookieName := cookieNamePrefix + name
		startPos := len(cookieName) + 1
		token.value = v[startPos : strings.Index(v, ";")]
	}

	if !found {
		return token, fmt.Errorf("could not find '%s' token cookie in header", name)
	}

	return token, nil
}

func NewTokenFromFile(name TokenName) (Token, error) {
	token := Token{}
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
