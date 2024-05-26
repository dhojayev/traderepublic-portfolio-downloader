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

var PaymentInboundSepaDirectDebit01 = tests.TestCase{
	ResponseJSON: `{
		"id": "ddc4ed4f-0314-42cf-8a65-930da1354348",
		"sections": [
		  {
			"action": null,
			"data": {
			  "icon": "logos/timeline_plus_circle/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2023-07-23T21:05:22.543+0000"
			},
			"title": "Du hast 500,00 € per Lastschrift hinzugefügt",
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
				  "text": "Lastschrift",
				  "trend": null,
				  "type": "text"
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
				"title": "Betrag"
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
				"detail": "23.07.2023",
				"id": "cfc08704-eb56-44f1-83a0-c39aba9055ca",
				"postboxType": "PAYMENT_INBOUND_INVOICE",
				"title": "Abrechnung Einzahlung"
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
				Icon:      "logos/timeline_plus_circle/v2",
				Status:    "executed",
				Timestamp: "2023-07-23T21:05:22.543+0000",
			},
			Title: "Du hast 500,00 € per Lastschrift hinzugefügt",
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
							Text: "Lastschrift",
							Type: "text",
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
						Title: "Betrag",
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
					Detail:      "23.07.2023",
					ID:          "cfc08704-eb56-44f1-83a0-c39aba9055ca",
					PostboxType: "PAYMENT_INBOUND_INVOICE",
					Title:       "Abrechnung Einzahlung",
				},
			},
			Title: "Dokumente",
			Type:  "documents",
		},
	},
	EventType: transactions.EventTypePaymentInboundSepaDirectDebit,
	Transaction: transaction.Model{
		UUID:      "ddc4ed4f-0314-42cf-8a65-930da1354348",
		Documents: []document.Model{},
		Type:      transaction.TypeDeposit,
		Status:    "executed",
		Total:     500,
	},
	CSVEntry: filesystem.CSVEntry{
		ID:        "ddc4ed4f-0314-42cf-8a65-930da1354348",
		Status:    "executed",
		Type:      "Deposit",
		AssetType: "Other",
		Credit:    500,
	},
}

func init() {
	PaymentInboundSepaDirectDebit01.Transaction.Timestamp, _ = time.Parse(internal.DefaultTimeFormat, "2023-07-23T21:05:22.543+0000")
	PaymentInboundSepaDirectDebit01.CSVEntry.Timestamp = internal.DateTime{Time: PaymentInboundSepaDirectDebit01.Transaction.Timestamp}
}
