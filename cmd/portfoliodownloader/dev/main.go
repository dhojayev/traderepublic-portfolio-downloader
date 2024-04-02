package main

import (
	"github.com/alexflint/go-arg"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/util"
)

const (
	responsesBaseDir string = "responses"
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

	application, err := CreateLocalApp(responsesBaseDir, logger)
	if err != nil {
		logger.Panic(err)
	}

	if !args.LocalMode {
		input, err := util.ReadPassword("pin")
		if err != nil {
			logger.Panic(err)
		}

		pin := auth.Pin(input)

		application, err = CreateRemoteApp(args.PhoneNumber, pin, logger)
		if err != nil {
			logger.Panic(err)
		}
	}

	if err := application.Run(); err != nil {
		logger.Panic(err)
	}
}
