package events

import "encoding/json"

type LiquidateOrderMsg struct {
	EventType string `json:"e,omitempty"`
	EventTime int64  `json:"E,omitempty"`
	Data      LOrder `json:"o,omitempty"`
}

type LOrder struct {
	Symbol                       string `json:"s,omitempty"`
	Pair                         string `json:"ps,omitempty"`
	Side                         string `json:"S,omitempty"`
	OrderType                    string `json:"o,omitempty"`
	TimeInForce                  string `json:"f,omitempty"`
	OriginalQuantity             string `json:"q,omitempty"`
	Price                        string `json:"p,omitempty"`
	AveragePrice                 string `json:"ap,omitempty"`
	BestBidPrice                 string `json:"b,omitempty"`
	OrderStatus                  string `json:"X,omitempty"`
	OrderLastFilledQuantity      string `json:"l,omitempty"`
	OrderLastAccumulatedQuantity string `json:"z,omitempty"`
	OrderTradeTime               string `json:"T,omitempty"`
}

func (msg *LiquidateOrderMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}
