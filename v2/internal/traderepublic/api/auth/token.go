package auth

// Token represents authentication tokens.
type Token struct {
	session string
	refresh string
}

func NewToken() Token {
	return Token{}
}

// NewTokenWithValues creates a new token with the given session and refresh values.
func NewTokenWithValues(sessionValue, refreshValue string) Token {
	return Token{
		session: sessionValue,
		refresh: refreshValue,
	}
}

// SessionToken returns the session token value.
func (t Token) SessionToken() string {
	return t.session
}

// RefreshToken returns the refresh token value.
func (t Token) RefreshToken() string {
	return t.refresh
}
