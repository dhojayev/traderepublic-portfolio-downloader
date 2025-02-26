package timeline_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

var CardSuccessfulTransaction01 = fakes.TimelineFakeData{
	RawResponse: []byte(`{
	"action": {
		"payload": "6221f5fb-b8fa-4ad6-8c99-c3fb3c31da10",
		"type": "timelineDetail"
	},
	"amount": {
		"currency": "EUR",
		"fractionDigits": 2,
		"value": -5.95
	},
	"badge": null,
	"eventType": "card_successful_transaction",
	"icon": "logos/merchant-45180dc7-8917-45c9-b926-6ae7b3befe28/v2",
	"id": "6221f5fb-b8fa-4ad6-8c99-c3fb3c31da10",
	"status": "EXECUTED",
	"subAmount": null,
	"subtitle": null,
	"timestamp": "2024-05-27T13:51:55.167+0000",
	"title": "Aldi"
}`),
	Unmarshalled: transactions.ResponseItem{
		Action: transactions.ResponseItemAction{
			Payload: "6221f5fb-b8fa-4ad6-8c99-c3fb3c31da10",
			Type:    "timelineDetail",
		},
		Amount: transactions.ResponseItemAmount{
			Currency:       "EUR",
			FractionDigits: 2,
			Value:          -5.95,
		},
		EventType: "card_successful_transaction",
		Icon:      "logos/merchant-45180dc7-8917-45c9-b926-6ae7b3befe28/v2",
		ID:        "6221f5fb-b8fa-4ad6-8c99-c3fb3c31da10",
		Status:    "EXECUTED",
		Timestamp: "2024-05-27T13:51:55.167+0000",
		Title:     "Aldi",
	},
}

func init() {
	TestCases["CardSuccessfulTransaction01"] = CardSuccessfulTransaction01
}
