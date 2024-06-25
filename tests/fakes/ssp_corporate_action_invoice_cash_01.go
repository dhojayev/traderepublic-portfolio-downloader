package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

var SSPCorporateActionInvoiceCash01 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`{
		"items": [
				{
					"action": {
						"payload": "c578258f-563d-49f8-85fc-fb652d8354d9",
						"type": "timelineDetail"
					},
					"amount": {
						"currency": "EUR",
						"fractionDigits": 2,
						"value": 0.27
					},
					"badge": null,
					"eventType": "ssp_corporate_action_invoice_cash",
					"icon": "logos/DE000A0F5UF5/v2",
					"id": "c578258f-563d-49f8-85fc-fb652d8354d9",
					"status": "EXECUTED",
					"subAmount": null,
					"subtitle": "Bardividende",
					"timestamp": "2024-06-01T06:22:43.505+0000",
					"title": "NASDAQ100 USD (Dist)"
				}
		]
}`),
		Unmarshalled: transactions.ResponseItem{
			Action: transactions.ResponseItemAction{
				Payload: "c578258f-563d-49f8-85fc-fb652d8354d9",
				Type:    "timelineDetail",
			},
			Amount: transactions.ResponseItemAmount{
				Currency:       "EUR",
				FractionDigits: 2,
				Value:          0.27,
			},
			EventType: "ssp_corporate_action_invoice_cash",
			Icon:      "logos/DE000A0F5UF5/v2",
			ID:        "c578258f-563d-49f8-85fc-fb652d8354d9",
			Status:    "EXECUTED",
			Subtitle:  "Bardividende",
			Timestamp: "2024-06-01T06:22:43.505+0000",
			Title:     "NASDAQ100 USD (Dist)",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
  "id": "c578258f-563d-49f8-85fc-fb652d8354d9",
  "sections": [
    {
      "action": null,
      "data": {
        "icon": "logos/DE000A0F5UF5/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2024-06-01T06:22:43.505+0000"
      },
      "title": "Du hast 0,27 € erhalten",
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
            "text": "Bardividende",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Event"
        },
        {
          "detail": {
            "action": null,
            "text": "NASDAQ100 USD (Dist)",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Wertpapier"
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
            "text": "1.203461",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Aktien"
        },
        {
          "detail": {
            "action": null,
            "text": "0,25 $",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Dividende pro Aktie"
        },
        {
          "detail": {
            "action": null,
            "text": "0,00 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Steuer"
        },
        {
          "detail": {
            "action": null,
            "text": "0,27 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Gesamt"
        }
      ],
      "title": "Geschäft",
      "type": "table"
    },
    {
      "action": null,
      "data": [
        {
          "action": {
            "payload": "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox/2024/6/01/84928492849829348/9043285923849489238498239489.pdf",
            "type": "browserModal"
          },
          "detail": "01.06.2024",
          "id": "31987973-143c-41e5-aef5-d20812b8912f",
          "postboxType": "CA_INCOME_INVOICE",
          "title": "Dokumente"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    }
  ]
}`),
	},
	EventType: transactions.EventTypeSSPCorporateActionInvoiceCash,
	Transaction: transaction.Model{
		UUID: "c578258f-563d-49f8-85fc-fb652d8354d9",
		Instrument: transaction.Instrument{
			ISIN: "DE000A0F5UF5",
			Name: "NASDAQ100 USD (Dist)",
			Icon: "logos/DE000A0F5UF5/v2",
		},
		Type:   transaction.TypeDividendPayout,
		Status: "executed",
		Shares: 1.203461,
		Rate:   0.25,
		Total:  0.27,
		Documents: []document.Model{
			{
				TransactionUUID: "c578258f-563d-49f8-85fc-fb652d8354d9",
				ID:              "31987973-143c-41e5-aef5-d20812b8912f",
				URL:             "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox/2024/6/01/84928492849829348/9043285923849489238498239489.pdf",
				Detail:          "01.06.2024",
				Title:           "Dokumente",
				Filepath:        "2024-06/c578258f-563d-49f8-85fc-fb652d8354d9/Dokumente.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "c578258f-563d-49f8-85fc-fb652d8354d9",
		Status:     "executed",
		Type:       transaction.TypeDividendPayout,
		AssetType:  transaction.InstrumentTypeOther,
		Name:       "NASDAQ100 USD (Dist)",
		Instrument: "DE000A0F5UF5",
		Profit:     0.27,
		Shares:     1.203461,
		Rate:       0.25,
		Credit:     0.27,
	},
}

func init() {
	SSPCorporateActionInvoiceCash01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2024-06-01T06:22:43.505+0000")
	SSPCorporateActionInvoiceCash01.CSVEntry.Timestamp = internal.DateTime{Time: SSPCorporateActionInvoiceCash01.Transaction.Timestamp}

	RegisterSupported("SSPCorporateActionInvoiceCash01", SSPCorporateActionInvoiceCash01)
}
