package events

import (
	"encoding/json"
)

type MarkPriceMsg struct {
	EventType      string `json:"e,omitempty"`
	EventTime      int64  `json:"E,omitempty"`
	Symbol         string `json:"s,omitempty"`
	Price          string `json:"p,omitempty"`
	EstimatedPrice string `json:"P,omitempty"`
	Rate           string `json:"r,omitempty"`
	FundingTime    int64  `json:"T,omitempty"`
}

func (msg *MarkPriceMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}
