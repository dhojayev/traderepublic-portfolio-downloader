package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
)

var InterestPayoutCreated01 = tests.TestCase{
	ResponseJSON: `{
		"id": "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
		"sections": [
		  {
			"action": null,
			"data": {
			  "icon": "logos/timeline_interest_new/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2023-11-06T11:22:52.544+0000"
			},
			"title": "Du hast 0,07 EUR erhalten",
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
				  "text": "283,33 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Durchschnittssaldo"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "2 %",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Jahressatz"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "Guthaben",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Vermögenswert"
			  }
			],
			"title": "Übersicht",
			"type": "table"
		  },
		  {
			"action": {
			  "payload": {
				"sections": [
				  {
					"action": null,
					"title": "Transaktion",
					"type": "title"
				  },
				  {
					"action": null,
					"data": [
					  {
						"detail": {
						  "action": null,
						  "text": "0,02 €",
						  "trend": null,
						  "type": "text"
						},
						"style": "plain",
						"title": "Kapitalertragsteuer"
					  },
					  {
						"detail": {
						  "action": null,
						  "text": "0,02 €",
						  "trend": null,
						  "type": "text"
						},
						"style": "highlighted",
						"title": "Steuern"
					  }
					],
					"title": null,
					"type": "table"
				  },
				  {
					"action": null,
					"data": [
					  {
						"detail": {
						  "action": null,
						  "text": "+ 0,09 €",
						  "trend": null,
						  "type": "text"
						},
						"style": "plain",
						"title": "Angefallen"
					  },
					  {
						"detail": {
						  "action": null,
						  "text": "0,02 €",
						  "trend": null,
						  "type": "text"
						},
						"style": "plain",
						"title": "Steuern"
					  },
					  {
						"detail": {
						  "action": null,
						  "text": "+ 0,07 €",
						  "trend": null,
						  "type": "text"
						},
						"style": "highlighted",
						"title": "Gesamt"
					  }
					],
					"title": null,
					"type": "table"
				  }
				]
			  },
			  "type": "infoPage"
			},
			"data": [
			  {
				"detail": {
				  "action": null,
				  "text": "+ 0,09 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Angefallen"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "0,02 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Steuern"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "+ 0,07 €",
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
				"detail": "06.11.2023",
				"id": "f1b33e1e-0e44-4508-b2cd-d508715d9740",
				"postboxType": "INTEREST_PAYOUT_INVOICE",
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
				Icon:      "logos/timeline_interest_new/v2",
				Status:    "executed",
				Timestamp: "2023-11-06T11:22:52.544+0000",
			},
			Title: "Du hast 0,07 EUR erhalten",
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
							Text: "283,33 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Durchschnittssaldo",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "2 %",
							Type: "text",
						},
						Style: "plain",
						Title: "Jahressatz",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "Guthaben",
							Type: "text",
						},
						Style: "plain",
						Title: "Vermögenswert",
					},
				},
				Title: "Übersicht",
				Type:  "table",
			},
			{
				Data: []details.ResponseSectionTypeTableData{
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "+ 0,09 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Angefallen",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "0,02 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Steuern",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "+ 0,07 €",
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
					Detail:      "06.11.2023",
					ID:          "f1b33e1e-0e44-4508-b2cd-d508715d9740",
					PostboxType: "INTEREST_PAYOUT_INVOICE",
					Title:       "Abrechnung",
				},
			},
			Title: "Dokumente",
			Type:  "documents",
		},
	},
	EventType: transactions.EventTypeInterestPayoutCreated,
	Transaction: transaction.Model{
		UUID:       "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
		Instrument: transaction.Instrument{},
		Type:       "",
		Status:     "executed",
		Yield:      0,
		Profit:     0,
		Shares:     0,
		Rate:       0,
		Commission: 0,
		Total:      0,
		TaxAmount:  0,
		Documents: []document.Model{
			{
				ID:    "",
				URL:   "",
				Date:  "",
				Title: "",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:             "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
		Status:         "executed",
		Type:           "",
		AssetType:      "",
		Name:           "",
		Instrument:     "",
		Shares:         0,
		Rate:           0,
		Yield:          0,
		Profit:         0,
		Commission:     0,
		Debit:          0,
		Credit:         0,
		TaxAmount:      0,
		InvestedAmount: 0,
	},
}
