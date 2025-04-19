package internal

import "time"

const (
	// RestAPIBaseURI Trade Republic's REST API base URI.
	RestAPIBaseURI = "https://api.traderepublic.com/api/v1"

	// WebsocketBaseHost Trade Republic's websocket base host.
	WebsocketBaseHost = "api.traderepublic.com"

	// HTTPUserAgent used for all HTTP communications.
	HTTPUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"

	// SessionRefreshInterval represents in how many seconds session has to be refreshed to keep it alive.
	SessionRefreshInterval = 120

	AuthTokenFilename = "./.auth"

	// ResponseActionTypeTimelineDetail represents the value the app will look for in order to determine
	// if any details can be fetched.
	ResponseActionTypeTimelineDetail = "timelineDetail"

	// ResponseTimeFormat represents the default date time format in the response.
	ResponseTimeFormat = "2006-01-02T15:04:05-0700"

	// ResponseTimeFormatAlt represents the alternative date time format in the response.
	ResponseTimeFormatAlt = time.RFC3339Nano

	// CSVFilename filename under which a CSV file with transaction entries has to be saved.
	CSVFilename = "./transactions.csv"

	// TransactionDocumentsBaseDir base directory under which downloaded transaction documents are saved.
	TransactionDocumentsBaseDir = "./documents/transactions"

	// ActivityLogDocumentsBaseDir base directory under which downloaded activity documents are saved.
	ActivityLogDocumentsBaseDir = "./documents/activity"
)
