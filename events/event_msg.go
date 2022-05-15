package events

type EventHandler interface {
	Unmarshal(in []byte) error
}
