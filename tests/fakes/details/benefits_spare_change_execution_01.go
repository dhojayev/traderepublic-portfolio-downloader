package details_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/document"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/instrument"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/portfolio/transaction"
)

var BenefitsSpareChangeExecution01 = Fake{
	RawResponse: []byte(`{
		"id": "265cb9c0-664a-45d4-b179-3061f196dd2a",
		"sections": [
		  {
			"action": {
			  "payload": "DE000A0F5UF5",
			  "type": "instrumentDetail"
			},
			"data": {
			  "icon": "logos/DE000A0F5UF5/v2",
			  "status": "executed",
			  "timestamp": "2024-01-04T12:26:52.110+0000"
			},
			"title": "Du hast 1,09 € investiert",
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
				  "text": "Round up",
				  "type": "text"
				},
				"style": "plain",
				"title": "Ordertyp"
			  },
			  {
				"detail": {
				  "text": "NASDAQ100 USD (Dist)",
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
			"data": [
			  {
				"detail": {
				  "text": "0.006882",
				  "type": "text"
				},
				"style": "plain",
				"title": "Aktien"
			  },
			  {
				"detail": {
				  "text": "158,38 €",
				  "type": "text"
				},
				"style": "plain",
				"title": "Aktienkurs"
			  },
			  {
				"detail": {
				  "text": "Kostenlos",
				  "type": "text"
				},
				"style": "plain",
				"title": "Gebühr"
			  },
			  {
				"detail": {
				  "text": "1,09 €",
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
				"id": "9df4c2e1-0de2-4900-aa8c-af5371ed58f6",
				"postboxType": "BENEFIT_DEACTIVATED",
				"title": "Deaktivierung"
			  },
			  {
				"action": {
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				  "type": "browserModal"
				},
				"id": "3a8ebf86-a2bb-463e-8bfd-28fd705359ff",
				"postboxType": "SAVINGS_PLAN_EXECUTED_V2",
				"title": "Abrechnung Ausführung"
			  },
			  {
				"action": {
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				  "type": "browserModal"
				},
				"id": "e2dfa755-e039-45c7-b7bb-1ac024844f75",
				"postboxType": "COSTS_INFO_SAVINGS_PLAN_V2",
				"title": "Kosteninformation"
			  }
			],
			"title": "Documents",
			"type": "documents"
		  }
		]
	  }`),
	Unmarshalled: details.NormalizedResponse{
		ID: "265cb9c0-664a-45d4-b179-3061f196dd2a",
		Header: details.NormalizedResponseHeaderSection{
			Action: details.NormalizedResponseSectionAction{
				Payload: "DE000A0F5UF5",
				Type:    "instrumentDetail",
			},
			Data: details.NormalizedResponseHeaderSectionData{
				Icon:      "logos/DE000A0F5UF5/v2",
				Status:    "executed",
				Timestamp: "2024-01-04T12:26:52.110+0000",
			},
			Title: "Du hast 1,09 € investiert",
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
							Text: "Round up",
							Type: "text",
						},
						Style: "plain",
						Title: "Ordertyp",
					},
					{
						Detail: details.NormalizedResponseTableSectionDataDetail{
							Text: "NASDAQ100 USD (Dist)",
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
							Text: "0.006882",
							Type: "text",
						},
						Style: "plain",
						Title: "Aktien",
					},
					{
						Detail: details.NormalizedResponseTableSectionDataDetail{
							Text: "158,38 €",
							Type: "text",
						},
						Style: "plain",
						Title: "Aktienkurs",
					},
					{
						Detail: details.NormalizedResponseTableSectionDataDetail{
							Text: "Kostenlos",
							Type: "text",
						},
						Style: "plain",
						Title: "Gebühr",
					},
					{
						Detail: details.NormalizedResponseTableSectionDataDetail{
							Text: "1,09 €",
							Type: "text",
						},
						Style: "plain",
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
					ID:          "9df4c2e1-0de2-4900-aa8c-af5371ed58f6",
					PostboxType: "BENEFIT_DEACTIVATED",
					Title:       "Deaktivierung",
				},
				{
					Action: details.NormalizedResponseSectionAction{
						Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
						Type:    "browserModal",
					},
					ID:          "3a8ebf86-a2bb-463e-8bfd-28fd705359ff",
					PostboxType: "SAVINGS_PLAN_EXECUTED_V2",
					Title:       "Abrechnung Ausführung",
				},
				{
					Action: details.NormalizedResponseSectionAction{
						Payload: "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
						Type:    "browserModal",
					},
					ID:          "e2dfa755-e039-45c7-b7bb-1ac024844f75",
					PostboxType: "COSTS_INFO_SAVINGS_PLAN_V2",
					Title:       "Kosteninformation",
				},
			},
			Title: "Documents",
			Type:  "documents",
		},
	},
	// EventType: transactions.EventTypeBenefitsSpareChangeExecution,
	Model: transaction.Model{
		UUID: "265cb9c0-664a-45d4-b179-3061f196dd2a",
		Instrument: instrument.Model{
			ISIN: "DE000A0F5UF5",
			Name: "NASDAQ100 USD (Dist)",
			Icon: "logos/DE000A0F5UF5/v2",
			Type: instrument.TypeETF,
		},
		Type:   transaction.TypeRoundUp,
		Status: "executed",
		Shares: 0.006882,
		Rate:   158.38,
		Total:  1.09,
		Documents: []document.Model{
			{
				TransactionUUID: "265cb9c0-664a-45d4-b179-3061f196dd2a",
				ID:              "9df4c2e1-0de2-4900-aa8c-af5371ed58f6",
				URL:             "https://traderepublic-postbox-platform-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Title:           "Deaktivierung",
				Filepath:        "2024-01/265cb9c0-664a-45d4-b179-3061f196dd2a/Deaktivierung.pdf",
			},
			{
				TransactionUUID: "265cb9c0-664a-45d4-b179-3061f196dd2a",
				ID:              "3a8ebf86-a2bb-463e-8bfd-28fd705359ff",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Title:           "Abrechnung Ausführung",
				Filepath:        "2024-01/265cb9c0-664a-45d4-b179-3061f196dd2a/Abrechnung Ausführung.pdf",
			},
			{
				TransactionUUID: "265cb9c0-664a-45d4-b179-3061f196dd2a",
				ID:              "e2dfa755-e039-45c7-b7bb-1ac024844f75",
				URL:             "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				Title:           "Kosteninformation",
				Filepath:        "2024-01/265cb9c0-664a-45d4-b179-3061f196dd2a/Kosteninformation.pdf",
			},
		},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:         "265cb9c0-664a-45d4-b179-3061f196dd2a",
		Status:     "executed",
		Type:       transaction.TypeRoundUp,
		AssetType:  string(instrument.TypeETF),
		Name:       "NASDAQ100 USD (Dist)",
		Instrument: "DE000A0F5UF5",
		Shares:     0.006882,
		Rate:       158.38,
		Debit:      1.09,
	},
}

func init() {
	BenefitsSpareChangeExecution01.Model.Timestamp, _ = internal.ParseTimestamp("2024-01-04T12:26:52.110+0000")
	BenefitsSpareChangeExecution01.CSVEntry.Timestamp = internal.DateTime{Time: BenefitsSpareChangeExecution01.Model.Timestamp}

	TestCasesSupported["BenefitsSpareChangeExecution01"] = BenefitsSpareChangeExecution01
}
