package ws

const (
	Ws_coin_futures = "dstream.binance.com/ws"

	Ws_format   = "%s://%s"
	Subscribe   = "SUBSCRIBE"
	Unsubscribe = "UNSUBSCRIBE"
)

// WsClient - an interface of websocket client
type WsClient interface {
	ErrHandler(err error)
	MsgHandler(msg []byte)
	GetEndpoint(cfg WsConfig) string
	GetServices(cfg WsConfig) []string
}

//
type WsConfig struct {
	UseSSL  bool
	Name    string
	Symbols []string
	Service string
}

type SubReq struct {
	Method string   `json:"method,omitempty"`
	Params []string `json:"params,omitempty"`
	Id     int64    `json:"id,omitempty"`
}
