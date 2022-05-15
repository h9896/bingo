package events

import (
	"encoding/json"
)

type PartialBookDepthMsg struct {
	EventType           string     `json:"e,omitempty"`
	EventTime           int64      `json:"E,omitempty"`
	TransactionTime     int64      `json:"T,omitempty"`
	Symbol              string     `json:"s,omitempty"`
	Pair                string     `json:"ps,omitempty"`
	FirstUpdateID       int64      `json:"U,omitempty"`
	FinalUpdateID       int64      `json:"u,omitempty"`
	FinalUpdateStreamID int64      `json:"pu,omitempty"`
	Bids                [][]string `json:"b,omitempty"`
	Asks                [][]string `json:"a,omitempty"`
}

func (msg *PartialBookDepthMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}
