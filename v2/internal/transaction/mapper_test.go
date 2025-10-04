package transaction_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	gocache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataMapper_Map(t *testing.T) {
	t.Parallel()

	detailsPath := "../../debug/responses/timeline_detail_v2_received"
	_, err := os.Stat(detailsPath)
	if err != nil && os.IsNotExist(err) {
		t.Skip()
	}

	entries, err := os.ReadDir(detailsPath)
	require.NoError(t, err)

	cache := gocache.New(gocache.NoExpiration, gocache.NoExpiration)

	instrPath := "../../debug/responses/instrument_received"
	_, err = os.Stat(instrPath)
	if err == nil {
		files, err := os.ReadDir(instrPath)
		require.NoError(t, err)

		for _, file := range files {
			contents, err := os.ReadFile(filepath.Join(instrPath, file.Name()))
			require.NoError(t, err)

			var instr traderepublic.InstrumentJson

			err = instr.UnmarshalJSON(contents)
			require.NoError(t, err)

			cache.Set(instr.Isin, instr, gocache.NoExpiration)
		}
	}

	resolver := transaction.NewTypeResolver()
	mapper := transaction.NewDataMapper(cache)

	for _, entry := range entries {
		contents, err := os.ReadFile(filepath.Join(detailsPath, entry.Name()))
		require.NoError(t, err)

		var details traderepublic.TimelineDetailsJson

		err = details.UnmarshalJSON(contents)
		require.NoError(t, err)

		t.Run("it can map all fields "+entry.Name(), func(t *testing.T) {
			t.Parallel()

			model := transaction.Model{}
			err := resolver.SetType(details, &model)
			if err != nil {
				switch {
				case errors.Is(err, transaction.ErrCancelledTransactionReceived), errors.Is(err, transaction.ErrIgnoredTransactionReceived), errors.Is(err, transaction.ErrUnknownTransactionReceived):
					t.Skip()
				}

				require.NoError(t, err)
			}

			err = mapper.Map(details, &model)
			if errors.Is(err, transaction.ErrTransactionWithoutTypeReceived) {
				t.Skip()
			}

			require.NoError(t, err)

			assert.NotEmpty(t, model.ID)
			assert.NotEmpty(t, model.Status)
			assert.NotEmpty(t, model.Timestamp)
			assert.NotEmpty(t, model.ISIN)
			assert.NotEmpty(t, model.AssetName)
			assert.NotEmpty(t, model.AssetType)
			assert.NotEmpty(t, model.Shares)
			assert.NotEmpty(t, model.SharePrice)
			assert.NotNil(t, model.Fee)
			assert.NotEmpty(t, model.Debit)
		})
	}
}
