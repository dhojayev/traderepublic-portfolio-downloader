package bus

import "github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"

const (
	TopicTimelineTransactions = string(traderepublic.WsSubRequestJsonTypeTimelineTransactions)
	TopicTimelineDetailsV2    = string(traderepublic.WsSubRequestJsonTypeTimelineDetailV2)
	TopicInstrument           = string(traderepublic.WsSubRequestJsonTypeInstrument)

	EventNameTimelineTransactionsReceived = "timeline_transactions_received"
	EventNameTimelineDetailV2Received     = "timeline_detail_v2_received"
	EventNameInstrumentReceived           = "instrument_received"
)
