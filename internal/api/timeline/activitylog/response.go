package activitylog

type Response struct {
	Type  string       `json:"type"`
	Items ResponseItem `json:"data"`
}

type ResponseItem struct {
	Action    ResponseItemAction `json:"action,omitempty"`
	EventType string             `json:"eventType"`
	Icon      string             `json:"icon"`
	ID        string             `json:"id"`
	Subtitle  string             `json:"subtitle,omitempty"`
	Timestamp string             `json:"timestamp"`
	Title     string             `json:"title"`
}

type ResponseItemAction struct {
	Payload string `json:"payload"`
	Type    string `json:"type"`
}