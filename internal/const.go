package internal

const (
	// RestAPIBaseURI Trade Republic's REST API base URI.
	RestAPIBaseURI = "https://api.traderepublic.com/api/v1"

	// WebsocketBaseHost Trade Republic's websocket base host.
	WebsocketBaseHost = "api.traderepublic.com"

	// HTTPUserAgent used for all HTTP communications.
	HTTPUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15"

	// CookieNamePrefix prefix used by Trade Republic for its auth cookies.
	CookieNamePrefix = "tr_"

	// SessionRefreshInterval represents in how many seconds session has to be refreshed to keep it alive.
	SessionRefreshInterval = 60

	// ResponseActionTypeTimelineDetail represents the value the app will look for in order to determine
	// if any details can be fetched.
	ResponseActionTypeTimelineDetail = "timelineDetail"

	// CSVFilename filename under which a CSV file with transaction entries has to be saved.
	CSVFilename = "./transactions.csv"

	// DocumentsBaseDir base directory under which downloaded documents have to be saved.
	DocumentsBaseDir = "./documents/transactions"
)
