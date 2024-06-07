package main

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

const (
	assetFilepath = "./assets/transactions.csv"
)

func main() {
	examples := fakes.TestCasesSupported
	logger := log.New()
	factory := transaction.NewCSVEntryFactory(logger)
	csvWriter := filesystem.NewCSVWriter(logger)

	if err := os.Remove(assetFilepath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}

	for _, example := range examples {
		csvEntry, err := factory.Make(example.Transaction)
		if err != nil {
			panic(err)
		}

		if err := csvWriter.Write(assetFilepath, csvEntry); err != nil {
			panic(fmt.Sprintf("could not save transaction to file: %s", err))
		}
	}
}
