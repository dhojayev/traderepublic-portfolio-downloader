package timelinedetails_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/timelinedetails"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizer_Normalize(t *testing.T) {
	t.Parallel()

	path := "../../debug/responses/timeline_detail_v2_received"
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		t.Skip()
	}

	entries, err := os.ReadDir(path)
	require.NoError(t, err)

	for _, entry := range entries {
		contents, err := os.ReadFile(filepath.Join(path, entry.Name()))
		require.NoError(t, err)

		var details traderepublic.TimelineDetailsJson

		err = details.UnmarshalJSON(contents)
		require.NoError(t, err)

		builder := transaction.NewModelBuilder()
		normalizer := timelinedetails.NewNormalizer(builder)

		t.Run("it normalizes "+entry.Name(), func(t *testing.T) {
			t.Parallel()

			model, err := normalizer.Normalize(details)
			require.NoError(t, err)

			assert.NotEmpty(t, model.ID)
			assert.NotEmpty(t, model.Status)
			assert.NotEmpty(t, model.Timestamp)
		})

		t.Run("it finds order type in "+entry.Name(), func(t *testing.T) {
			t.Parallel()

			header, err := details.SectionHeader()
			require.NoError(t, err)

			overview, err := details.FindSection(traderepublic.SectionOverview)
			require.NoError(t, err)

			_, err = overview.FindData(traderepublic.DataCardPayment)
			if err != nil && header.Action == nil {
				t.Skip()
			}

			orderType, err := overview.FindData(traderepublic.DataOrderType)
			require.NoError(t, err, "Unexpected contents: "+string(contents))
			require.NotEmpty(t, orderType)
		})
	}
}
