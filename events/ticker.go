package events

import "encoding/json"

type MinTickertMsg struct {
	EventType       string `json:"e,omitempty"`
	EventTime       int64  `json:"E,omitempty"`
	Symbol          string `json:"s,omitempty"`
	Pair            string `json:"ps,omitempty"`
	OpenPrice       string `json:"o,omitempty"`
	ClosePrice      string `json:"c,omitempty"`
	HighPrice       string `json:"h,omitempty"`
	LowPrice        string `json:"l,omitempty"`
	Volume          string `json:"v,omitempty"`
	BaseAssetVolume string `json:"q,omitempty"`
}

func (msg *MinTickertMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}

type IndividualSymbolTickertMsg struct {
	EventType            string `json:"e,omitempty"`
	EventTime            int64  `json:"E,omitempty"`
	Symbol               string `json:"s,omitempty"`
	Pair                 string `json:"ps,omitempty"`
	PriceChange          string `json:"p,omitempty"`
	PriceChangeRate      string `json:"P,omitempty"`
	WeightedAveragePrice string `json:"w,omitempty"`
	LastQuantity         string `json:"Q,omitempty"`
	OpenPrice            string `json:"o,omitempty"`
	ClosePrice           string `json:"c,omitempty"`
	HighPrice            string `json:"h,omitempty"`
	LowPrice             string `json:"l,omitempty"`
	Volume               string `json:"v,omitempty"`
	BaseAssetVolume      string `json:"q,omitempty"`
	StatisticsOpenTime   int64  `json:"O,omitempty"`
	StatisticsCloseTime  int64  `json:"C,omitempty"`
	FirstTradeID         int64  `json:"F,omitempty"`
	LastTradeID          int64  `json:"L,omitempty"`
	TotalNumberTrades    int64  `json:"n,omitempty"`
}

func (msg *IndividualSymbolTickertMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}

type BookTickertMsg struct {
	EventType         string `json:"e,omitempty"`
	OrderBookUpdateId int64  `json:"u,omitempty"`
	EventTime         int64  `json:"E,omitempty"`
	Symbol            string `json:"s,omitempty"`
	Pair              string `json:"ps,omitempty"`
	BestBidPrice      string `json:"b,omitempty"`
	BestBidQty        string `json:"B,omitempty"`
	BestAskPrice      string `json:"a,omitempty"`
	BestAskQty        string `json:"A,omitempty"`
	TransactionTime   int64  `json:"T,omitempty"`
}

func (msg *BookTickertMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}
