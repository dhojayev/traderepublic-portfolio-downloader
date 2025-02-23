package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/instrument"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
)

var SavingsPlanInvoiceCreated01 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`[
{
      "action": {
        "payload": "398687aa-748b-47aa-951a-8583c5d141bf",
        "type": "timelineDetail"
      },
      "amount": {
        "currency": "EUR",
        "fractionDigits": 2,
        "value": -0.999886
      },
      "badge": null,
      "eventType": "SAVINGS_PLAN_INVOICE_CREATED",
      "icon": "logos/US0231351067/v2",
      "id": "398687aa-748b-47aa-951a-8583c5d141bf",
      "status": "EXECUTED",
      "subAmount": null,
      "subtitle": "Sparplan ausgeführt",
      "timestamp": "2024-10-16T14:55:20.954+0000",
      "title": "Amazon.com"
}
]`),
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
"id": "398687aa-748b-47aa-951a-8583c5d141bf",
  "sections": [
    {
      "action": {
        "payload": "US0231351067",
        "type": "instrumentDetail"
      },
      "data": {
        "icon": "logos/US0231351067/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2024-10-09T15:20:39.847+0000"
      },
      "title": "Du hast 1,00 € investiert",
      "type": "header"
    },
    {
      "action": null,
      "data": [
        {
          "detail": {
            "functionalStyle": "EXECUTED",
            "text": "Ausgeführt",
            "type": "status"
          },
          "style": "plain",
          "title": "Status"
        },
        {
          "detail": {
            "action": null,
            "text": "Sparplan",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Orderart"
        },
        {
          "detail": {
            "action": null,
            "text": "Amazon.com",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Asset"
        },
        {
          "detail": {
            "icon": "logos/bank_traderepublic/v2",
            "text": "Guthaben",
            "type": "iconWithText"
          },
          "style": "plain",
          "title": "Zahlung"
        }
      ],
      "title": "Übersicht",
      "type": "table"
    },
    {
      "action": null,
      "data": [
        {
          "detail": {
            "action": {
              "payload": {
                "savingsPlanId": "7e8259b7-aef3-46e4-bcd9-991841325f7c"
              },
              "type": "openSavingsPlanOverview"
            },
            "amount": "1,00 €",
            "icon": "logos/US023135****/v2",
            "status": "executed",
            "subtitle": "Wöchentlich",
            "timestamp": "2024-10-09T15:20:39.847Z",
            "title": "Amazon.com",
            "type": "embeddedTimelineItem"
          },
          "style": "plain",
          "title": ""
        }
      ],
      "title": "Sparplan",
      "type": "table"
    },
    {
      "action": null,
      "data": [
        {
          "detail": {
            "action": null,
            "text": "0,005941",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Anteile"
        },
        {
          "detail": {
            "action": null,
            "text": "168,30 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Anteilspreis"
        },
        {
          "detail": {
            "action": null,
            "text": "Gratis",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Gebühr"
        },
        {
          "detail": {
            "action": null,
            "text": "1,00 €",
            "trend": null,
            "type": "text"
          },
          "style": "highlighted",
          "title": "Gesamt"
        }
      ],
      "title": "Transaktion",
      "type": "table"
    },
    {
      "action": null,
      "data": [
        {
          "action": {
            "payload": "****",
            "type": "browserModal"
          },
          "detail": "09.10.2024",
          "id": "d3bacd7f-c63d-4d68-aef8-49130951ce9c",
          "postboxType": "SAVINGS_PLAN_EXECUTED_V2",
          "title": "Abrechnung Ausführung"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    }
  ]
}`),
		Normalized: details.NormalizedResponse{
			ID: "398687aa-748b-47aa-951a-8583c5d141bf",
			Header: details.NormalizedResponseHeaderSection{
				Action: details.NormalizedResponseSectionAction{
					Payload: "US0231351067",
					Type:    "instrumentDetail",
				},
				Data: details.NormalizedResponseHeaderSectionData{
					Icon:      "logos/US0231351067/v2",
					Status:    "executed",
					Timestamp: "2024-10-09T15:20:39.847+0000",
				},
				Title: "Du hast 1,00 € investiert",
				Type:  "header",
			},
			Overview: details.NormalizedResponseOverviewSection{
				NormalizedResponseTableSection: details.NormalizedResponseTableSection{
					Data: []details.NormalizedResponseTableSectionData{
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								FunctionalStyle: "EXECUTED",
								Text:            "Ausgeführt",
								Type:            "status",
							},
							Style: "plain",
							Title: "Status",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "Sparplan",
								Type: "text",
							},
							Style: "plain",
							Title: "Orderart",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "Amazon.com",
								Type: "text",
							},
							Style: "plain",
							Title: "Asset",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Icon: "logos/bank_traderepublic/v2",
								Text: "Guthaben",
								Type: "iconWithText",
							},
							Style: "plain",
							Title: "Zahlung",
						},
					},
					Title: "Übersicht",
					Type:  "table",
				},
			},
			Transaction: details.NormalizedResponseTransactionSection{
				NormalizedResponseTableSection: details.NormalizedResponseTableSection{
					Data: []details.NormalizedResponseTableSectionData{
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "0,005941",
								Type: "text",
							},
							Style: "plain",
							Title: "Anteile",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "168,30 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Anteilspreis",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "Gratis",
								Type: "text",
							},
							Style: "plain",
							Title: "Gebühr",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "1,00 €",
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
			Documents: details.NormalizedResponseDocumentsSection{
				Data: []details.NormalizedResponseDocumentsSectionData{
					{
						Action: details.NormalizedResponseSectionAction{
							Payload: "****",
							Type:    "browserModal",
						},
						Detail:      "09.10.2024",
						ID:          "d3bacd7f-c63d-4d68-aef8-49130951ce9c",
						PostboxType: "SAVINGS_PLAN_EXECUTED_V2",
						Title:       "Abrechnung Ausführung",
					},
				},
				Title: "Dokumente",
				Type:  "documents",
			},
		},
	},
	EventType: transactions.EventTypeSavingsPlanInvoiceCreated,
	Transaction: transaction.Model{
		UUID: "398687aa-748b-47aa-951a-8583c5d141bf",
		Instrument: instrument.Model{
			ISIN: "US0231351067",
			Name: "Amazon.com",
			Icon: "logos/US0231351067/v2",
			Type: instrument.TypeOther,
		},
		Type:   transaction.TypePurchase,
		Status: "executed",
		Shares: 0.005941,
		Rate:   168.30,
		Total:  1,
		Documents: []document.Model{
			{
				TransactionUUID: "398687aa-748b-47aa-951a-8583c5d141bf",
				ID:              "d3bacd7f-c63d-4d68-aef8-49130951ce9c",
				URL:             "****",
				Detail:          "09.10.2024",
				Title:           "Abrechnung Ausführung",
				Filepath:        "2024-10/398687aa-748b-47aa-951a-8583c5d141bf/Abrechnung Ausführung.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "398687aa-748b-47aa-951a-8583c5d141bf",
		Status:     "executed",
		Type:       transaction.TypePurchase,
		AssetType:  string(instrument.TypeOther),
		Name:       "Amazon.com",
		Instrument: "US0231351067",
		Shares:     0.005941,
		Rate:       168.30,
		Debit:      1,
	},
}

func init() {
	SavingsPlanInvoiceCreated01.Transaction.Timestamp, _ = internal.ParseTimestamp("2024-10-09T15:20:39.847+0000")
	SavingsPlanInvoiceCreated01.CSVEntry.Timestamp = internal.DateTime{Time: SavingsPlanInvoiceCreated01.Transaction.Timestamp}

	RegisterSupported("SavingsPlanInvoiceCreated01", SavingsPlanInvoiceCreated01)
}
