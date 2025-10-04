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
		expectedType := &transaction.SavingsPlanType{}
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
		b.SetID(expectedID)
		b.SetStatus(expectedStatus)
		b.SetTimestamp(expectedTimestamp)
		b.SetType(expectedType)
		b.SetAssetType(expectedAssetType)
		b.SetAssetName(expectedName)
		b.SetISIN(expectedInstrument)
		b.SetShares(expectedShares)
		b.SetRate(expectedRate)
		b.SetYield(expectedYield)
		b.SetProfit(expectedProfit)
		b.SetCommission(expectedCommission)
		b.SetDebit(expectedDebit)
		b.SetCredit(expectedCredit)
		b.SetTaxAmount(expectedTaxAmount)
		b.SetInvestedAmount(expectedInvestedAmount)
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
		assert.Equal(t, expectedRate, model.SharePrice)
		assert.Equal(t, expectedYield, model.Yield)
		assert.Equal(t, expectedProfit, model.Gain)
		assert.Equal(t, expectedCommission, model.Fee)
		assert.Equal(t, expectedDebit, model.Debit)
		assert.Equal(t, expectedCredit, model.Credit)
		assert.Equal(t, expectedTaxAmount, model.TaxAmount)
		assert.Equal(t, expectedInvestedAmount, model.InvestedAmount)
		assert.Equal(t, expectedDocuments, model.Documents)
	})
}
