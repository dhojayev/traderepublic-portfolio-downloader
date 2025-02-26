package timeline_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

var InterestPayout01 = fakes.TimelineFakeData{
	RawResponse: []byte(`{
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
}

func init() {
	TestCases["InterestPayout01"] = InterestPayout01
}
