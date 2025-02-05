package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/instrument"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
)

var PaymentInbound01 = TransactionTestCase{
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
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
	  }`),
		Normalized: details.NormalizedResponse{
			ID: "1ae661c0-b3f1-4a81-a909-79567161b014",
			Header: details.NormalizedResponseHeaderSection{
				Data: details.NormalizedResponseHeaderSectionData{
					Icon:      "logos/timeline_plus_circle/v2",
					Status:    "executed",
					Timestamp: "2023-05-21T08:25:53.360+0000",
				},
				Title: "Du hast 200,00 € erhalten",
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
								Text: "John Doe",
								Type: "text",
							},
							Style: "plain",
							Title: "Von",
						},
						{
							Detail: details.NormalizedResponseTableSectionDataDetail{
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
		Instrument: instrument.Model{
			Icon: "logos/timeline_plus_circle/v2",
			Type: instrument.TypeCash,
		},
	},
	DepotTransactionCSVEntry: filesystem.DepotTransactionCSVEntry{
		ID:        "1ae661c0-b3f1-4a81-a909-79567161b014",
		Status:    "executed",
		Type:      transaction.TypeDeposit,
		AssetType: string(instrument.TypeCash),
		Credit:    200,
	},
}

func init() {
	PaymentInbound01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2023-05-21T08:25:53.360+0000")
	PaymentInbound01.DepotTransactionCSVEntry.Timestamp = internal.DateTime{Time: PaymentInbound01.Transaction.Timestamp}

	RegisterSupported("PaymentInbound01", PaymentInbound01)
}
