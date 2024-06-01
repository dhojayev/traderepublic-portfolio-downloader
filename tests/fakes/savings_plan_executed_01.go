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

var SavingsPlanExecuted01 = tests.TestCase{
	ResponseJSON: `{
		"id": "7c9be07c-7b88-4a49-a4be-425094388b8e",
		"sections": [
		  {
			"action": {
			  "payload": "IE00BK1PV551",
			  "type": "instrumentDetail"
			},
			"data": {
			  "icon": "logos/IE00BK1PV551/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2023-11-11T13:40:59.926+0000"
			},
			"title": "Du hast 500,00 € investiert",
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
				  "text": "Sparplan",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Orderart"
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
			  },
			  {
				"detail": {
				  "icon": "logos/bank_commerzbank/v2",
				  "text": "·· 0000",
				  "type": "iconWithText"
				},
				"style": "plain",
				"title": "Zahlung"
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
				  "action": {
					"payload": {
					  "savingsPlanId": "f9c615ca-959c-4cf1-b8b9-10541673fba5"
					},
					"type": "openSavingsPlanOverview"
				  },
				  "amount": "500,00 €",
				  "icon": "logos/IE00BK1PV551/v2",
				  "status": "executed",
				  "subtitle": "Wöchentlich",
				  "timestamp": "2023-11-02T16:41:39.944Z",
				  "title": "MSCI World USD (Dist)",
				  "type": "embeddedTimelineItem"
				},
				"style": "plain",
				"title": ""
			  }
			],
			"title": "Sparplan",
			"type": "table"
		  },
		  {
			"action": null,
			"data": [
			  {
				"detail": {
				  "action": null,
				  "text": "6,887811",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Anteile"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "72,592 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Anteilspreis"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "Gratis",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Gebühr"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "500,00 €",
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
				"detail": "11.11.2023",
				"id": "0ac3aea7-6d68-4815-8f25-9c8997ef790d",
				"postboxType": "SAVINGS_PLAN_EXECUTED_V2",
				"title": "Abrechnung Ausführung"
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
				Payload: "IE00BK1PV551",
				Type:    "instrumentDetail",
			},
			Data: details.ResponseSectionTypeHeaderData{
				Icon:      "logos/IE00BK1PV551/v2",
				Status:    "executed",
				Timestamp: "2023-11-11T13:40:59.926+0000",
			},
			Title: "Du hast 500,00 € investiert",
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
							Text: "Sparplan",
							Type: "text",
						},
						Style: "plain",
						Title: "Orderart",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "MSCI World USD (Dist)",
							Type: "text",
						},
						Style: "plain",
						Title: "Asset",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Icon: "logos/bank_commerzbank/v2",
							Text: "·· 0000",
							Type: "iconWithText",
						},
						Style: "plain",
						Title: "Zahlung",
					},
				},
				Title: "Übersicht",
				Type:  "table",
			},
			{
				Data: []details.ResponseSectionTypeTableData{
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Action: details.ResponseAction{
								Payload: map[string]any{
									"savingsPlanId": "f9c615ca-959c-4cf1-b8b9-10541673fba5",
								},
								Type: "openSavingsPlanOverview",
							},
							Amount:    "500,00 €",
							Icon:      "logos/IE00BK1PV551/v2",
							Status:    "executed",
							Subtitle:  "Wöchentlich",
							Timestamp: "2023-11-02T16:41:39.944Z",
							Title:     "MSCI World USD (Dist)",
							Type:      "embeddedTimelineItem",
						},
						Style: "plain",
						Title: "",
					},
				},
				Title: "Sparplan",
				Type:  "table",
			},
			{
				Data: []details.ResponseSectionTypeTableData{
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "6,887811",
							Type: "text",
						},
						Style: "plain",
						Title: "Anteile",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "72,592 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Anteilspreis",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "Gratis",
							Type: "text",
						},
						Style: "plain",
						Title: "Gebühr",
					},
					{
						Detail: details.ResponseSectionTypeTableDataDetail{
							Text: "500,00 €",
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
					Detail:      "11.11.2023",
					ID:          "0ac3aea7-6d68-4815-8f25-9c8997ef790d",
					PostboxType: "SAVINGS_PLAN_EXECUTED_V2",
					Title:       "Abrechnung Ausführung",
				},
			},
			Title: "Dokumente",
			Type:  "documents",
		},
	},
	EventType: transactions.EventTypeSavingsPlanExecuted,
	Transaction: transaction.Model{
		UUID: "7c9be07c-7b88-4a49-a4be-425094388b8e",
		Instrument: transaction.Instrument{
			ISIN: "IE00BK1PV551",
			Name: "MSCI World USD (Dist)",
			Icon: "logos/IE00BK1PV551/v2",
		},
		Type:   transaction.TypePurchase,
		Status: "executed",
		Shares: 6.887811,
		Rate:   72.592,
		Total:  500,
		Documents: []document.Model{
			{
				ID:       "0ac3aea7-6d68-4815-8f25-9c8997ef790d",
				URL:      "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail:   "11.11.2023",
				Title:    "Abrechnung Ausführung",
				Filepath: "2023-11/7c9be07c-7b88-4a49-a4be-425094388b8e/Abrechnung Ausführung.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "7c9be07c-7b88-4a49-a4be-425094388b8e",
		Status:     "executed",
		Type:       "Purchase",
		AssetType:  "ETF",
		Name:       "MSCI World USD (Dist)",
		Instrument: "IE00BK1PV551",
		Shares:     6.887811,
		Rate:       72.592,
		Debit:      500,
	},
}

func init() {
	SavingsPlanExecuted01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2023-11-11T13:40:59.926+0000")
	SavingsPlanExecuted01.CSVEntry.Timestamp = internal.DateTime{Time: SavingsPlanExecuted01.Transaction.Timestamp}
}
