package timeline_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/transactions"
)

var InterestPayoutCreated02 = Fake{
	RawResponse: []byte(`{
      "action": {
        "payload": "79bea4ff-5552-45a1-85b5-8977e29f3d04",
        "type": "timelineDetail"
      },
      "amount": {
        "currency": "EUR",
        "fractionDigits": 2,
        "value": 20.28
      },
      "badge": null,
      "eventType": "INTEREST_PAYOUT_CREATED",
      "icon": "logos/timeline_interest_new/v2",
      "id": "79bea4ff-5552-45a1-85b5-8977e29f3d04",
      "status": "EXECUTED",
      "subAmount": null,
      "subtitle": "4,00% p.a.",
      "timestamp": "2024-03-03T20:45:47.367+0000",
      "title": "Zinsen"
  }`),
	Unmarshalled: transactions.ResponseItem{
		Action: transactions.ResponseItemAction{
			Payload: "79bea4ff-5552-45a1-85b5-8977e29f3d04",
			Type:    "timelineDetail",
		},
		Amount: transactions.ResponseItemAmount{
			Currency:       "EUR",
			FractionDigits: 2,
			Value:          20.28,
		},
		EventType: "INTEREST_PAYOUT_CREATED",
		Icon:      "logos/timeline_interest_new/v2",
		ID:        "79bea4ff-5552-45a1-85b5-8977e29f3d04",
		Status:    "EXECUTED",
		Subtitle:  "4,00% p.a.",
		Timestamp: "2024-03-03T20:45:47.367+0000",
		Title:     "Zinsen",
	},
}

func init() {
	TestCases["InterestPayoutCreated02"] = InterestPayoutCreated02
}
