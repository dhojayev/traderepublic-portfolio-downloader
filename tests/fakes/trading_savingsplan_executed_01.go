package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/instrument"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
)

var TradingSavingsplanExecuted01 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`{"items": {
      "action": {
        "payload": "a8e1d370-2edd-43a3-bbfa-7894695add44",
        "type": "timelineDetail"
      },
      "amount": {
        "currency": "EUR",
        "fractionDigits": 2,
        "value": -500
      },
      "badge": null,
      "cashAccountNumber": "**********",
      "deleted": false,
      "eventType": "trading_savingsplan_executed",
      "hidden": false,
      "icon": "logos/IE00B3RBWM25/v2",
      "id": "a8e1d370-2edd-43a3-bbfa-7894695add44",
      "status": "EXECUTED",
      "subAmount": null,
      "subtitle": "Sparplan ausgeführt",
      "timestamp": "2025-02-17T13:37:21.936+0000",
      "title": "FTSE All-World USD (Dist)"
    }
	  }`),
		Unmarshalled: transactions.ResponseItem{
			Action: transactions.ResponseItemAction{
				Payload: "a8e1d370-2edd-43a3-bbfa-7894695add44",
				Type:    "timelineDetail",
			},
			Amount: transactions.ResponseItemAmount{
				Currency:       "EUR",
				FractionDigits: 2,
				Value:          -500,
			},
			EventType: "trading_savingsplan_executed",
			Icon:      "logos/IE00B3RBWM25/v2",
			ID:        "a8e1d370-2edd-43a3-bbfa-7894695add44",
			Status:    "EXECUTED",
			Subtitle:  "Sparplan ausgeführt",
			Timestamp: "2025-02-17T13:37:21.936+0000",
			Title:     "FTSE All-World USD (Dist)",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
  "id": "a8e1d370-2edd-43a3-bbfa-7894695add44",
  "sections": [
    {
      "action": {
        "payload": "IE00B3RBWM25",
        "type": "instrumentDetail"
      },
      "data": {
        "icon": "logos/IE00B3RBWM25/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2025-02-17T13:37:21.936+0000"
      },
      "title": "Du hast 500,00 € gespart",
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
            "type": "text"
          },
          "style": "plain",
          "title": "Orderart"
        },
        {
          "detail": {
            "action": null,
            "text": "FTSE All-World USD (Dist)",
            "type": "text"
          },
          "style": "plain",
          "title": "Asset"
        },
        {
          "detail": {
            "icon": "logos/bank_traderepublic/v2",
            "text": "Cash",
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
                "savingsPlanId": "a8e1d370-2edd-43a3-bbfa-7894695add44"
              },
              "type": "openSavingsPlanOverview"
            },
            "amount": "500.00 €",
            "icon": "logos/IE00B3RBWM25/v2",
            "status": "executed",
            "subtitle": "Zwei pro Monat",
            "timestamp": "2025-02-17T13:37:21.936691Z",
            "title": "FTSE All-World USD (Dist)",
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
            "text": "3,616374",
            "type": "text"
          },
          "style": "plain",
          "title": "Aktien"
        },
        {
          "detail": {
            "action": null,
            "text": "138.26 €",
            "type": "text"
          },
          "style": "plain",
          "title": "Aktienkurs"
        },
        {
          "detail": {
            "action": null,
            "text": "Kostenlos",
            "type": "text"
          },
          "style": "plain",
          "title": "Gebühr"
        },
        {
          "detail": {
            "action": null,
            "text": "500.00 €",
            "type": "text"
          },
          "style": "plain",
          "title": "Summe"
        }
      ],
      "title": "Transaktion",
      "type": "table"
    },
    {
      "data": [
        {
          "action": {
            "payload": "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
            "type": "browserModal"
          },
          "id": "f8f13202-05a5-4d00-a2ed-d22069e746bf",
          "postboxType": "SECURITIES_SETTLEMENT_SAVINGS_PLAN",
          "title": "Abrechnungsausführung"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    },
    {
      "action": null,
      "data": [
        {
          "detail": {
            "action": {
              "payload": {
                "contextCategory": "NHC",
                "contextParams": {
                  "chat_flow_key": "NHC_0029_wealth_Failed_savings_plan_execution",
                  "primId": "N***********",
                  "savingsPlanId": "********-****-****-****-************",
                  "timelineEventId": "********-****-****-****-************"
                }
              },
              "type": "customerSupportChat"
            },
            "icon": "",
            "style": "highlighted",
            "type": "listItemAvatarDefault"
          },
          "style": "plain",
          "title": ""
        }
      ],
      "title": "",
      "type": "table"
    }
  ]
}`),
	},
	EventType: transactions.EventTypeTradingSavingsPlanExecuted,
	Transaction: transaction.Model{
		UUID: "a8e1d370-2edd-43a3-bbfa-7894695add44",
		Instrument: instrument.Model{
			ISIN: "IE00B3RBWM25",
			Name: "FTSE All-World USD (Dist)",
			Icon: "logos/IE00B3RBWM25/v2",
			Type: instrument.TypeETF,
		},
		Type:   transaction.TypePurchase,
		Status: "executed",
		Shares: 3.616374,
		Rate:   138.26,
		Total:  500,
		Documents: []document.Model{
			{
				TransactionUUID: "a8e1d370-2edd-43a3-bbfa-7894695add44",
				ID:              "f8f13202-05a5-4d00-a2ed-d22069e746bf",
				URL:             "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Title:           "Abrechnungsausführung",
				Filepath:        "2025-02/a8e1d370-2edd-43a3-bbfa-7894695add44/Abrechnungsausführung.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "a8e1d370-2edd-43a3-bbfa-7894695add44",
		Status:     "executed",
		Type:       transaction.TypePurchase,
		AssetType:  string(instrument.TypeETF),
		Name:       "FTSE All-World USD (Dist)",
		Instrument: "IE00B3RBWM25",
		Shares:     3.616374,
		Rate:       138.26,
		Debit:      500,
		Documents: []string{
			"2025-02/a8e1d370-2edd-43a3-bbfa-7894695add44/Abrechnungsausführung.pdf",
		},
	},
}

func init() {
	TradingSavingsplanExecuted01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2025-02-17T13:37:21.936+0000")
	TradingSavingsplanExecuted01.CSVEntry.Timestamp = internal.DateTime{Time: TradingSavingsplanExecuted01.Transaction.Timestamp}

	RegisterSupported("TradingSavingsplanExecuted01", TradingSavingsplanExecuted01)
}
