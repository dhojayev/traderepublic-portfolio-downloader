package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/file"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/instrument"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/message"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/timelinedetails"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/timelinetransactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/writer"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
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

	wHandler := file.NewRawResponseHandler(writer.NewResponseWriter())
	eventBus := bus.New()

	eventBus.Subscribe(bus.TopicTimelineTransactionsReceived, wHandler.Handle)
	eventBus.Subscribe(bus.TopicTimelineDetailsV2Received, wHandler.Handle)
	eventBus.Subscribe(bus.TopicInstrumentReceived, wHandler.Handle)

	wsclient := traderepublic.NewWSClient(traderepublic.NewPublisher(), ctx)

	cache := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	msgClient := message.NewClient(eventBus, credentialsService, wsclient)
	ttHandler := timelinetransactions.NewHandler(eventBus, msgClient)
	tdHandler := timelinedetails.NewHandler(eventBus)
	instrHandler := instrument.NewHandler(msgClient, cache)

	mapper := transaction.NewDataMapper(cache)
	resolver := transaction.NewTypeResolver()
	trnHandler := transaction.NewHandler(resolver, mapper, eventBus)
	csvWriter := file.NewCSVWriter()
	csvHandler := file.NewCSVHandler(internal.CSVFilename, csvWriter)

	eventBus.Subscribe(bus.TopicTimelineTransactionsReceived, ttHandler.Handle)
	eventBus.Subscribe(bus.TopicTimelineDetailsV2Received, tdHandler.Handle)
	eventBus.Subscribe(bus.TopicInstrumentFetch, instrHandler.HandleFetch)
	eventBus.Subscribe(bus.TopicInstrumentReceived, instrHandler.HandleReceived)
	eventBus.Subscribe(bus.TopicTimelineDetailsV2Received, trnHandler.Handle)
	eventBus.Subscribe(bus.TopicModelReady, csvHandler.Handle)

	app := NewApp(auth.NewClient(console.NewInputHandler(), apiClient), credentialsService, msgClient, eventBus)

	err = app.Run()
	if err != nil {
		log.Error("Error running app", "error", err)
	}

	time.Sleep(time.Minute * 1)
}
