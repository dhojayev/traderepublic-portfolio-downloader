package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/instrument"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
)

var InterestPayoutCreated02 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`{
		"items": [
			{
				"action": {
					"payload": "79bea4ff-5552-45a1-85b5-8977e29f3d04",
					"type": "timelineDetail"
				},
				"amount": {
					"currency": "EUR",
					"fractionDigits": 2,
					"value": 20.28
				},
				"badge": null,
				"eventType": "INTEREST_PAYOUT_CREATED",
				"icon": "logos/timeline_interest_new/v2",
				"id": "79bea4ff-5552-45a1-85b5-8977e29f3d04",
				"status": "EXECUTED",
				"subAmount": null,
				"subtitle": "4,00% p.a.",
				"timestamp": "2024-03-03T20:45:47.367+0000",
				"title": "Zinsen"
				}
			]
		}`),
		Unmarshalled: transactions.ResponseItem{},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
  "id": "79bea4ff-5552-45a1-85b5-8977e29f3d04",
  "sections": [
    {
      "action": null,
      "data": {
        "icon": "logos/timeline_interest_new/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2024-03-03T20:45:47.367+0000"
      },
      "title": "Du hast 20,28 EUR erhalten",
      "type": "header"
    },
    {
      "action": null,
      "data": [
        {
          "detail": {
            "functionalStyle": "EXECUTED",
            "text": "Abgeschlossen",
            "type": "status"
          },
          "style": "plain",
          "title": "Status"
        },
        {
          "detail": {
            "action": null,
            "text": "6.294,86 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Durchschnittssaldo"
        },
        {
          "detail": {
            "action": null,
            "text": "4 %",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Jahressatz"
        },
        {
          "detail": {
            "action": null,
            "text": "Guthaben",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Vermögenswert"
        }
      ],
      "title": "Übersicht",
      "type": "table"
    },
    {
      "action": null,
      "data": [
        {
          "action": {
            "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
            "type": "browserModal"
          },
          "detail": "03.03.2024",
          "id": "2acd16e5-efc9-4bff-b171-3a31c300a628",
          "postboxType": "INTEREST_PAYOUT_INVOICE",
          "title": "Abrechnung"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    }
  ]
}`),
	},
	EventType: transactions.EventTypeInterestPayoutCreated,
	Transaction: transaction.Model{
		UUID:   "79bea4ff-5552-45a1-85b5-8977e29f3d04",
		Type:   transaction.TypeInterestPayout,
		Status: "executed",
		Total:  20.28,
		Instrument: instrument.Model{
			Icon: "logos/timeline_interest_new/v2",
			Type: instrument.TypeCash,
		},
		Documents: []document.Model{
			{
				TransactionUUID: "79bea4ff-5552-45a1-85b5-8977e29f3d04",
				ID:              "2acd16e5-efc9-4bff-b171-3a31c300a628",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Detail:          "03.03.2024",
				Title:           "Abrechnung",
				Filepath:        "2024-03/79bea4ff-5552-45a1-85b5-8977e29f3d04/Abrechnung.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:        "79bea4ff-5552-45a1-85b5-8977e29f3d04",
		Status:    "executed",
		Type:      transaction.TypeInterestPayout,
		AssetType: string(instrument.TypeCash),
		Credit:    20.28,
		Documents: []string{
			"2024-03/79bea4ff-5552-45a1-85b5-8977e29f3d04/Abrechnung.pdf",
		},
	},
}

func init() {
	InterestPayoutCreated02.Transaction.Timestamp, _ = internal.ParseTimestamp("2024-03-03T20:45:47.367+0000")
	InterestPayoutCreated02.CSVEntry.Timestamp = internal.DateTime{Time: InterestPayoutCreated02.Transaction.Timestamp}

	RegisterSupported("InterestPayoutCreated02", InterestPayoutCreated02)
}
