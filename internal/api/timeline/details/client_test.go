package details_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
)

type content struct {
	data []byte
}

func (c content) Data() []byte {
	return c.data
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filepath                 string
		expectedHeaderSection    details.ResponseSectionTypeHeaderNew
		expectedTableSections    []details.ResponseSectionTypeTableNew
		expectedDocumentsSection details.ResponseSectionTypeDocumentsNew
	}{
		{
			filepath: "../../../../tests/data/transaction-details/payment-inbound-01.json",
			expectedHeaderSection: details.ResponseSectionTypeHeaderNew{
				Data: details.ResponseSectionTypeHeaderDataNew{
					Icon:      "logos/timeline_plus_circle/v2",
					Status:    "executed",
					Timestamp: "2023-05-21T08:25:53.360+0000",
				},
				Title: "Du hast 200,00 € erhalten",
				Type:  "header",
			},
			expectedTableSections: []details.ResponseSectionTypeTableNew{
				{
					Data: []details.ResponseSectionTypeTableDataNew{
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								FunctionalStyle: "EXECUTED",
								Text:            "Abgeschlossen",
								Type:            "status",
							},
							Style: "plain",
							Title: "Status",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "John Doe",
								Type: "text",
							},
							Style: "plain",
							Title: "Von",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "DE78 0000 0000 0000 0000 00",
								Type: "text",
							},
							Style: "plain",
							Title: "IBAN",
						},
					},
					Title: "Übersicht",
					Type:  "table",
				},
			},
		},
		{
			filepath: "../../../../tests/data/transaction-details/payment-inbound-sepa-direct-debit-01.json",
			expectedHeaderSection: details.ResponseSectionTypeHeaderNew{
				Data: details.ResponseSectionTypeHeaderDataNew{
					Icon:      "logos/timeline_plus_circle/v2",
					Status:    "executed",
					Timestamp: "2023-07-23T21:05:22.543+0000",
				},
				Title: "Du hast 500,00 € per Lastschrift hinzugefügt",
				Type:  "header",
			},
			expectedTableSections: []details.ResponseSectionTypeTableNew{
				{
					Data: []details.ResponseSectionTypeTableDataNew{
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								FunctionalStyle: "EXECUTED",
								Text:            "Ausgeführt",
								Type:            "status",
							},
							Style: "plain",
							Title: "Status",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "Lastschrift",
								Type: "text",
							},
							Style: "plain",
							Title: "Zahlung",
						},
					},
					Title: "Übersicht",
					Type:  "table",
				},
				{
					Data: []details.ResponseSectionTypeTableDataNew{
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "Gratis",
								Type: "text",
							},
							Style: "plain",
							Title: "Gebühr",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "500,00 €",
								Type: "text",
							},
							Style: "highlighted",
							Title: "Betrag",
						},
					},
					Title: "Transaktion",
					Type:  "table",
				},
			},
			expectedDocumentsSection: details.ResponseSectionTypeDocumentsNew{
				Data: []details.ResponseSectionTypeDocumentDataNew{
					{
						Action: details.ResponseActionNew{
							Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
							Type:    "browserModal",
						},
						Detail:      "23.07.2023",
						ID:          "cfc08704-eb56-44f1-83a0-c39aba9055ca",
						PostboxType: "PAYMENT_INBOUND_INVOICE",
						Title:       "Abrechnung Einzahlung",
					},
				},
				Title: "Dokumente",
				Type:  "documents",
			},
		},
		{
			filepath: "../../../../tests/data/transaction-details/interest-payout-created-01.json",
			expectedHeaderSection: details.ResponseSectionTypeHeaderNew{
				Data: details.ResponseSectionTypeHeaderDataNew{
					Icon:      "logos/timeline_interest_new/v2",
					Status:    "executed",
					Timestamp: "2023-11-06T11:22:52.544+0000",
				},
				Title: "Du hast 0,07 EUR erhalten",
				Type:  "header",
			},
			expectedTableSections: []details.ResponseSectionTypeTableNew{
				{
					Data: []details.ResponseSectionTypeTableDataNew{
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								FunctionalStyle: "EXECUTED",
								Text:            "Abgeschlossen",
								Type:            "status",
							},
							Style: "plain",
							Title: "Status",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "283,33 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Durchschnittssaldo",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "2 %",
								Type: "text",
							},
							Style: "plain",
							Title: "Jahressatz",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "Guthaben",
								Type: "text",
							},
							Style: "plain",
							Title: "Vermögenswert",
						},
					},
					Title: "Übersicht",
					Type:  "table",
				},
				{
					Data: []details.ResponseSectionTypeTableDataNew{
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "+ 0,09 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Angefallen",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "0,02 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Steuern",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetailNew{
								Text: "+ 0,07 €",
								Type: "text",
							},
							Style: "highlighted",
							Title: "Gesamt",
						},
					},
					Title: "Transaktion",
					Type:  "table",
				},
			},
			expectedDocumentsSection: details.ResponseSectionTypeDocumentsNew{
				Data: []details.ResponseSectionTypeDocumentDataNew{
					{
						Action: details.ResponseActionNew{
							Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
							Type:    "browserModal",
						},
						Detail:      "06.11.2023",
						ID:          "f1b33e1e-0e44-4508-b2cd-d508715d9740",
						PostboxType: "INTEREST_PAYOUT_INVOICE",
						Title:       "Abrechnung",
					},
				},
				Title: "Dokumente",
				Type:  "documents",
			},
		},
	}

	controller := gomock.NewController(t)
	readerMock := portfolio.NewMockReaderInterface(controller)
	client := details.NewClient(readerMock)

	for _, testCase := range testCases {
		readerMock.
			EXPECT().
			Read("timelineDetailV2", gomock.Any()).
			DoAndReturn(func(_ string, _ map[string]any) (portfolio.OutputDataInterface, error) {
				fileContents, err := os.ReadFile(testCase.filepath)

				return content{data: fileContents}, err
			})

		actual, err := client.Get("1ae661c0-b3f1-4a81-a909-79567161b014")
		assert.Nil(t, err)

		headerSection, err := actual.SectionTypeHeader()
		assert.Nil(t, err)

		assert.Equal(t, testCase.expectedHeaderSection, headerSection)

		tableSections, err := actual.SectionsTypeTable()
		assert.Nil(t, err)

		assert.Equal(t, testCase.expectedTableSections, tableSections)

		// do not compare documents section if no expected value provided.
		if !reflect.DeepEqual(testCase.expectedDocumentsSection, details.ResponseSectionTypeDocumentsNew{}) {
			documentsSection, err := actual.SectionTypeDocuments()
			assert.Nil(t, err)

			assert.Equal(t, testCase.expectedDocumentsSection, documentsSection)
		}
	}
}
