package activitylog_test

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/activitylog"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestClient_Get(t *testing.T) {
	t.Parallel()

	testCases := fakes.ActivityLogTestCasesSupported

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := reader.NewMockInterface(controller)
	client := activitylog.NewClient(readerMock, logger)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read(activitylog.RequestDataType, gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (reader.JSONResponse, error) {
				return reader.NewJSONResponse(testCase.ActivityLogData.Raw), nil
			})

		var actual []activitylog.ResponseItem
		err := client.List(&actual)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		if err != nil {
			continue
		}

		assert.Len(t, actual, 1, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.ActivityLogData.Unmarshalled, actual[0], fmt.Sprintf("case %d", i))
	}
}
