package transactions

const (
	transactionActionTypeTimelineDetail = "timelineDetail"
)

type Response struct {
	Type  string       `json:"type"`
	Items ResponseItem `json:"data"`
}

type ResponseItem struct {
	Action    ResponseItemAction `json:"action,omitempty"`
	Amount    ResponseItemAmount `json:"amount"`
	Badge     any                `json:"badge,omitempty"`
	EventType string             `json:"eventType"`
	Icon      string             `json:"icon"`
	ID        string             `json:"id"`
	Status    string             `json:"status"`
	SubAmount ResponseItemAmount `json:"subAmount,omitempty"`
	Subtitle  string             `json:"subtitle,omitempty"`
	Timestamp string             `json:"timestamp"`
	Title     string             `json:"title"`
}

type ResponseItemAction struct {
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

type ResponseItemAmount struct {
	Currency       string  `json:"currency"`
	FractionDigits uint8   `json:"fractionDigits"`
	Value          float32 `json:"value"`
}

func (a ResponseItemAction) HasTimelineDetail() bool {
	return a.Type == transactionActionTypeTimelineDetail
}
