package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/file"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/timelinedetails"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/timelinetransactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/message/publisher"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

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

	wHandler := file.NewEventWriterHandler(writer.NewResponseWriter())
	eventBus := bus.New()

	eventBus.Subscribe(bus.TopicTimelineTransactions, wHandler.Handle)
	eventBus.Subscribe(bus.TopicTimelineDetailsV2, wHandler.Handle)

	wsclient := websocketclient.NewClient(publisher.NewPublisher(), ctx)

	messageClient := message.NewClient(eventBus, credentialsService, wsclient)
	ttHandler := timelinetransactions.NewHandler(eventBus, messageClient)
	tdHandler := timelinedetails.NewHandler()

	eventBus.Subscribe(bus.TopicTimelineTransactions, ttHandler.Handle)
	eventBus.Subscribe(bus.TopicTimelineDetailsV2, tdHandler.Handle)

	app := NewApp(auth.NewClient(console.NewInputHandler(), apiClient), credentialsService, messageClient, eventBus)

	err = app.Run()
	if err != nil {
		log.Error("Error running app", "error", err)
	}

	time.Sleep(time.Minute * 1)
}
