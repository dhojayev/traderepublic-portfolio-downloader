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

var Credit01 = TransactionTestCase{
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
		"id": "23cf72a9-3888-4918-898c-c3bc38346ba1",
		"sections": [
		  {
			"action": null,
			"data": {
			  "icon": "logos/IE00BK1PV551/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2023-12-13T12:44:28.857+0000"
			},
			"title": "Du hast 2,94 € erhalten",
			"type": "header"
		  },
		  {
			"action": null,
			"data": [
			  {
				"detail": {
				  "action": null,
				  "text": "Ausschüttung",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Ereignis"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "MSCI World USD (Dist)",
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
				  "text": "10,344033",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Anteile"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "0,28 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Dividende je Aktie"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "+ 2,94 €",
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
				"detail": "13.12.2023",
				"id": "df244c67-8907-4365-bb89-ce26e1fadea5",
				"postboxType": "INCOME",
				"title": "Abrechnung"
			  }
			],
			"title": "Dokumente",
			"type": "documents"
		  }
		]
	  }`),
		Normalized: details.NormalizedResponse{
			ID: "23cf72a9-3888-4918-898c-c3bc38346ba1",
			Header: details.NormalizedResponseHeaderSection{
				Action: details.NormalizedResponseSectionAction{
					Payload: nil,
					Type:    "",
				},
				Data: details.NormalizedResponseHeaderSectionData{
					Icon:      "logos/IE00BK1PV551/v2",
					Status:    "executed",
					Timestamp: "2023-12-13T12:44:28.857+0000",
				},
				Title: "Du hast 2,94 € erhalten",
				Type:  "header",
			},
			Overview: details.NormalizedResponseOverviewSection{
				NormalizedResponseTableSection: details.NormalizedResponseTableSection{
					Data: []details.NormalizedResponseTableSectionData{
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "Ausschüttung",
								Type: "text",
							},
							Style: "plain",
							Title: "Ereignis",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "MSCI World USD (Dist)",
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
								Text: "10,344033",
								Type: "text",
							},
							Style: "plain",
							Title: "Anteile",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "0,28 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Dividende je Aktie",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "+ 2,94 €",
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
						Detail:      "13.12.2023",
						ID:          "df244c67-8907-4365-bb89-ce26e1fadea5",
						PostboxType: "INCOME",
						Title:       "Abrechnung",
					},
				},
				Title: "Dokumente",
				Type:  "documents",
			},
		},
	},
	EventType: transactions.EventTypeCredit,
	Transaction: transaction.Model{
		UUID: "23cf72a9-3888-4918-898c-c3bc38346ba1",
		Instrument: transaction.Instrument{
			ISIN: "IE00BK1PV551",
			Name: "MSCI World USD (Dist)",
			Icon: "logos/IE00BK1PV551/v2",
		},
		Documents: []document.Model{
			{
				TransactionUUID: "23cf72a9-3888-4918-898c-c3bc38346ba1",
				ID:              "df244c67-8907-4365-bb89-ce26e1fadea5",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail:          "13.12.2023",
				Title:           "Abrechnung",
				Filepath:        "2023-12/23cf72a9-3888-4918-898c-c3bc38346ba1/Abrechnung.pdf",
			},
		},
		Type:   transaction.TypeDividendPayout,
		Status: "executed",
		Shares: 10.344033,
		Rate:   0.28,
		Total:  2.94,
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "23cf72a9-3888-4918-898c-c3bc38346ba1",
		Status:     "executed",
		Type:       "Dividends",
		AssetType:  "ETF",
		Name:       "MSCI World USD (Dist)",
		Instrument: "IE00BK1PV551",
		Shares:     10.344033,
		Rate:       0.28,
		Profit:     2.94,
		Credit:     2.94,
	},
}

func init() {
	Credit01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2023-12-13T12:44:28.857+0000")
	Credit01.CSVEntry.Timestamp = internal.DateTime{Time: Credit01.Transaction.Timestamp}

	RegisterSupported("Credit01", Credit01)
}
