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

var OrderExecuted03 = TestCase{
	TimelineDetailsData: TimelineDetailsData{
		Raw: `{
		"id": "a3b8e625-a6e9-4269-9529-01ebb86d69bb",
		"sections": [
		  {
			"action": {
			  "payload": "US6701002056",
			  "type": "instrumentDetail"
			},
			"data": {
			  "icon": "logos/US6701002056/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2024-03-11T11:23:59.448+0000"
			},
			"title": "Du hast 482,99 €  erhalten",
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
				  "text": "Novo Nordisk (ADR)",
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
				  "text": "0,21 %",
				  "trend": "positive",
				  "type": "text"
				},
				"style": "plain",
				"title": "Rendite"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "1,04 €",
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
			"action": null,
			"data": [
			  {
				"detail": {
				  "action": null,
				  "text": "5",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Anteile"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "96,80 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Aktienkurs"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "0,01 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Steuern"
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
				  "text": "+ 482,99 €",
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
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				  "type": "browserModal"
				},
				"detail": "11.03.2024",
				"id": "f17b2237-0e32-410e-b38b-8638600ffbb0",
				"postboxType": "SECURITIES_SETTLEMENT",
				"title": "Abrechnung"
			  },
			  {
				"action": {
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				  "type": "browserModal"
				},
				"detail": "27.02.2024",
				"id": "3c214355-dc5a-488a-b780-b28fb66b66c8",
				"postboxType": "CONFIRM_ORDER_CREATE_V2",
				"title": "Auftragsbestätigung"
			  },
			  {
				"action": {
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				  "type": "browserModal"
				},
				"detail": "12.02.2024",
				"id": "21a13acc-7f3c-4156-8365-be8089006ac4",
				"postboxType": "COSTS_INFO_SELL_V2",
				"title": "Kosteninformation"
			  }
			],
			"title": "Dokumente",
			"type": "documents"
		  }
		]
	  }`,
	},
	EventType: transactions.EventTypeOrderExecuted,
	Transaction: transaction.Model{
		UUID: "a3b8e625-a6e9-4269-9529-01ebb86d69bb",
		Instrument: transaction.Instrument{
			ISIN: "US6701002056",
			Name: "Novo Nordisk (ADR)",
			Icon: "logos/US6701002056/v2",
		},
		Type:       transaction.TypeSale,
		Status:     "executed",
		Shares:     5,
		Rate:       96.80,
		Yield:      0.21,
		Profit:     1.04,
		Commission: 1,
		Total:      482.99,
		TaxAmount:  0.01,
		Documents: []document.Model{
			{
				ID:       "f17b2237-0e32-410e-b38b-8638600ffbb0",
				URL:      "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Detail:   "11.03.2024",
				Title:    "Abrechnung",
				Filepath: "2024-03/a3b8e625-a6e9-4269-9529-01ebb86d69bb/Abrechnung.pdf",
			},
			{
				ID:       "3c214355-dc5a-488a-b780-b28fb66b66c8",
				URL:      "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Detail:   "27.02.2024",
				Title:    "Auftragsbestätigung",
				Filepath: "2024-02/a3b8e625-a6e9-4269-9529-01ebb86d69bb/Auftragsbestätigung.pdf",
			},
			{
				ID:       "21a13acc-7f3c-4156-8365-be8089006ac4",
				URL:      "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Detail:   "12.02.2024",
				Title:    "Kosteninformation",
				Filepath: "2024-02/a3b8e625-a6e9-4269-9529-01ebb86d69bb/Kosteninformation.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "a3b8e625-a6e9-4269-9529-01ebb86d69bb",
		Status:     "executed",
		Type:       "Sale",
		AssetType:  "Other",
		Name:       "Novo Nordisk (ADR)",
		Instrument: "US6701002056",
		Shares:     -5,
		Rate:       96.80,
		Yield:      0.21,
		Profit:     1.04,
		Commission: 1,
		Credit:     482.99,
		TaxAmount:  0.01,
	},
}

func init() {
	OrderExecuted03.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2024-03-11T11:23:59.448+0000")
	OrderExecuted03.CSVEntry.Timestamp = internal.DateTime{Time: OrderExecuted03.Transaction.Timestamp}

	RegisterSupported(OrderExecuted03)
}
