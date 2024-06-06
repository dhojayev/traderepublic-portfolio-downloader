package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
)

var PaymentOutbound01 = TestCase{
	TimelineDetailsData: TimelineDetailsData{
		Raw: `{
  "id": "a2597441-45f4-4ae2-a881-ab4a65aa0f0e",
  "sections": [
    {
      "action": null,
      "data": {
        "icon": "logos/timeline_minus_circle/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2024-01-11T08:55:22.185+0000"
      },
      "title": "Du hast 1,00 € gesendet",
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
            "text": "Mr. Bean",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "An"
        },
        {
          "detail": {
            "action": null,
            "text": "DE14 1234 5678 9012 3456 78",
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
}`,
	},
	EventType: transactions.EventTypePaymentOutbound,
	Transaction: transaction.Model{
		UUID:      "a2597441-45f4-4ae2-a881-ab4a65aa0f0e",
		Type:      transaction.TypeWithdrawal,
		Status:    "executed",
		Total:     1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	CSVEntry: filesystem.CSVEntry{
		ID:        "a2597441-45f4-4ae2-a881-ab4a65aa0f0e",
		Status:    "executed",
		Type:      "Withdrawal",
		AssetType: "Other",
		Debit:     1,
	},
}

func init() {
	PaymentOutbound01.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2024-01-11T08:55:22.185+0000")
	PaymentOutbound01.CSVEntry.Timestamp = internal.DateTime{Time: PaymentOutbound01.Transaction.Timestamp}

	RegisterSupported(PaymentOutbound01)
}
