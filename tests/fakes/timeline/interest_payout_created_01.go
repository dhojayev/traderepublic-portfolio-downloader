package timeline_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests/fakes"
)

var InterestPayoutCreated01 = fakes.TimelineFakeData{
	RawResponse: []byte(`{
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
}

func init() {
	TestCases["InterestPayoutCreated01"] = InterestPayoutCreated01
}
