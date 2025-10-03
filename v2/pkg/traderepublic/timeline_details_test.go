package traderepublic_test

import (
	"os"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimelineDetailsJson_SectionHeader(t *testing.T) {
	t.Parallel()

	details := getTestData(t)

	header, err := details.SectionHeader()
	require.NoError(t, err)

	t.Run("it can find status", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, traderepublic.HeaderSectionDataStatusExecuted, header.Data.Status)
	})

	t.Run("it can find timestamp", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "2025-01-02T14:52:18.686+0000", header.Data.Timestamp)
	})

	t.Run("it can find isin", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "IE00B0M63177", header.Action.Payload)
	})

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

	details := getTestData(t)
	transaction, err := details.SectionTransaction()
	require.NoError(t, err)

	t.Run("it can find shares", func(t *testing.T) {
		t.Parallel()

		shares, err := transaction.DataPayment(traderepublic.PaymentShares)
		require.NoError(t, err)

		assert.Equal(t, "2,481328", shares.Detail.Text)
	})

	t.Run("it can find share price", func(t *testing.T) {
		t.Parallel()

		sharePrice, err := transaction.DataPayment(traderepublic.PaymentSharePrice)
		require.NoError(t, err)

		assert.Equal(t, "40,301 €", sharePrice.Detail.Text)
	})

	t.Run("it can find commission", func(t *testing.T) {
		t.Parallel()

		commission, err := transaction.DataPayment(traderepublic.PaymentCommission)
		require.NoError(t, err)

		assert.Equal(t, "Gratis", commission.Detail.Text)
	})

	t.Run("it can find total", func(t *testing.T) {
		t.Parallel()

		total, err := transaction.DataPayment(traderepublic.PaymentTotal)
		require.NoError(t, err)

		assert.Equal(t, "100,00 €", total.Detail.Text)
	})
}

func getTestData(t *testing.T) traderepublic.TimelineDetailsJson {
	t.Helper()

	response, err := os.ReadFile("../../tests/fakes/fe9f80f9-329c-44db-bd98-22c192bd93fc.json")
	require.NoError(t, err)

	var details traderepublic.TimelineDetailsJson

	err = details.UnmarshalJSON(response)
	require.NoError(t, err)

	return details
}
