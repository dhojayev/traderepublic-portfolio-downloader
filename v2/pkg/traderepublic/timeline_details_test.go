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

	t.Run("it can find isin", func(t *testing.T) {
		t.Parallel()

		expected := "IE00B0M63177"
		details := getTestData(t)

		header, err := details.SectionHeader()
		require.NoError(t, err)

		actual := header.Action.Payload

		assert.Equal(t, expected, actual)
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

	t.Run("it can find shares", func(t *testing.T) {
		t.Parallel()

		details := getTestData(t)
		transaction, err := details.SectionTransaction()
		require.NoError(t, err)

		shares, err := transaction.DataShares()
		require.NoError(t, err)

		assert.Equal(t, "2,481328", shares.Detail.Text)
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
