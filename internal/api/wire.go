package api

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	NewClient,
	NewWSClient,
)
