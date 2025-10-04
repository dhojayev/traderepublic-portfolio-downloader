package transaction_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeResolver_Resolve(t *testing.T) {
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

		resolver := transaction.NewTypeResolver()

		t.Run("it resolves type in  "+entry.Name(), func(t *testing.T) {
			t.Parallel()

			resolveedType, err := resolver.Resolve(details)
			require.NoError(t, err)

			assert.NotEqual(t, transaction.TypeUnknown, resolveedType)
		})
	}
}
