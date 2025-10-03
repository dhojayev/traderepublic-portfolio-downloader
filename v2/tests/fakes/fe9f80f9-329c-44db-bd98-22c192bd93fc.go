package fakes

import "github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"

func Modelfe9f80f9329c44dbbd9822c192bd93fc() (transaction.Model, error) {
	var model transaction.Model

	timestamp, err := transaction.ParseTimestamp("2025-01-02T14:52:18.686+0000")
	if err != nil {
		return model, err
	}

	model = transaction.Model{
		ID:             "fe9f80f9-329c-44db-bd98-22c192bd93fc",
		Status:         "executed",
		Timestamp:      transaction.CSVDateTime{Time: timestamp},
		Type:           "",
		AssetType:      "",
		AssetName:      "MSCI EM USD (Dist)",
		ISIN:           "IE00B0M63177",
		Shares:         2.481328,
		Rate:           40.301,
		Yield:          0,
		Profit:         0,
		Commission:     0,
		Debit:          0,
		Credit:         100.00,
		TaxAmount:      0,
		InvestedAmount: 0,
		Documents:      []string{},
	}
}
