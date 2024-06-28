package console

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewAuthService,

	wire.Bind(new(AuthServiceInterface), new(*AuthService)),
)
