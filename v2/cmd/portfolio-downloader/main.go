package main

import (
	"log/slog"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
)

func main() {
	var args Args

	_ = godotenv.Load(".env")
	debugMode := os.Getenv("DEBUG") == "true"
	arg.MustParse(&args)

	logLevel := slog.LevelInfo
	if args.DebugMode || debugMode {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: debugMode,
		Level:     logLevel,
	}))

	logger.Info("Hello, world!", "debug", os.Getenv("DEBUG"))
	logger.Debug("Hello, world!")
}
