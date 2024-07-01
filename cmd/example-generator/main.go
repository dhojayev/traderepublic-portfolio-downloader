package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

const (
	assetFilepath = "./assets/transactions.csv"
)

func main() {
	examples := fakes.TransactionTestCasesSupported
	keys := make([]string, 0, len(examples))

	for k := range examples {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	logger := log.New()
	factory := transaction.NewCSVEntryFactory(logger)
	csvWriter := filesystem.NewCSVWriter(logger)

	if err := os.Remove(assetFilepath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Fatal(err)
		}
	}

	for _, key := range keys {
		csvEntry, err := factory.Make(examples[key].Transaction)
		if err != nil {
			logger.Fatal(err)
		}

		if err := csvWriter.Write(assetFilepath, csvEntry); err != nil {
			logger.Fatal(fmt.Sprintf("could not save transaction to file: %s", err))
		}
	}
}
