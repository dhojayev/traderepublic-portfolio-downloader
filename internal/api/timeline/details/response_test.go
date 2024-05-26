package details_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

func TestResponseContents(t *testing.T) {
	t.Parallel()

	testCases := []tests.TestCase{
		fakes.PaymentInbound01,
		fakes.PaymentInboundSepaDirectDebit01,
		fakes.InterestPayoutCreated01,
		fakes.SavingsPlanExecuted01,
		fakes.OrderExecuted02,
		fakes.Credit01,
		fakes.BenefitsSpareChangeExecution01,
	}

	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	client := details.NewClient(readerMock)

	for i, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				return filesystem.NewOutputData([]byte(testCase.ResponseJSON)), nil
			})

		actual, err := client.Get("1ae661c0-b3f1-4a81-a909-79567161b014")
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		headerSection, err := actual.SectionTypeHeader()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		assert.Equal(t, testCase.Response.HeaderSection, headerSection)

		tableSections, err := actual.SectionsTypeTable()
		assert.NoError(t, err, fmt.Sprintf("case %d", i))

		assert.Equal(t, testCase.Response.TableSections, tableSections, fmt.Sprintf("case %d", i))

		// do not compare documents section if no expected value provided.
		if !reflect.DeepEqual(testCase.Response.DocumentsSection, details.ResponseSectionTypeDocuments{}) {
			documentsSection, err := actual.SectionTypeDocuments()
			assert.NoError(t, err, fmt.Sprintf("case %d", i))

			assert.Equal(t, testCase.Response.DocumentsSection, documentsSection, fmt.Sprintf("case %d", i))
		}
	}
}
