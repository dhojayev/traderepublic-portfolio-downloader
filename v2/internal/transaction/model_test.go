package transaction_test

import (
	"testing"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/stretchr/testify/assert"
)

func TestModelBuilder(t *testing.T) {
	t.Parallel()

	t.Run("Builds Model with all fields", func(t *testing.T) {
		t.Parallel()

		expectedID := "12345"
		expectedStatus := "completed"
		expectedTimestamp := transaction.CSVDateTime{time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)}
		expectedType := transaction.TypeDeposit
		expectedAssetType := "stock"
		expectedName := "AAPL"
		expectedInstrument := "AAPL"
		expectedShares := 10.0
		expectedRate := 150.75
		expectedYield := 2.5
		expectedProfit := 30.0
		expectedCommission := 1.5
		expectedDebit := 1560.75
		expectedCredit := 0.0
		expectedTaxAmount := 5.0
		expectedInvestedAmount := 1495.25
		expectedDocuments := []string{"doc1.pdf", "doc2.pdf"}

		b := transaction.NewModelBuilder()
		b.WithID(expectedID)
		b.WithStatus(expectedStatus)
		b.WithTimestamp(expectedTimestamp)
		b.WithType(expectedType)
		b.WithAssetType(expectedAssetType)
		b.WithName(expectedName)
		b.WithInstrument(expectedInstrument)
		b.WithShares(expectedShares)
		b.WithRate(expectedRate)
		b.WithYield(expectedYield)
		b.WithProfit(expectedProfit)
		b.WithCommission(expectedCommission)
		b.WithDebit(expectedDebit)
		b.WithCredit(expectedCredit)
		b.WithTaxAmount(expectedTaxAmount)
		b.WithInvestedAmount(expectedInvestedAmount)
		b.AddDocument("doc1.pdf")
		b.AddDocument("doc2.pdf")

		model := b.Build()

		assert.Equal(t, expectedID, model.ID)
		assert.Equal(t, expectedStatus, model.Status)
		assert.Equal(t, expectedTimestamp, model.Timestamp)
		assert.Equal(t, expectedType, model.Type)
		assert.Equal(t, expectedAssetType, model.AssetType)
		assert.Equal(t, expectedName, model.AssetName)
		assert.Equal(t, expectedInstrument, model.ISIN)
		assert.Equal(t, expectedShares, model.Shares)
		assert.Equal(t, expectedRate, model.Rate)
		assert.Equal(t, expectedYield, model.Yield)
		assert.Equal(t, expectedProfit, model.Profit)
		assert.Equal(t, expectedCommission, model.Commission)
		assert.Equal(t, expectedDebit, model.Debit)
		assert.Equal(t, expectedCredit, model.Credit)
		assert.Equal(t, expectedTaxAmount, model.TaxAmount)
		assert.Equal(t, expectedInvestedAmount, model.InvestedAmount)
		assert.Equal(t, expectedDocuments, model.Documents)
	})
}
