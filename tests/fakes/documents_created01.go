package fakes

import "github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/activitylog"

var DocumentsCreated01 = ActivityLogTestCase{
	ActivityLogData: ActivityLogTestData{
		Raw: []byte(`{
		"items": [
			{
				"action": {
					"payload": "55600b18-d064-4346-a143-2df1013de583",
					"type": "timelineDetail"
				},
				"eventType": "DOCUMENTS_CREATED",
				"icon": "logos/timeline_document/v2",
				"id": "55600b18-d064-4346-a143-2df1013de583",
				"subtitle": "Hinzugefügt",
				"timestamp": "2023-09-12T06:56:19.127+0000",
				"title": "Rechtliche Dokumente"
			}
		]
	}`),
		Unmarshalled: activitylog.ResponseItem{
			Action: activitylog.ResponseItemAction{
				Payload: "55600b18-d064-4346-a143-2df1013de583",
				Type:    "timelineDetail",
			},
			EventType: "DOCUMENTS_CREATED",
			Icon:      "logos/timeline_document/v2",
			ID:        "55600b18-d064-4346-a143-2df1013de583",
			Subtitle:  "Hinzugefügt",
			Timestamp: "2023-09-12T06:56:19.127+0000",
			Title:     "Rechtliche Dokumente",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
  "id": "55600b18-d064-4346-a143-2df1013de583",
  "sections": [
    {
      "action": null,
      "data": {
        "icon": "logos/timeline_document/v2",
        "status": "executed",
        "subtitleText": null,
        "timestamp": "2023-09-12T06:56:19.127+0000"
      },
      "title": "Du hast neue Dokumente erhalten",
      "type": "header"
    },
    {
      "action": null,
      "data": [
        {
          "action": {
            "payload": "https://assets.traderepublic.com/documents/DE/CONTRACT_CUSTOMER_AGREEMENT_20230914171325.pdf",
            "type": "browserModal"
          },
          "detail": "Aktuell von Trade Republic verwendete Kundenvereinbarung",
          "id": "55600b18-d064-4346-a143-2df1013de583",
          "postboxType": "DOCUMENTS_CREATED",
          "title": "Kundenvereinbarung"
        },
        {
          "action": {
            "payload": "https://assets.traderepublic.com/documents/DE/FURTHER_INFORMATION_DATA_PROTECTION_20230831080812.pdf",
            "type": "browserModal"
          },
          "detail": "Wie wir Deine vertraglichen Daten schützen",
          "id": "55600b18-d064-4346-a143-2df1013de583",
          "postboxType": "DOCUMENTS_CREATED",
          "title": "Datenschutzinformationen"
        },
        {
          "action": {
            "payload": "https://assets.traderepublic.com/documents/DE/FURTHER_INFORMATION_PRICE_LIST_20230426144143.pdf",
            "type": "browserModal"
          },
          "detail": "Unsere Preise, fair und transparent",
          "id": "55600b18-d064-4346-a143-2df1013de583",
          "postboxType": "DOCUMENTS_CREATED",
          "title": "Preis- \u0026 Leistungsverzeichnis"
        },
        {
          "action": {
            "payload": "https://assets.traderepublic.com/documents/DE/STOCK_PERK_TERMS_AND_CONDITIONS_20221205145701.pdf",
            "type": "browserModal"
          },
          "detail": "NCP-AGB",
          "id": "55600b18-d064-4346-a143-2df1013de583",
          "postboxType": "DOCUMENTS_CREATED",
          "title": "NCP-AGB"
        },
        {
          "action": {
            "payload": "https://assets.traderepublic.com/documents/DE/einlagensicherung_de_20220314125857.pdf",
            "type": "browserModal"
          },
          "detail": "Informationen zur Einlagensicherung",
          "id": "55600b18-d064-4346-a143-2df1013de583",
          "postboxType": "DOCUMENTS_CREATED",
          "title": "Informationen zur Einlagensicherung"
        },
        {
          "action": {
            "payload": "https://assets.traderepublic.com/documents/DE/einlagensicherung_de_20210215155905.pdf",
            "type": "browserModal"
          },
          "detail": "Informationen zur Einlagensicherung",
          "id": "55600b18-d064-4346-a143-2df1013de583",
          "postboxType": "DOCUMENTS_CREATED",
          "title": "Informationen zur Einlagensicherung"
        }
      ],
      "title": "Dokumente",
      "type": "documents"
    }
  ]
}`),
	},
}

func init() {
	RegisterActivityLogSupported("DocumentsCreated01", DocumentsCreated01)
}
