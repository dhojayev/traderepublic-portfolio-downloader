package traderepublic_test

import (
	"os"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimelineDetailsJson_Section(t *testing.T) {
	t.Parallel()

	t.Run("isin can be fetched", func(t *testing.T) {
		t.Parallel()

		expected := "IE00B0M63177"
		response := getTestFileContents(t)

		var details traderepublic.TimelineDetailsJson

		err := details.UnmarshalJSON(response)
		require.NoError(t, err)

		var header traderepublic.HeaderSection

		err = details.Section(&header)
		require.NoError(t, err)

		actual := header.Action.Payload

		assert.Equal(t, expected, actual)
	})

	t.Run("SectionNotFound", func(t *testing.T) {
		t.Parallel()

		d := &traderepublic.TimelineDetailsJson{
			Sections: []any{
				map[string]any{"Type": "Footer"},
			},
		}

		var header traderepublic.HeaderSection

		err := d.Section(&header)
		assert.ErrorIs(t, err, traderepublic.ErrSectionNotFound)
	})
}

func getTestFileContents(t *testing.T) []byte {
	t.Helper()

	response, err := os.ReadFile("../../tests/fakes/fe9f80f9-329c-44db-bd98-22c192bd93fc.json")
	require.NoError(t, err)

	return response
}
