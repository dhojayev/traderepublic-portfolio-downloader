package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
)

var Credit01 = tests.TestCase{
	ResponseJSON: `{
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
	  }`,
	Response: tests.Response{
		HeaderSection: details.ResponseSectionTypeHeader{
			Data: details.ResponseSectionTypeHeaderData{
				Icon:      "logos/IE00BK1PV551/v2",
				Status:    "executed",
				Timestamp: "2023-12-13T12:44:28.857+0000",
			},
			Title: "Du hast 2,94 € erhalten",
			Type:  "header",
		},
		TableSections: details.ResponseSectionsTypeTable{
			{
				Data: []details.ResponseSectionTypeTableData{
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "Ausschüttung",
							Type: "text",
						},
						Style: "plain",
						Title: "Ereignis",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
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
			{
				Data: []details.ResponseSectionTypeTableData{
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "10,344033",
							Type: "text",
						},
						Style: "plain",
						Title: "Anteile",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "0,28 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Dividende je Aktie",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
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
		DocumentsSection: details.ResponseSectionTypeDocuments{
			Data: []details.ResponseSectionTypeDocumentData{
				{
					Action: details.ResponseAction{
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
	EventType: transactions.EvenTypeCredit,
}
