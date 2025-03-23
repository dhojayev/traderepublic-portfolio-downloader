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

var OrderExecied06 = TransactionTestCase{
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
  "id": "e8aecbac-cc49-4a08-ac25-588afe45df70",
  "sections": [
    {
      "action": {
        "payload": "DE000A2TEDB8",
        "type": "instrumentDetail"
      },
      "data": {
        "icon": "logos/DE000A2TEDB8/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2025-03-12T15:21:56.707+0000"
      },
      "title": "Du hast 1.001,77 € erhalten",
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
            "displayValue": null,
            "text": "Limit Verkauf",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Orderart"
        },
        {
          "detail": {
            "action": null,
            "displayValue": null,
            "text": "Anleihe Feb. 2024",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Asset"
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
            "action": null,
            "displayValue": null,
            "text": "9 %",
            "trend": "positive",
            "type": "text"
          },
          "style": "plain",
          "title": "Rendite"
        },
        {
          "detail": {
            "action": null,
            "displayValue": null,
            "text": "1.001,77 €",
            "trend": "positive",
            "type": "text"
          },
          "style": "plain",
          "title": "Gewinn"
        }
      ],
      "title": "Performance",
      "type": "horizontalTable"
    },
    {
      "action": {
        "payload": {
          "sections": [
            {
              "action": null,
              "title": "Transaktion",
              "type": "title"
            },
            {
              "action": null,
              "data": [
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "50",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Anteile"
                },
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "123,22 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Aktienkurs"
                },
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "+ 5.008,33 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "highlighted",
                  "title": "Order"
                }
              ],
              "title": null,
              "type": "table"
            },
            {
              "action": null,
              "data": [
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "1,02 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Kapitalertragsteuer"
                },
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "0,00 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Solidaritätszuschlag"
                },
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "1,02 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "highlighted",
                  "title": "Steuern"
                }
              ],
              "title": null,
              "type": "table"
            },
            {
              "action": null,
              "data": [
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "+ 5.008,33 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Order"
                },
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "1,02 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Steuern"
                },
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "1,00 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Gebühr"
                },
                {
                  "detail": {
                    "action": null,
                    "displayValue": null,
                    "text": "+ 5.006,31 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "highlighted",
                  "title": "Gesamt"
                }
              ],
              "title": null,
              "type": "table"
            }
          ]
        },
        "type": "infoPage"
      },
      "data": [
        {
          "detail": {
            "action": null,
            "displayValue": null,
            "text": "50",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Anteile"
        },
        {
          "detail": {
            "action": null,
            "displayValue": null,
            "text": "123,22 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Aktienkurs"
        },
        {
          "detail": {
            "action": null,
            "displayValue": null,
            "text": "1,02 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Steuern"
        },
        {
          "detail": {
            "action": null,
            "displayValue": null,
            "text": "1,00 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Gebühr"
        },
        {
          "detail": {
            "action": null,
            "displayValue": null,
            "text": "+ 5.006,31 €",
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
            "payload": "/3d219fea-b09a-49e5-a31d-d33c13c38d5e/",
            "type": "browserModal"
          },
          "detail": "11.03.2025",
          "id": "3d219fea-b09a-49e5-a31d-d33c13c38d5e",
          "postboxType": "SECURITIES_SETTLEMENT",
          "title": "Abrechnung"
        },
        {
          "action": {
            "payload": "/ba06c258-cf61-4dcc-9201-efb4669575e9/",
            "type": "browserModal"
          },
          "detail": "12.03.2025",
          "id": "ba06c258-cf61-4dcc-9201-efb4669575e9",
          "postboxType": "CONFIRM_ORDER_CREATE_V2",
          "title": "Auftragsbestätigung"
        },
        {
          "action": {
            "payload": "/66ae51de-05e5-4268-972e-6079f77c39b7/",
            "type": "browserModal"
          },
          "detail": "13.03.2025",
          "id": "66ae51de-05e5-4268-972e-6079f77c39b7",
          "postboxType": "COSTS_INFO_SELL_V2",
          "title": "Kosteninformation"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    }
  ]
}`),
		Normalized: details.NormalizedResponse{},
	},
	EventType: transactions.EventTypeOrderExecuted,
	Transaction: transaction.Model{
		UUID: "e8aecbac-cc49-4a08-ac25-588afe45df70",
		Instrument: instrument.Model{
			ISIN: "DE000A2TEDB8",
			Name: "Anleihe Feb. 2024",
			Icon: "logos/DE000A2TEDB8/v2",
			Type: instrument.TypeOther,
		},
		Documents: []document.Model{
			{
				TransactionUUID: "e8aecbac-cc49-4a08-ac25-588afe45df70",
				ID:              "3d219fea-b09a-49e5-a31d-d33c13c38d5e",
				URL:             "/3d219fea-b09a-49e5-a31d-d33c13c38d5e/",
				Detail:          "11.03.2025",
				Title:           "Abrechnung",
				Filepath:        "2025-03/e8aecbac-cc49-4a08-ac25-588afe45df70/Abrechnung.pdf",
			},
			{
				TransactionUUID: "e8aecbac-cc49-4a08-ac25-588afe45df70",
				ID:              "ba06c258-cf61-4dcc-9201-efb4669575e9",
				URL:             "/ba06c258-cf61-4dcc-9201-efb4669575e9/",
				Detail:          "12.03.2025",
				Title:           "Auftragsbestätigung",
				Filepath:        "2025-03/e8aecbac-cc49-4a08-ac25-588afe45df70/Auftragsbestätigung.pdf",
			},
			{
				TransactionUUID: "e8aecbac-cc49-4a08-ac25-588afe45df70",
				ID:              "66ae51de-05e5-4268-972e-6079f77c39b7",
				URL:             "/66ae51de-05e5-4268-972e-6079f77c39b7/",
				Detail:          "13.03.2025",
				Title:           "Kosteninformation",
				Filepath:        "2025-03/e8aecbac-cc49-4a08-ac25-588afe45df70/Kosteninformation.pdf",
			},
		},
		Type:       transaction.TypeSale,
		Status:     "executed",
		Yield:      9,
		Profit:     1001.77,
		Shares:     50,
		Rate:       123.22,
		Commission: 1,
		Total:      5006.31,
		TaxAmount:  1.02,
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "e8aecbac-cc49-4a08-ac25-588afe45df70",
		Status:     "executed",
		Type:       "Sale",
		AssetType:  "Other",
		Name:       "Anleihe Feb. 2024",
		Instrument: "DE000A2TEDB8",
		Yield:      9,
		Profit:     1001.77,
		Shares:     -50,
		Rate:       123.22,
		Commission: 1,
		Debit:      0,
		Credit:     5006.31,
		TaxAmount:  1.02,
		Documents: []string{
			"2025-03/e8aecbac-cc49-4a08-ac25-588afe45df70/Abrechnung.pdf",
			"2025-03/e8aecbac-cc49-4a08-ac25-588afe45df70/Auftragsbestätigung.pdf",
			"2025-03/e8aecbac-cc49-4a08-ac25-588afe45df70/Kosteninformation.pdf",
		},
	},
}

func init() {
	OrderExecied06.Transaction.Timestamp, _ = internal.ParseTimestamp("2025-03-12T15:21:56.707+0000")
	OrderExecied06.CSVEntry.Timestamp = internal.DateTime{Time: OrderExecied06.Transaction.Timestamp}

	RegisterSupported("OrderExecied06", OrderExecied06)
}
