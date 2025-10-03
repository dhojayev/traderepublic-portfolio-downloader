package traderepublic_test

import (
	"os"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	filepath         string
	status           traderepublic.HeaderSectionDataStatus
	timestamp        string
	isin             string
	shares           string
	sharePrice       string
	dividendPerShare string
	fee              string
	tax              string
	total            string
}

func TestTimelineDetailsJson_SectionHeader(t *testing.T) {
	t.Parallel()

	testCases := getTestData(t)

	for _, testCase := range testCases {
		contents, err := os.ReadFile(testCase.filepath)
		require.NoError(t, err)

		var details traderepublic.TimelineDetailsJson

		err = details.UnmarshalJSON(contents)
		require.NoError(t, err)

		header, err := details.SectionHeader()
		require.NoError(t, err)

		t.Run("it can find status", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.status, header.Data.Status)
		})

		t.Run("it can find timestamp", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.timestamp, header.Data.Timestamp)
		})

		t.Run("it can find isin", func(t *testing.T) {
			t.Parallel()

			if testCase.isin == "" {
				t.Skip()
			}

			require.NotNil(t, header.Action)
			assert.Equal(t, testCase.isin, header.Action.Payload)
		})
	}

	t.Run("error returned on header section not found", func(t *testing.T) {
		t.Parallel()

		d := &traderepublic.TimelineDetailsJson{
			Sections: []any{
				map[string]any{"Type": "Footer"},
			},
		}

		_, err := d.SectionHeader()
		assert.ErrorIs(t, err, traderepublic.ErrSectionNotFound)
	})

}

func TestTimelineDetailsJson_SectionTrnsaction(t *testing.T) {
	t.Parallel()

	testCases := getTestData(t)

	for _, testCase := range testCases {
		contents, err := os.ReadFile(testCase.filepath)
		require.NoError(t, err)

		var details traderepublic.TimelineDetailsJson

		err = details.UnmarshalJSON(contents)
		require.NoError(t, err)

		transaction, err := details.SectionTable(traderepublic.SectionTableTransaction)
		require.NoError(t, err)

		t.Run("it can find shares", func(t *testing.T) {
			t.Parallel()

			shares, err := transaction.DataPayment(traderepublic.DataShares)
			require.NoError(t, err)

			assert.Equal(t, testCase.shares, shares.Detail.Text)
		})

		t.Run("it can find share price", func(t *testing.T) {
			t.Parallel()

			if testCase.sharePrice == "" {
				t.Skip()
			}

			sharePrice, err := transaction.DataPayment(traderepublic.DataSharePrice)
			require.NoError(t, err)

			assert.Equal(t, testCase.sharePrice, sharePrice.Detail.Text)
		})

		t.Run("it can find dividend per share", func(t *testing.T) {
			t.Parallel()

			if testCase.dividendPerShare == "" {
				t.Skip()
			}

			dividendPerShare, err := transaction.DataPayment(traderepublic.DataDividendPerShare)
			require.NoError(t, err)

			assert.Equal(t, testCase.dividendPerShare, dividendPerShare.Detail.Text)
		})

		t.Run("it can find fee", func(t *testing.T) {
			t.Parallel()

			if testCase.fee == "" {
				t.Skip()
			}

			fee, err := transaction.DataPayment(traderepublic.DataFee)
			require.NoError(t, err)

			assert.Equal(t, testCase.fee, fee.Detail.Text)
		})

		t.Run("it can find tax", func(t *testing.T) {
			t.Parallel()

			if testCase.tax == "" {
				t.Skip()
			}

			tax, err := transaction.DataPayment(traderepublic.DataTax)
			require.NoError(t, err)

			assert.Equal(t, testCase.tax, tax.Detail.Text)
		})

		t.Run("it can find total", func(t *testing.T) {
			t.Parallel()

			total, err := transaction.DataPayment(traderepublic.DataTotal)
			require.NoError(t, err)

			assert.Equal(t, testCase.total, total.Detail.Text)
		})
	}
}

func getTestData(t *testing.T) []testCase {
	t.Helper()

	return []testCase{
		{
			filepath:   "../../tests/fakes/fe9f80f9-329c-44db-bd98-22c192bd93fc.json",
			status:     traderepublic.HeaderSectionDataStatusExecuted,
			timestamp:  "2025-01-02T14:52:18.686+0000",
			isin:       "IE00B0M63177",
			shares:     "2.481328",
			sharePrice: "€40.301",
			fee:        "Free",
			total:      "€100.00",
		},
		{
			filepath:         "../../tests/fakes/a0e4c36a-e0ee-4183-a725-09fb1c6b3c33.json",
			status:           traderepublic.HeaderSectionDataStatusExecuted,
			timestamp:        "2024-06-26T15:22:31.478Z",
			shares:           "30.447001",
			dividendPerShare: "0,15 $",
			tax:              "0,00 €",
			total:            "4,13 €",
		},
	}
}
