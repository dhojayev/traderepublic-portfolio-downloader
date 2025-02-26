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

var InterestPayout01 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`{
						"action": {
							"payload": "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
							"type": "timelineDetail"
						},
						"amount": {
							"currency": "EUR",
							"fractionDigits": 2,
							"value": 40.55
						},
						"badge": null,
						"eventType": "INTEREST_PAYOUT",
						"icon": "logos/timeline_interest_new/v2",
						"id": "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
						"status": "EXECUTED",
						"subAmount": null,
						"subtitle": "3,25 % p.a.",
						"timestamp": "2024-12-01T00:57:13.464+0000",
						"title": "Zinsen"
					}`),
		Unmarshalled: transactions.ResponseItem{
			Action: transactions.ResponseItemAction{
				Payload: "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
				Type:    "timelineDetail",
			},
			Amount: transactions.ResponseItemAmount{
				Currency:       "EUR",
				FractionDigits: 2,
				Value:          40.55,
			},
			EventType: "INTEREST_PAYOUT",
			Icon:      "logos/timeline_interest_new/v2",
			ID:        "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
			Status:    "EXECUTED",
			Subtitle:  "3,25 % p.a.",
			Timestamp: "2024-12-01T00:57:13.464+0000",
			Title:     "Zinsen",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
  "id": "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
  "sections": [
    {
      "data": {
        "icon": "logos/timeline_interest_new/v2",
        "status": "executed",
        "timestamp": "2024-12-01T01:57:20.170969+01:00"
      },
      "title": "Du hast €40.55 erhalten",
      "type": "header"
    },
    {
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
            "text": "€20,617.87",
            "type": "text"
          },
          "style": "plain",
          "title": "Durchschnittssaldo"
        },
        {
          "detail": {
            "text": "3.25 %",
            "type": "text"
          },
          "style": "plain",
          "title": "Jährliche Rate"
        },
        {
          "detail": {
            "text": "Cash",
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
      "steps": [
        {
          "content": {
            "cta": null,
            "subtitle": null,
            "timestamp": "2024-12-01T01:57:20.170969+01:00",
            "title": "Zinsberechnung"
          },
          "leading": {
            "avatar": {
              "status": "completed",
              "type": "bullet"
            },
            "connection": {
              "order": "first"
            }
          }
        },
        {
          "content": {
            "cta": null,
            "subtitle": null,
            "timestamp": "2024-12-01T01:57:20.170969+01:00",
            "title": "Zinszahlung"
          },
          "leading": {
            "avatar": {
              "status": "completed",
              "type": "bullet"
            },
            "connection": {
              "order": "last"
            }
          }
        }
      ],
      "title": "Status",
      "type": "steps"
    },
    {
      "data": [
        {
          "detail": {
            "text": "€55.07",
            "type": "text"
          },
          "style": "plain",
          "title": "Angesammelt"
        },
        {
          "detail": {
            "text": "€14.52",
            "type": "text"
          },
          "style": "plain",
          "title": "Steuern"
        },
        {
          "detail": {
            "text": "€40.55",
            "type": "text"
          },
          "style": "plain",
          "title": "Gesamt"
        }
      ],
      "title": "Transaktion",
      "type": "table"
    },
    {
      "data": [
        {
          "action": {
            "payload": "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
            "type": "browserModal"
          },
          "id": "0432d0d3-7f06-4e5b-bf54-76068c32dce3",
          "postboxType": "INTEREST_PAYOUT_INVOICE",
          "title": "Abrechnung"
        }
      ],
      "title": "Dokument",
      "type": "documents"
    }
  ]
}`),
		Normalized: details.NormalizedResponse{
			ID: "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
			Header: details.NormalizedResponseHeaderSection{
				Data: details.NormalizedResponseHeaderSectionData{
					Icon:      "logos/timeline_interest_new/v2",
					Status:    "executed",
					Timestamp: "2024-12-01T01:57:20.170969+01:00",
				},
				Title: "Du hast €40.55 erhalten",
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
								Text: "€20,617.87",
								Type: "text",
							},
							Style: "plain",
							Title: "Durchschnittssaldo",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "3.25 %",
								Type: "text",
							},
							Style: "plain",
							Title: "Jährliche Rate",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "Cash",
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
								Text: "€55.07",
								Type: "text",
							},
							Style: "plain",
							Title: "Angesammelt",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "€14.52",
								Type: "text",
							},
							Style: "plain",
							Title: "Steuern",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
								Text: "€40.55",
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
							Payload: "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
							Type:    "browserModal",
						},
						ID:          "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
						PostboxType: "INTEREST_PAYOUT_INVOICE",
						Title:       "Abrechnung",
					},
				},
				Title: "Dokumente",
				Type:  "documents",
			},
		},
	},
	EventType: transactions.EventTypeInterestPayout,
	Transaction: transaction.Model{
		UUID:      "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
		Type:      transaction.TypeInterestPayout,
		Status:    "executed",
		Total:     40.55,
		TaxAmount: 14.52,
		Instrument: instrument.Model{
			Icon: "logos/timeline_interest_new/v2",
			Type: instrument.TypeCash,
			Name: "Cash",
		},
		Documents: []document.Model{
			{
				TransactionUUID: "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
				ID:              "0432d0d3-7f06-4e5b-bf54-76068c32dce3",
				URL:             "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Title:           "Abrechnung",
				Filepath:        "2024-12/4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e/Abrechnung.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:        "4b33616d-1f9b-4e84-a88e-6dd12cdc0b7e",
		Status:    "executed",
		Type:      transaction.TypeInterestPayout,
		AssetType: string(instrument.TypeCash),
		Name:      "Savings account",
		Credit:    40.55,
		TaxAmount: 14.52,
	},
}

func init() {
	InterestPayout01.Transaction.Timestamp, _ = internal.ParseTimestamp("2024-12-01T01:57:20.170969+01:00")
	InterestPayout01.CSVEntry.Timestamp = internal.DateTime{Time: InterestPayout01.Transaction.Timestamp}

	RegisterSupported("InterestPayout01", InterestPayout01)
}
