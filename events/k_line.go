package events

import (
	"encoding/json"
)

type KLineMsg struct {
	EventType string      `json:"e,omitempty"`
	EventTime int64       `json:"E,omitempty"`
	Symbol    string      `json:"s,omitempty"`
	Data      Candlestick `json:"k,omitempty"`
}
type Candlestick struct {
	StartTime               int64  `json:"t,omitempty"`
	CloseTime               int64  `json:"T,omitempty"`
	Symbol                  string `json:"s,omitempty"`
	Interval                string `json:"i,omitempty"`
	FirstTradeID            int64  `json:"f,omitempty"`
	LastTradeId             int64  `json:"L,omitempty"`
	OpenPrice               string `json:"o,omitempty"`
	ClosePrice              string `json:"c,omitempty"`
	HighPrice               string `json:"h,omitempty"`
	LowPrice                string `json:"l,omitempty"`
	Volume                  string `json:"v,omitempty"`
	NumberOfTrades          int64  `json:"n,omitempty"`
	CloseFlag               bool   `json:"x,omitempty"`
	BaseAssetVolume         string `json:"q,omitempty"`
	TakerBuyVolume          string `json:"V,omitempty"`
	TakerBuyBaseAssetVolume string `json:"Q,omitempty"`
	Ignore                  string `json:"B,omitempty"`
}

func (msg *KLineMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}

type KLineContractMsg struct {
	EventType    string      `json:"e,omitempty"`
	EventTime    int64       `json:"E,omitempty"`
	Pair         string      `json:"ps,omitempty"`
	ContractType string      `json:"ct,omitempty"`
	Data         Candlestick `json:"k,omitempty"`
}

func (msg *KLineContractMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}
