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

var TradeInvoice01 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`{
		"items":
		[
			{
				"action": {
					"payload": "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
					"type": "timelineDetail"
				},
				"amount": {
					"currency": "EUR",
					"fractionDigits": 2,
					"value": -411.45
				},
				"badge": null,
				"cashAccountNumber": null,
				"deleted": false,
				"eventType": "TRADE_INVOICE",
				"hidden": false,
				"icon": "logos/US00206R1023/v2",
				"id": "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
				"status": "EXECUTED",
				"subAmount": null,
				"subtitle": "Kauforder",
				"timestamp": "2024-06-17T12:02:07.095+0000",
				"title": "AT\u0026T"
    	}
		]
	}`),
		Unmarshalled: transactions.ResponseItem{
			Action: transactions.ResponseItemAction{
				Payload: "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
				Type:    "timelineDetail",
			},
			Amount: transactions.ResponseItemAmount{
				Currency:       "EUR",
				FractionDigits: 2,
				Value:          -411.45,
			},
			EventType: "TRADE_INVOICE",
			Icon:      "logos/US00206R1023/v2",
			ID:        "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
			Status:    "EXECUTED",
			Subtitle:  "Kauforder",
			Timestamp: "2024-06-17T12:02:07.095+0000",
			Title:     "AT\u0026T",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
  "id": "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
  "sections": [
    {
      "action": {
        "payload": "US00206R1023",
        "type": "instrumentDetail"
      },
      "data": {
        "icon": "logos/US00206R1023/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2024-06-17T12:02:07.095+0000"
      },
      "title": "Du hast 411,45 €  investiert",
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
            "text": "Kauf",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Orderart"
        },
        {
          "detail": {
            "action": null,
            "text": "AT\u0026T",
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
            "text": "25",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Anteile"
        },
        {
          "detail": {
            "action": null,
            "text": "16,42 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Aktienkurs"
        },
        {
          "detail": {
            "action": null,
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
            "text": "411,45 €",
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
            "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/***",
            "type": "browserModal"
          },
          "detail": "17.06.2024",
          "id": "6ad6cf20-7863-4b24-96ec-2a9261c302ba",
          "postboxType": "SECURITIES_SETTLEMENT",
          "title": "Abrechnung"
        },
        {
          "action": {
            "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/***",
            "type": "browserModal"
          },
          "detail": "17.06.2024",
          "id": "52de9800-373f-4361-8ab9-478b1d01b45e",
          "postboxType": "COSTS_INFO_BUY_V2",
          "title": "Kosteninformation"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    }
  ]
}`),
	},
	EventType: "TRADE_INVOICE",
	Transaction: transaction.Model{
		UUID: "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
		Instrument: instrument.Model{
			ISIN: "US00206R1023",
			Name: "AT&T",
			Icon: "logos/US00206R1023/v2",
			Type: instrument.TypeOther,
		},
		Type:       transaction.TypePurchase,
		Status:     "executed",
		Shares:     25,
		Rate:       16.42,
		Commission: 1,
		Total:      411.45,
		Documents: []document.Model{
			{
				TransactionUUID: "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
				ID:              "6ad6cf20-7863-4b24-96ec-2a9261c302ba",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/***",
				Detail:          "17.06.2024",
				Title:           "Abrechnung",
				Filepath:        "2024-06/91aa5f02-27f0-4a9e-8733-f90dee41f2cc/Abrechnung.pdf",
			},
			{
				TransactionUUID: "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
				ID:              "52de9800-373f-4361-8ab9-478b1d01b45e",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/***",
				Detail:          "17.06.2024",
				Title:           "Kosteninformation",
				Filepath:        "2024-06/91aa5f02-27f0-4a9e-8733-f90dee41f2cc/Kosteninformation.pdf",
			},
		},
	},
	DepotTransactionCSVEntry: filesystem.DepotTransactionCSVEntry{
		ID:         "91aa5f02-27f0-4a9e-8733-f90dee41f2cc",
		Status:     "executed",
		Type:       transaction.TypePurchase,
		AssetType:  string(instrument.TypeOther),
		Name:       "AT&T",
		Instrument: "US00206R1023",
		Shares:     25,
		Rate:       16.42,
		Commission: 1,
		Debit:      411.45,
	},
}

func init() {
	TradeInvoice01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2024-06-17T12:02:07.095+0000")
	TradeInvoice01.DepotTransactionCSVEntry.Timestamp = internal.DateTime{Time: TradeInvoice01.Transaction.Timestamp}

	RegisterSupported("TradeInvoice01", TradeInvoice01)
}
