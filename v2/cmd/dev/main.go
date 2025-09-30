package main

import (
	"log/slog"
	"os"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/cmd/dev/app"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/timelinedetails"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	debugMode := os.Getenv("DEBUG") == "true"
	logLevel := slog.LevelInfo

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: debugMode,
		Level:     logLevel,
	}))

	slog.SetDefault(log)

	eventBus := bus.New()
	faker := app.NewFaker(eventBus)
	tdHandler := timelinedetails.NewHandler()

	eventBus.Subscribe(bus.TopicTimelineDetailsV2, tdHandler.Handle)

	err := faker.Run()
	if err != nil {
		panic(err)
	}
}
