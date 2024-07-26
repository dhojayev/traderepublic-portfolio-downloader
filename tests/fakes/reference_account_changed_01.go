package fakes

import "github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/timeline/activitylog"

var ReferenceAccountChanged01 = ActivityLogTestCase{
	ActivityLogData: ActivityLogTestData{
		Raw: []byte(`{
		"items": [
			{
			"action": {
				"payload": { 
					"verificationTransfer": false
					},
				"type": "redirectCashAccountScreen"
			},
			"eventType": "REFERENCE_ACCOUNT_CHANGED",
			"icon": "logos/timeline_bank/v2",
			"id": "5b23e2a2-ff54-49a5-8f08-4ecb58385dff",
			"subtitle": "Geändert",
			"timestamp": "2023-11-16T17:28:56.013+0000",
			"title": "Auszahlungskonto"
			}
		]
		}`),
		Unmarshalled: activitylog.ResponseItem{
			Action: activitylog.ResponseItemAction{
				Payload: map[string]any{
					"verificationTransfer": false,
				},
				Type: "redirectCashAccountScreen",
			},
			EventType: "REFERENCE_ACCOUNT_CHANGED",
			Icon:      "logos/timeline_bank/v2",
			ID:        "5b23e2a2-ff54-49a5-8f08-4ecb58385dff",
			Subtitle:  "Geändert",
			Timestamp: "2023-11-16T17:28:56.013+0000",
			Title:     "Auszahlungskonto",
		},
	},
	TimelineDetailsData: TimelineDetailsTestData{},
}

func init() {
	RegisterActivityLogSupported("ReferenceAccountChanged01", ReferenceAccountChanged01)
}
