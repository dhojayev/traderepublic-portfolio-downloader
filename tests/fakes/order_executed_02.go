package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
)

var OrderExecuted02 = tests.TestCase{
	ResponseJSON: `{
		"id": "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
		"sections": [
		  {
			"action": {
			  "payload": "DE000A0F5UF5",
			  "type": "instrumentDetail"
			},
			"data": {
			  "icon": "logos/DE000A0F5UF5/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2023-11-23T15:45:24.252+0000"
			},
			"title": "Du hast 136,14 €  investiert",
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
				  "text": "NASDAQ100 USD (Dist)",
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
				  "text": "1",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Anteile"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "135,14 €",
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
				  "text": "136,14 €",
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
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				  "type": "browserModal"
				},
				"detail": "23.11.2023",
				"id": "c9a1c524-1c54-4689-8b2f-0f1bcbb91c9d",
				"postboxType": "SECURITIES_SETTLEMENT",
				"title": "Abrechnung"
			  },
			  {
				"action": {
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				  "type": "browserModal"
				},
				"detail": "23.11.2023",
				"id": "b26233a9-ee80-4da9-8404-08e722fe830b",
				"postboxType": "INFO",
				"title": "Basisinformationsblatt"
			  },
			  {
				"action": {
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				  "type": "browserModal"
				},
				"detail": "23.11.2023",
				"id": "b582015c-7a5c-47d0-8d33-6391d414cdc7",
				"postboxType": "COSTS_INFO_BUY_V2",
				"title": "Kosteninformation"
			  }
			],
			"title": "Dokumente",
			"type": "documents"
		  }
		]
	  }`,
	Response: tests.Response{
		HeaderSection: details.ResponseSectionTypeHeader{
			Action: details.ResponseAction{
				Payload: "DE000A0F5UF5",
				Type:    "instrumentDetail",
			},
			Data: details.ResponseSectionTypeHeaderData{
				Icon:      "logos/DE000A0F5UF5/v2",
				Status:    "executed",
				Timestamp: "2023-11-23T15:45:24.252+0000",
			},
			Title: "Du hast 136,14 €  investiert",
			Type:  "header",
		},
		TableSections: details.ResponseSectionsTypeTable{
			{
				Data: []details.ResponseSectionTypeTableData{
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							FunctionalStyle: "EXECUTED",
							Text:            "Ausgeführt",
							Type:            "status",
						},
						Style: "plain",
						Title: "Status",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "Kauf",
							Type: "text",
						},
						Style: "plain",
						Title: "Orderart",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "NASDAQ100 USD (Dist)",
							Type: "text",
						},
						Style: "plain",
						Title: "Asset",
					},
				},
				Title: "Übersicht",
				Type:  "table",
			},
			{
				Data: []details.ResponseSectionTypeTableData{
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "1",
							Type: "text",
						},
						Style: "plain",
						Title: "Anteile",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "135,14 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Aktienkurs",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "1,00 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Gebühr",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "136,14 €",
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
		DocumentsSection: details.ResponseSectionTypeDocuments{
			Data: []details.ResponseSectionTypeDocumentData{
				{
					Action: details.ResponseAction{
						Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
						Type:    "browserModal",
					},
					Detail:      "23.11.2023",
					ID:          "c9a1c524-1c54-4689-8b2f-0f1bcbb91c9d",
					PostboxType: "SECURITIES_SETTLEMENT",
					Title:       "Abrechnung",
				},
				{
					Action: details.ResponseAction{
						Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
						Type:    "browserModal",
					},
					Detail:      "23.11.2023",
					ID:          "b26233a9-ee80-4da9-8404-08e722fe830b",
					PostboxType: "INFO",
					Title:       "Basisinformationsblatt",
				},
				{
					Action: details.ResponseAction{
						Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
						Type:    "browserModal",
					},
					Detail:      "23.11.2023",
					ID:          "b582015c-7a5c-47d0-8d33-6391d414cdc7",
					PostboxType: "COSTS_INFO_BUY_V2",
					Title:       "Kosteninformation",
				},
			},
			Title: "Dokumente",
			Type:  "documents",
		},
	},
	EventType: transactions.EventTypeOrderExecuted,
	Transaction: transaction.Model{
		UUID: "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
		Instrument: transaction.Instrument{
			ISIN: "DE000A0F5UF5",
			Name: "NASDAQ100 USD (Dist)",
			Icon: "logos/DE000A0F5UF5/v2",
		},
		Type:       transaction.TypePurchase,
		Status:     "executed",
		Shares:     1,
		Rate:       135.14,
		Commission: 1,
		Total:      136.14,
		Documents: []document.Model{
			{
				ID:     "c9a1c524-1c54-4689-8b2f-0f1bcbb91c9d",
				URL:    "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail: "23.11.2023",
				Title:  "Abrechnung",
			},
			{
				ID:     "b26233a9-ee80-4da9-8404-08e722fe830b",
				URL:    "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail: "23.11.2023",
				Title:  "Basisinformationsblatt",
			},
			{
				ID:     "b582015c-7a5c-47d0-8d33-6391d414cdc7",
				URL:    "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail: "23.11.2023",
				Title:  "Kosteninformation",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
		Status:     "executed",
		Type:       "Purchase",
		AssetType:  "ETF",
		Name:       "NASDAQ100 USD (Dist)",
		Instrument: "DE000A0F5UF5",
		Shares:     1,
		Rate:       135.14,
		Commission: 1,
		Debit:      136.14,
	},
}

func init() {
	OrderExecuted02.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2023-11-23T15:45:24.252+0000")
	OrderExecuted02.CSVEntry.Timestamp = internal.DateTime{Time: OrderExecuted02.Transaction.Timestamp}

	OrderExecuted02.Transaction.Documents[0].Timestamp, _ = time.Parse(document.ResolverTimeFormat, "23.11.2023")
	OrderExecuted02.Transaction.Documents[1].Timestamp, _ = time.Parse(document.ResolverTimeFormat, "23.11.2023")
	OrderExecuted02.Transaction.Documents[2].Timestamp, _ = time.Parse(document.ResolverTimeFormat, "23.11.2023")
}
