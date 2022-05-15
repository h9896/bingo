package events

import (
	"encoding/json"
)

type AggregateMsg struct {
	EventType        string `json:"e,omitempty"`
	EventTime        int64  `json:"E,omitempty"`
	AggregateTradeID int    `json:"a,omitempty"`
	Symbol           string `json:"s,omitempty"`
	Price            string `json:"p,omitempty"`
	Quantity         string `json:"q,omitempty"`
	FirstTradeID     int64  `json:"f,omitempty"`
	LastTradeID      int64  `json:"l,omitempty"`
	TradeTime        int64  `json:"T,omitempty"`
	MarketMaker      bool   `json:"m,omitempty"`
}

func (msg *AggregateMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}
