package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
)

var PaymentInbound01 = tests.TestCase{
	TimelineDetailsData: tests.TimelineDetailsData{
		Raw: `{
		"id": "1ae661c0-b3f1-4a81-a909-79567161b014",
		"sections": [
		  {
			"action": null,
			"data": {
			  "icon": "logos/timeline_plus_circle/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2023-05-21T08:25:53.360+0000"
			},
			"title": "Du hast 200,00 € erhalten",
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
				  "text": "John Doe",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Von"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "DE78 0000 0000 0000 0000 00",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "IBAN"
			  }
			],
			"title": "Übersicht",
			"type": "table"
		  }
		]
	  }`,
		Unmarshalled: tests.TimelineDetailsResponseSections{
			Header: details.ResponseSectionTypeHeader{
				Data: details.ResponseSectionTypeHeaderData{
					Icon:      "logos/timeline_plus_circle/v2",
					Status:    "executed",
					Timestamp: "2023-05-21T08:25:53.360+0000",
				},
				Title: "Du hast 200,00 € erhalten",
				Type:  "header",
			},
			Table: details.ResponseSectionsTypeTable{
				{
					Data: []details.ResponseSectionTypeTableData{
						{
							Detail: details.ResponseSectionTypeTableDataDetail{
								FunctionalStyle: "EXECUTED",
								Text:            "Abgeschlossen",
								Type:            "status",
							},
							Style: "plain",
							Title: "Status",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetail{
								Text: "John Doe",
								Type: "text",
							},
							Style: "plain",
							Title: "Von",
						},
						{
							Detail: details.ResponseSectionTypeTableDataDetail{
								Text: "DE78 0000 0000 0000 0000 00",
								Type: "text",
							},
							Style: "plain",
							Title: "IBAN",
						},
					},
					Title: "Übersicht",
					Type:  "table",
				},
			},
		},
	},
	EventType: transactions.EventTypePaymentInbound,
	Transaction: transaction.Model{
		UUID:   "1ae661c0-b3f1-4a81-a909-79567161b014",
		Type:   transaction.TypeDeposit,
		Status: "executed",
		Total:  200,
	},
	CSVEntry: filesystem.CSVEntry{
		ID:        "1ae661c0-b3f1-4a81-a909-79567161b014",
		Status:    "executed",
		Type:      "Deposit",
		AssetType: "Other",
		Credit:    200,
	},
}

func init() {
	PaymentInbound01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2023-05-21T08:25:53.360+0000")
	PaymentInbound01.CSVEntry.Timestamp = internal.DateTime{Time: PaymentInbound01.Transaction.Timestamp}
}
