package database

import "github.com/google/wire"

var SqliteInMemorySet = wire.NewSet(
	NewSQLiteInMemory,
)

var SqliteOnFilesystemSet = wire.NewSet(
	NewSQLiteOnFS,
)
