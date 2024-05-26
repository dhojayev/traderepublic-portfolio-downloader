package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
)

var PaymentInbound01 = tests.TestCase{
	ResponseJSON: `{
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
	Response: tests.Response{
		HeaderSection: details.ResponseSectionTypeHeader{
			Data: details.ResponseSectionTypeHeaderData{
				Icon:      "logos/timeline_plus_circle/v2",
				Status:    "executed",
				Timestamp: "2023-05-21T08:25:53.360+0000",
			},
			Title: "Du hast 200,00 € erhalten",
			Type:  "header",
		},
		TableSections: details.ResponseSectionsTypeTable{
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
	EventType: transactions.EventTypePaymentInbound,
}
