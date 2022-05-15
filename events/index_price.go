package events

import (
	"encoding/json"
)

type IndexPriceMsg struct {
	EventType string `json:"e,omitempty"`
	EventTime int64  `json:"E,omitempty"`
	Pair      string `json:"i,omitempty"`
	Price     string `json:"p,omitempty"`
}

func (msg *IndexPriceMsg) Unmarshal(in []byte) error {
	return json.Unmarshal(in, msg)
}
