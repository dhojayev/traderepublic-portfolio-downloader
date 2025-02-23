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

var OrderExecuted02 = TransactionTestCase{
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
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
	  }`),
		Normalized: details.NormalizedResponse{
			ID: "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
			Header: details.NormalizedResponseHeaderSection{
				Action: details.NormalizedResponseSectionAction{
					Payload: "DE000A0F5UF5",
					Type:    "instrumentDetail",
				},
				Data: details.NormalizedResponseHeaderSectionData{
					Icon:      "logos/DE000A0F5UF5/v2",
					Status:    "executed",
					Timestamp: "2023-11-23T15:45:24.252+0000",
				},
				Title: "Du hast 136,14 €  investiert",
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
								Text: "Kauf",
								Type: "text",
							},
							Style: "plain",
							Title: "Orderart",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
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
			},
			Transaction: details.NormalizedResponseTransactionSection{
				NormalizedResponseTableSection: details.NormalizedResponseTableSection{
					Data: []details.NormalizedResponseTableSectionData{
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "1",
								Type: "text",
							},
							Style: "plain",
							Title: "Anteile",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "135,14 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Aktienkurs",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "1,00 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Gebühr",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
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
			Documents: details.NormalizedResponseDocumentsSection{
				Data: []details.NormalizedResponseDocumentsSectionData{
					{
						Action: details.NormalizedResponseSectionAction{
							Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
							Type:    "browserModal",
						},
						Detail:      "23.11.2023",
						ID:          "c9a1c524-1c54-4689-8b2f-0f1bcbb91c9d",
						PostboxType: "SECURITIES_SETTLEMENT",
						Title:       "Abrechnung",
					},
					{
						Action: details.NormalizedResponseSectionAction{
							Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
							Type:    "browserModal",
						},
						Detail:      "23.11.2023",
						ID:          "b26233a9-ee80-4da9-8404-08e722fe830b",
						PostboxType: "INFO",
						Title:       "Basisinformationsblatt",
					},
					{
						Action: details.NormalizedResponseSectionAction{
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
	},
	EventType: transactions.EventTypeOrderExecuted,
	Transaction: transaction.Model{
		UUID: "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
		Instrument: instrument.Model{
			ISIN: "DE000A0F5UF5",
			Name: "NASDAQ100 USD (Dist)",
			Icon: "logos/DE000A0F5UF5/v2",
			Type: instrument.TypeETF,
		},
		Type:       transaction.TypePurchase,
		Status:     "executed",
		Shares:     1,
		Rate:       135.14,
		Commission: 1,
		Total:      136.14,
		Documents: []document.Model{
			{
				TransactionUUID: "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
				ID:              "c9a1c524-1c54-4689-8b2f-0f1bcbb91c9d",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail:          "23.11.2023",
				Title:           "Abrechnung",
				Filepath:        "2023-11/1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad/Abrechnung.pdf",
			},
			{
				TransactionUUID: "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
				ID:              "b26233a9-ee80-4da9-8404-08e722fe830b",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail:          "23.11.2023",
				Title:           "Basisinformationsblatt",
				Filepath:        "2023-11/1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad/Basisinformationsblatt.pdf",
			},
			{
				TransactionUUID: "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
				ID:              "b582015c-7a5c-47d0-8d33-6391d414cdc7",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail:          "23.11.2023",
				Title:           "Kosteninformation",
				Filepath:        "2023-11/1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad/Kosteninformation.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "1d9ad3b5-e65c-41f6-9c7d-96baa2a2ecad",
		Status:     "executed",
		Type:       transaction.TypePurchase,
		AssetType:  string(instrument.TypeETF),
		Name:       "NASDAQ100 USD (Dist)",
		Instrument: "DE000A0F5UF5",
		Shares:     1,
		Rate:       135.14,
		Commission: 1,
		Debit:      136.14,
	},
}

func init() {
	OrderExecuted02.Transaction.Timestamp, _ = internal.ParseTimestamp("2023-11-23T15:45:24.252+0000")
	OrderExecuted02.CSVEntry.Timestamp = internal.DateTime{Time: OrderExecuted02.Transaction.Timestamp}

	RegisterSupported("OrderExecuted01", OrderExecuted02)
}
