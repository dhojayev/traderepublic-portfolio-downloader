package fakes

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/details"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/timeline/transactions"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/filesystem"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio/transaction"
	"github.com/dhojayev/traderepublic-portfolio-downloader/tests"
)

var OrderExecuted05 = tests.TestCase{
	ResponseJSON: `{
  "id": "eb6ee8c7-2cff-4dcc-ab70-3ca7f31f0371",
  "sections": [
    {
      "action": {
        "payload": "DE000A2TEDB8",
        "type": "instrumentDetail"
      },
      "data": {
        "icon": "logos/DE0007500001/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2023-09-12T06:35:52.879+0000"
      },
      "title": "Du hast 5.001,01 €  investiert",
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
            "text": "Kauf",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Orderart"
        },
        {
          "detail": {
            "action": null,
            "text": "Anleihe Feb. 2024",
            "trend": null,
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
      "action": {
        "payload": {
          "sections": [
            {
              "action": null,
              "title": "Transaktion",
              "type": "title"
            },
            {
              "action": null,
              "data": [
                {
                  "detail": {
                    "action": null,
                    "text": "4.948,04 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Nominale"
                },
                {
                  "detail": {
                    "action": null,
                    "text": "99,27 %",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Kurs"
                },
                {
                  "detail": {
                    "action": null,
                    "text": "4.911,92 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "highlighted",
                  "title": "Order"
                }
              ],
              "title": null,
              "type": "table"
            },
            {
              "action": null,
              "data": [
                {
                  "detail": {
                    "action": null,
                    "text": "4.911,92 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Order"
                },
                {
                  "detail": {
                    "action": null,
                    "text": "88,09 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Stückzinsen"
                },
                {
                  "detail": {
                    "action": null,
                    "text": "1,00 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "plain",
                  "title": "Gebühr"
                },
                {
                  "detail": {
                    "action": null,
                    "text": "5.001,01 €",
                    "trend": null,
                    "type": "text"
                  },
                  "style": "highlighted",
                  "title": "Gesamt"
                }
              ],
              "title": null,
              "type": "table"
            }
          ]
        },
        "type": "infoPage"
      },
      "data": [
        {
          "detail": {
            "action": null,
            "text": "1,00 €",
            "trend": null,
            "type": "text"
          },
          "style": "plain",
          "title": "Gebühr"
        },
        {
          "detail": {
            "action": null,
            "text": "5.001,01 €",
            "trend": null,
            "type": "text"
          },
          "style": "highlighted",
          "title": "Gesamt"
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
            "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/1209301239102390129309321.pdf",
            "type": "browserModal"
          },
          "detail": "12.09.2023",
          "id": "f5d837cd-831e-42f5-87bf-9939a68bd138",
          "postboxType": "SECURITIES_SETTLEMENT",
          "title": "Abrechnung 2"
        },
        {
          "action": {
            "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/1209301239102390129309322.pdf",
            "type": "browserModal"
          },
          "detail": "12.09.2023",
          "id": "0ae93fa9-9d7e-4043-90de-1326966ed141",
          "postboxType": "SECURITIES_SETTLEMENT",
          "title": "Abrechnung 1"
        },
        {
          "action": {
            "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/1209301239102390129309323.pdf",
            "type": "browserModal"
          },
          "detail": "12.09.2023",
          "id": "e9026951-cb24-4ff4-9bea-156f4b0d4693",
          "postboxType": "COSTS_INFO_BUY_V2",
          "title": "Kosteninformation 2"
        },
        {
          "action": {
            "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox/1209301239102390129309324.pdf",
            "type": "browserModal"
          },
          "detail": "12.09.2023",
          "id": "d2628369-b1c7-4ecf-814b-fe8c6c266ed1",
          "postboxType": "COSTS_INFO_BUY_V2",
          "title": "Kosteninformation 1"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    }
  ]
}`,
	Response:    tests.Response{},
	EventType:   transactions.EventTypeOrderExecuted,
	Transaction: transaction.Model{},
	CSVEntry:    filesystem.CSVEntry{},
}

func init() {
	OrderExecuted05.Transaction.Timestamp, _ = time.Parse(details.ResponseTimeFormat, "2023-09-12T06:35:52.879+0000")
	OrderExecuted05.CSVEntry.Timestamp = internal.DateTime{Time: OrderExecuted05.Transaction.Timestamp}

	tests.RegisterUnknown(OrderExecuted05)
}
