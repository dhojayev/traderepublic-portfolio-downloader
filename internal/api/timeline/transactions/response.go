package transactions

const (
	transactionActionTypeTimelineDetail = "timelineDetail"
)

type Response struct {
	Type  string              `json:"type"`
	Items TransactionResponse `json:"data"`
}

type TransactionResponse struct {
	Action    TransactionActionResponse `json:"action,omitempty"`
	Amount    TransactionAmountResponse `json:"amount"`
	Badge     string                    `json:"badge,omitempty"`
	EventType string                    `json:"eventType"`
	Icon      string                    `json:"icon"`
	ID        string                    `json:"id"`
	Status    string                    `json:"status"`
	SubAmount string                    `json:"subAmount,omitempty"`
	Subtitle  string                    `json:"subtitle,omitempty"`
	Timestamp string                    `json:"timestamp"`
	Title     string                    `json:"title"`
}

type TransactionActionResponse struct {
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

type TransactionAmountResponse struct {
	Currency       string  `json:"currency"`
	FractionDigits uint8   `json:"fractionDigits"`
	Value          float32 `json:"value"`
}

func (a TransactionActionResponse) HasTimelineDetail() bool {
	return a.Type == transactionActionTypeTimelineDetail
}
