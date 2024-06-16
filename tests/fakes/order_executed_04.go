package fakes

var OrderExecuted04 = TransactionTestCase{
	TimelineDetailsData: TimelineDetailsTestData{
		Raw: []byte(`{
		"id": "61f297f9-f9c3-46c4-a15c-cdd50d5544ad",
		"sections": [
		  {
			"action": {
			  "payload": "XF000AVAX016",
			  "type": "instrumentDetail"
			},
			"data": {
			  "icon": "logos/XF000AVAX016/v2",
			  "status": "executed",
			  "subtitleText": null,
			  "timestamp": "2024-03-12T15:21:56.707+0000"
			},
			"title": "Du hast 2.517,95 €  erhalten",
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
				  "text": "Verkauf",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Orderart"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "Avalanche",
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
			"action": null,
			"data": [
			  {
				"detail": {
				  "action": null,
				  "text": "1,71 %",
				  "trend": "positive",
				  "type": "text"
				},
				"style": "plain",
				"title": "Rendite"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "42,44 €",
				  "trend": "positive",
				  "type": "text"
				},
				"style": "plain",
				"title": "Gewinn"
			  }
			],
			"title": "Performance",
			"type": "horizontalTable"
		  },
		  {
			"action": null,
			"data": [
			  {
				"detail": {
				  "action": null,
				  "text": "65",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Anteile"
			  },
			  {
				"detail": {
				  "action": null,
				  "text": "38,75 €",
				  "trend": null,
				  "type": "text"
				},
				"style": "plain",
				"title": "Aktienkurs"
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
				  "text": "+ 2.517,95 €",
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
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				  "type": "browserModal"
				},
				"detail": "12.03.2024",
				"id": "7afcc1b3-42f4-4ecd-a40d-b17efd8b1478",
				"postboxType": "CRYPTO_SECURITIES_SETTLEMENT",
				"title": "Abrechnung"
			  },
			  {
				"action": {
				  "payload": "https://traderepublic-data-production.s3.eu-central-1.amazonaws.com/timeline/postbox",
				  "type": "browserModal"
				},
				"detail": "12.03.2024",
				"id": "f8c83c0f-5294-4854-9645-58c65bab8170",
				"postboxType": "COSTS_INFO_SELL_V2",
				"title": "Kosteninformation"
			  }
			],
			"title": "Dokumente",
			"type": "documents"
		  }
		]
	  }`),
	},
}
