package main

import (
	"log/slog"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message/publisher"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
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

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: debugMode,
		Level:     logLevel,
	}))

	slog.SetDefault(log)

	credentialsService := auth.NewFileCredentialsService("")

	apiClient, err := api.NewClient()
	if err != nil {
		log.Error("Error creating API client", "error", err)

		return
	}

	messageClient := message.NewClient(credentialsService, websocketclient.NewClient(publisher.NewPublisher()))
	app := NewApp(auth.NewClient(console.NewInputHandler(), apiClient), credentialsService, messageClient)

	err = app.Run()
	if err != nil {
		log.Error("Error running app", "error", err)
	}
}
