package bus

import "github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"

const (
	TopicTimelineTransactions = string(traderepublic.WebsocketSubRequestSchemaJsonTypeTimelineTransactions)
	TopicTimelineDetailsV2    = string(traderepublic.WebsocketSubRequestSchemaJsonTypeTimelineDetailV2)
)
