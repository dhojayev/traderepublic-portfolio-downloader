package fakes

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
)

var CardSuccessfulTransaction02 = TransactionTestCase{
	TimelineTransactionsData: TimelineTransactionsTestData{
		Raw: []byte(`{
		"items": 
			[
				{
					"action": {
						"payload": "9aa0f0a1-1b68-412d-8f90-71ef77a10f45",
						"type": "timelineDetail"
					},
					"amount": {
						"currency": "EUR",
						"fractionDigits": 2,
						"value": -157.93
					},
					"badge": null,
					"eventType": "card_successful_transaction",
					"icon": "logos/merchant-fallback-entertainment/v2",
					"id": "9aa0f0a1-1b68-412d-8f90-71ef77a10f45",
					"status": "EXECUTED",
					"subAmount": {
						"currency": "CZK",
						"fractionDigits": 2,
						"value": -3900
					},
					"subtitle": null,
					"timestamp": "2024-05-23T11:37:27.519+0000",
					"title": "Home Depot"
				}
			]
}`),
		Unmarshalled: transactions.ResponseItem{
			Action: transactions.ResponseItemAction{
				Payload: "9aa0f0a1-1b68-412d-8f90-71ef77a10f45",
				Type:    "timelineDetail",
			},
			Amount: transactions.ResponseItemAmount{
				Currency:       "EUR",
				FractionDigits: 2,
				Value:          -157.93,
			},
			EventType: "card_successful_transaction",
			Icon:      "logos/merchant-fallback-entertainment/v2",
			ID:        "9aa0f0a1-1b68-412d-8f90-71ef77a10f45",
			SubAmount: transactions.ResponseItemAmount{
				Currency:       "CZK",
				FractionDigits: 2,
				Value:          -3900,
			},
			Status:    "EXECUTED",
			Timestamp: "2024-05-23T11:37:27.519+0000",
			Title:     "Home Depot",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte("{}"),
	},
}

func init() {
	RegisterUnsupported(CardSuccessfulTransaction02)
}
