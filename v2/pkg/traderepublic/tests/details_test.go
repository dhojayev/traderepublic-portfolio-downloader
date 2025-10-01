package tests_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimelineDetails(t *testing.T) {
	t.Parallel()

	t.Run("it unmarshals response", func(t *testing.T) {
		t.Parallel()

		expected := "IE00B0M63177"
		response, err := os.ReadFile("../../../tests/fakes/fe9f80f9-329c-44db-bd98-22c192bd93fc.json")
		require.NoError(t, err)

		var details traderepublic.TimelineDetailsJson

		err = details.UnmarshalJSON(response)
		require.NoError(t, err)

		var actual string

		for _, section := range details.Sections {
			var header traderepublic.HeaderSection

			data, err := json.Marshal(section)
			if err != nil {
				continue
			}

			err = header.UnmarshalJSON(data)
			if err != nil {
				continue
			}

			actual = header.Action.Payload

			break
		}

		assert.Equal(t, expected, actual)
	})
}
