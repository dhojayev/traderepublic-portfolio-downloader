package main

import (
	"github.com/alexflint/go-arg"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/util"
)

func main() {
	var args Args

	arg.MustParse(&args)

	logger := log.New()
	logger.SetFormatter(&nested.Formatter{})

	switch {
	case args.DebugMode:
		logger.SetLevel(log.DebugLevel)
	case args.TraceMode:
		logger.SetLevel(log.TraceLevel)
		logger.SetReportCaller(true)
	}

	if err := internal.GetRuntimeTimezone(logger); err != nil {
		logger.Panic(err)
	}

	input, err := util.ReadPassword("pin")
	if err != nil {
		logger.Panic(err)
	}

	pin := auth.Pin(input)

	application, err := CreateNonWritingApp(args.PhoneNumber, pin, logger)
	if err != nil {
		logger.Panic(err)
	}

	if args.WriteResponseFiles {
		application, err = CreateWritingApp(args.PhoneNumber, pin, logger)
		if err != nil {
			logger.Panic(err)
		}
	}

	if err := application.Run(); err != nil {
		logger.Panic(err)
	}
}
