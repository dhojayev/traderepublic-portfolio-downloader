package details_test

import (
	"fmt"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestResponseContents(t *testing.T) {
	t.Parallel()

	testCases := []fakes.TransactionTestCase{
		fakes.PaymentInbound01,
		fakes.PaymentInboundSepaDirectDebit01,
		fakes.InterestPayoutCreated01,
		fakes.SavingsPlanExecuted01,
		fakes.OrderExecuted02,
		fakes.Credit01,
		fakes.BenefitsSpareChangeExecution01,
	}

	logger := log.New()
	controller := gomock.NewController(t)
	readerMock := reader.NewMockInterface(controller)
	client := details.NewClient(readerMock, logger)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (reader.JSONResponse, error) {
				return reader.NewJSONResponse(testCase.TimelineDetailsData.Raw), nil
			})

		var actual details.Response

		err := client.Details("1ae661c0-b3f1-4a81-a909-79567161b014", &actual)
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		headerSection, err := actual.SectionTypeHeader()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		assert.Equal(t, testCase.TimelineDetailsData.Unmarshalled.Header, headerSection)

		tableSections, err := actual.SectionsTypeTable()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		assert.Equal(t, testCase.TimelineDetailsData.Unmarshalled.Table, tableSections, fmt.Sprintf("case %d", i))

		// do not compare documents section if no expected value provided.
		if !reflect.DeepEqual(testCase.TimelineDetailsData.Unmarshalled.Documents, details.ResponseSectionTypeDocuments{}) {
			documentsSection, err := actual.SectionTypeDocuments()
			assert.NoError(t, err, fmt.Sprintf("case %d", i))

			assert.Equal(t, testCase.TimelineDetailsData.Unmarshalled.Documents, documentsSection, fmt.Sprintf("case %d", i))
		}
	}
}
