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

var CardSuccessfulTransaction01 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`[
{
}
]`)},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
    	"id": "f729b13a-ed08-5e48-bdde-87f17c478e48",
		"sections": [
			{
				"data": {
					"icon": "merchant-logos/a95de37b-f62d-42fe-b226-0d3bcae8df0c",
					"status": "executed",
					"timestamp": "2024-04-07T18:03:33.954+0000"
				},
				"title": "Du hast 2,00 € ausgegeben",
				"type": "header"
			},
			{
				"data": [
					{
						"detail": {
							"functionalStyle": "EXECUTED",
							"text": "Fertig",
							"type": "status"
						},
						"style": "plain",
						"title": "Status"
					},
					{
						"detail": {
							"icon": "logos/bank_traderepublic/v2",
							"text": "··1234",
							"type": "iconWithText"
						},
						"style": "plain",
						"title": "Zahlung"
					},
					{
						"detail": {
							"text": "Stayery",
							"type": "text"
						},
						"style": "plain",
						"title": "Händler"
					}
				],
				"title": "Übersicht",
				"type": "table"
			},
			{
				"data": [
					{
						"detail": {
							"action": {
								"type": "benefitsSavebackOverview"
							},
							"amount": "0,02 €",
							"icon": "logos/IE0031442068/v2",
							"status": "executed",
							"subtitle": "Saveback",
							"timestamp": "2024-04-07T18:03:34.802+0000",
							"title": "Core S\u0026P 500 USD (Dist)",
							"type": "embeddedTimelineItem"
						},
						"style": "plain",
						"title": "Core S\u0026P 500 USD (Dist)"
					}
				],
				"title": "Vorteile",
				"type": "table"
			},
			{
				"data": [
					{
						"detail": {
							"action": {
								"payload": {
									"contextCategory": "card-dispute",
									"contextParams": {
										"card-dispute-txId": "f729b13a-ed08-5e48-bdde-87f17c478e48"
									},
									"transactionId": "f729b13a-ed08-5e48-bdde-87f17c478e48"
								},
								"type": "customerSupportChat"
							},
							"icon": "logos/timeline_communication/v2",
							"type": "listItemAvatarDefault"
						},
						"style": "highlighted",
						"title": "Problem melden"
					}
				],
				"title": "Hilfe",
				"type": "table"
			}
		]
}`),
		Normalized: details.NormalizedResponse{
			ID: "f729b13a-ed08-5e48-bdde-87f17c478e48",
			Header: details.NormalizedResponseHeaderSection{
				Data: details.NormalizedResponseHeaderSectionData{
					Icon:      "merchant-logos/a95de37b-f62d-42fe-b226-0d3bcae8df0c",
					Status:    "executed",
					Timestamp: "2024-04-07T18:03:33.954+0000",
				},
				Title: "Du hast 2,00 € ausgegeben",
				Type:  "header",
			},
		},
	},
	EventType: transactions.EventTypeCardSuccessfulTransaction,
	Transaction: transaction.Model{
		UUID:   "f729b13a-ed08-5e48-bdde-87f17c478e48",
		Type:   transaction.TypeCardPaymentTransaction,
		Status: "executed",
		Total:  2,
		Instrument: instrument.Model{
			Name: "Stayery",
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:    "f729b13a-ed08-5e48-bdde-87f17c478e48",
		Debit: 2,
		Type:  transaction.TypeCardPaymentTransaction,
	},
}

func init() {
	CardSuccessfulTransaction01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2024-04-07T18:03:33.954+0000")
	CardSuccessfulTransaction01.CSVEntry.Timestamp = internal.DateTime{Time: CardSuccessfulTransaction01.Transaction.Timestamp}
	RegisterSupported("CardSuccessfulTransaction01", CardSuccessfulTransaction01)
}
