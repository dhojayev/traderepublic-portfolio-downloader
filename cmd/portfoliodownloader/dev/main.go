package main

import (
	"github.com/alexflint/go-arg"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/cmd/portfoliodownloader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
)

const (
	responsesBaseDir = "responses"
)

func main() {
	var args Args

	arg.MustParse(&args)

	logger := log.New()
	logger.SetFormatter(&nested.Formatter{})

	logger.SetLevel(log.DebugLevel)

	if args.TraceMode {
		logger.SetLevel(log.TraceLevel)
		logger.SetReportCaller(true)
	}

	if err := internal.GetRuntimeTimezone(logger); err != nil {
		logger.Panic(err)
	}

	var application portfoliodownloader.App
	var err error

	if args.LocalMode {
		application, err = ProvideLocalApp(responsesBaseDir, logger)
	} else {
		application, err = ProvideRemoteApp(logger)
	}

	if err != nil {
		logger.Panic(err)
	}

	if err := application.Run(); err != nil {
		logger.Panic(err)
	}
}
