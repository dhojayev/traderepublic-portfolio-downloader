package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/instrument"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
)

var InterestPayoutCreated01 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`{
		"items": 
			[
				{
					"action": {
						"payload": "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
						"type": "timelineDetail"
					},
					"amount": {
						"currency": "EUR",
						"fractionDigits": 2,
						"value": 0.07
					},
					"badge": null,
					"eventType": "INTEREST_PAYOUT_CREATED",
					"icon": "logos/timeline_interest_new/v2",
					"id": "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
					"status": "EXECUTED",
					"subAmount": null,
					"subtitle": "2,00% p.a.",
					"timestamp": "2023-11-06T11:22:52.544+0000",
					"title": "Zinsen"
				}
			]
		}`),
		Unmarshalled: transactions.ResponseItem{
			Action: transactions.ResponseItemAction{
				Payload: "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
				Type:    "timelineDetail",
			},
			Amount: transactions.ResponseItemAmount{
				Currency:       "EUR",
				FractionDigits: 2,
				Value:          0.07,
			},
			EventType: "INTEREST_PAYOUT_CREATED",
			Icon:      "logos/timeline_interest_new/v2",
			ID:        "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
			Status:    "EXECUTED",
			Subtitle:  "2,00% p.a.",
			Timestamp: "2023-11-06T11:22:52.544+0000",
			Title:     "Zinsen",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
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
	  }`),
		Normalized: details.NormalizedResponse{
			ID: "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
			Header: details.NormalizedResponseHeaderSection{
				Data: details.NormalizedResponseHeaderSectionData{
					Icon:      "logos/timeline_interest_new/v2",
					Status:    "executed",
					Timestamp: "2023-11-06T11:22:52.544+0000",
				},
				Title: "Du hast 0,07 EUR erhalten",
				Type:  "header",
			},
			Overview: details.NormalizedResponseOverviewSection{
				NormalizedResponseTableSection: details.NormalizedResponseTableSection{
					Data: []details.NormalizedResponseTableSectionData{
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								FunctionalStyle: "EXECUTED",
								Text:            "Abgeschlossen",
								Type:            "status",
							},
							Style: "plain",
							Title: "Status",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "283,33 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Durchschnittssaldo",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "2 %",
								Type: "text",
							},
							Style: "plain",
							Title: "Jahressatz",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
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
			},
			Transaction: details.NormalizedResponseTransactionSection{
				NormalizedResponseTableSection: details.NormalizedResponseTableSection{
					Data: []details.NormalizedResponseTableSectionData{
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "+ 0,09 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Angefallen",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "0,02 €",
								Type: "text",
							},
							Style: "plain",
							Title: "Steuern",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
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
			Documents: details.NormalizedResponseDocumentsSection{
				Data: []details.NormalizedResponseDocumentsSectionData{
					{
						Action: details.NormalizedResponseSectionAction{
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
	},
	EventType: transactions.EventTypeInterestPayoutCreated,
	Transaction: transaction.Model{
		UUID:      "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
		Type:      transaction.TypeInterestPayout,
		Status:    "executed",
		Total:     0.07,
		TaxAmount: 0.02,
		Instrument: instrument.Model{
			Icon: "logos/timeline_interest_new/v2",
			Type: instrument.TypeCash,
		},
		Documents: []document.Model{
			{
				TransactionUUID: "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
				ID:              "f1b33e1e-0e44-4508-b2cd-d508715d9740",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/",
				Detail:          "06.11.2023",
				Title:           "Abrechnung",
				Filepath:        "2023-11/c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6/Abrechnung.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:        "c30c2952-ff0e-4fdb-bb8c-dfe1a8c35ce6",
		Status:    "executed",
		Type:      transaction.TypeInterestPayout,
		AssetType: string(instrument.TypeCash),
		Name:      "Savings account",
		Credit:    0.07,
		TaxAmount: 0.02,
	},
}

func init() {
	InterestPayoutCreated01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2023-11-06T11:22:52.544+0000")
	InterestPayoutCreated01.CSVEntry.Timestamp = internal.DateTime{Time: InterestPayoutCreated01.Transaction.Timestamp}

	RegisterSupported("InterestPayoutCreated01", InterestPayoutCreated01)
}
