package activitylog

import "github.com/dhojayev/traderepublic-portfolio-downloader/internal"

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

func (a ResponseItemAction) HasDetails() bool {
	return a.Type == internal.ResponseActionTypeTimelineDetail && a.Payload != ""
}

type ResponseItemAction struct {
	Payload string `json:"payload"`
	Type    string `json:"type"`
}
