package traderepublic_test

import (
	"os"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimelineDetails_HeaderSection(t *testing.T) {
	t.Parallel()

	t.Run("isin can be fetched", func(t *testing.T) {
		t.Parallel()

		var header traderepublic.HeaderSection

		expected := "IE00B0M63177"
		response, err := os.ReadFile("../../tests/fakes/fe9f80f9-329c-44db-bd98-22c192bd93fc.json")
		require.NoError(t, err)

		var details traderepublic.TimelineDetailsJson

		err = details.UnmarshalJSON(response)
		require.NoError(t, err)

		err = details.Section(&header)
		require.NoError(t, err)

		actual := header.Action.Payload

		assert.Equal(t, expected, actual)
	})
}
