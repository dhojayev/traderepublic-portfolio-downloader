package transaction_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataMapper_Map(t *testing.T) {
	t.Parallel()

	path := "../../debug/responses/timeline_detail_v2_received"
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		t.Skip()
	}

	entries, err := os.ReadDir(path)
	require.NoError(t, err)

	resolver := transaction.NewTypeResolver()
	mapper := transaction.NewDataMapper()

	ignoredTypes := []transaction.TransactionType{
		transaction.TypeCardPayment,
		transaction.TypeCardRefund,
		transaction.TypeIgnored,
	}

	for _, entry := range entries {
		contents, err := os.ReadFile(filepath.Join(path, entry.Name()))
		require.NoError(t, err)

		var details traderepublic.TimelineDetailsJson

		err = details.UnmarshalJSON(contents)
		require.NoError(t, err)

		t.Run("it can map all fields "+entry.Name(), func(t *testing.T) {
			t.Parallel()

			trnType, err := resolver.Resolve(details)
			require.NoError(t, err)

			if slices.Contains(ignoredTypes, trnType) {
				t.Skip()
			}

			if details.Id == "a8a104a8-2d0c-43d8-877f-3f01b884ed0e" {
				t.Log("test")
			}

			model, err := mapper.Map(trnType, details)
			if errors.Is(err, transaction.ErrIgnoredTransactionReceived) {
				t.Skip()
			}

			require.NoError(t, err)

			assert.NotEmpty(t, model.ID)
			assert.NotEmpty(t, model.Status)
			assert.NotEmpty(t, model.Timestamp)

			switch trnType {
			case transaction.TypeSavingsplan:
				assert.NotEmpty(t, model.ISIN)
				assert.NotEmpty(t, model.Shares)
			case transaction.TypeBuyOrder:
				assert.NotEmpty(t, model.ISIN)
				assert.NotEmpty(t, model.Shares)
			case transaction.TypeSellOrder:
				assert.NotEmpty(t, model.ISIN)
				assert.NotEmpty(t, model.Shares)
			case transaction.TypeDividendsIncome:
				assert.NotEmpty(t, model.ISIN)
				assert.NotEmpty(t, model.Shares)
			case transaction.TypeRoundUp:
				assert.NotEmpty(t, model.ISIN)
				assert.NotEmpty(t, model.Shares)
			case transaction.TypeSaveback:
				assert.NotEmpty(t, model.ISIN)
				assert.NotEmpty(t, model.Shares)
			case transaction.TypeDeposit:
				assert.Empty(t, model.ISIN)
				assert.Empty(t, model.Shares)
			case transaction.TypeWithdrawal:
				assert.Empty(t, model.ISIN)
				assert.Empty(t, model.Shares)
			case transaction.TypeInterestPayment:
				assert.Empty(t, model.ISIN)
				assert.Empty(t, model.Shares)
			}
		})
	}
}
