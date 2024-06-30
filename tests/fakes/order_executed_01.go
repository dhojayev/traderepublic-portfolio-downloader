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

var (
	OrderExecuted01 = TransactionTestCase{
		TimelineDetailsData: TimelineDetailsTestData{
			Raw: []byte(`{
			"id": "b20e367c-5542-4fab-9fd6-6faa5e7ab582",
			"sections": [
			  {
				"action": {
				  "payload": "DE000SH0MW59",
				  "type": "instrumentDetail"
				},
				"data": {
				  "icon": "logos/FR0003500008/v2",
				  "status": "executed",
				  "subtitleText": null,
				  "timestamp": "2022-03-29T09:43:31.570+0000"
				},
				"title": "Du hast 395,80 €  investiert",
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
					  "text": "CAC",
					  "trend": null,
					  "type": "text"
					},
					"style": "plain",
					"title": "Basiswert"
				  },
				  {
					"detail": {
					  "action": null,
					  "text": "Short Faktor Optionsschein 2",
					  "trend": null,
					  "type": "text"
					},
					"style": "plain",
					"title": "Produkt"
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
					  "text": "40",
					  "trend": null,
					  "type": "text"
					},
					"style": "plain",
					"title": "Anteile"
				  },
				  {
					"detail": {
					  "action": null,
					  "text": "9,87 €",
					  "trend": null,
					  "type": "text"
					},
					"style": "plain",
					"title": "Aktienkurs"
				  },
				  {
					"detail": {
					  "action": null,
					  "text": "1,00 €",
					  "trend": null,
					  "type": "text"
					},
					"style": "plain",
					"title": "Gebühr"
				  },
				  {
					"detail": {
					  "action": null,
					  "text": "395,80 €",
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
					  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/",
					  "type": "browserModal"
					},
					"detail": "29.03.2022",
					"id": "46e92aa7-df44-4a69-957c-183459753e66",
					"postboxType": "SECURITIES_SETTLEMENT",
					"title": "Abrechnung"
				  },
				  {
					"action": {
					  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/",
					  "type": "browserModal"
					},
					"detail": "29.03.2022",
					"id": "3c4ccef3-249d-4d10-a54a-18a82fb9475a",
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
		EventType: transactions.EventTypeOrderExecuted,
		Transaction: transaction.Model{
			UUID: "b20e367c-5542-4fab-9fd6-6faa5e7ab582",
			Instrument: transaction.Instrument{
				ISIN: "DE000SH0MW59",
				Name: "CAC",
				Icon: "logos/FR0003500008/v2",
			},
			Type:       transaction.TypePurchase,
			Status:     "executed",
			Shares:     40,
			Rate:       9.87,
			Commission: 1,
			Total:      395.80,
			Documents: []document.Model{
				{
					TransactionUUID: "b20e367c-5542-4fab-9fd6-6faa5e7ab582",
					ID:              "46e92aa7-df44-4a69-957c-183459753e66",
					URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/",
					Detail:          "29.03.2022",
					Title:           "Abrechnung",
					Filepath:        "2022-03/b20e367c-5542-4fab-9fd6-6faa5e7ab582/Abrechnung.pdf",
				},
				{
					TransactionUUID: "b20e367c-5542-4fab-9fd6-6faa5e7ab582",
					ID:              "3c4ccef3-249d-4d10-a54a-18a82fb9475a",
					URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/",
					Detail:          "29.03.2022",
					Title:           "Kosteninformation",
					Filepath:        "2022-03/b20e367c-5542-4fab-9fd6-6faa5e7ab582/Kosteninformation.pdf",
				},
			},
		},
		CSVEntry: filesystem.CSVEntry{
			ID:         "b20e367c-5542-4fab-9fd6-6faa5e7ab582",
			Status:     "executed",
			Type:       transaction.TypePurchase,
			AssetType:  transaction.InstrumentTypeOther,
			Name:       "CAC",
			Instrument: "DE000SH0MW59",
			Shares:     40,
			Rate:       9.87,
			Commission: 1,
			Debit:      395.80,
		},
	}
)

func init() {
	OrderExecuted01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2022-03-29T09:43:31.570+0000")
	OrderExecuted01.CSVEntry.Timestamp = internal.DateTime{Time: OrderExecuted01.Transaction.Timestamp}

	RegisterSupported("OrderExecuted01", OrderExecuted01)
}
