# bingo

[![GoDoc](https://godoc.org/github.com/h9896/bingo?status.svg)](https://godoc.org/github.com/h9896/bingo)
[![Go Report Card](https://goreportcard.com/badge/github.com/h9896/bingo)](https://goreportcard.com/report/github.com/h9896/bingo)
[![codecov](https://codecov.io/gh/h9896/bingo/branch/main/graph/badge.svg)](https://codecov.io/gh/h9896/bingo)

---

A Golang API for [binance](https://www.binance.com) base on [Doc](https://binance-docs.github.io/apidocs/delivery/en/#change-log).

|      Name      |                        Description                         |       Status       |
| :------------: | :--------------------------------------------------------: | :----------------: |
| Coin-M Futures | Perpetual or Quarterly Contracts settled in Cryptocurrency | Partical Implement |
| USD-M Futures  |  Perpetual or Quarterly Contracts settled in USDT or BUSD  |        ToDo        |
|  Spot/Margin   |                                                            |        ToDo        |

## Example

### Websocket

Create a client and implement the WsClient

```go
type coin struct {
}

func NewCoinMFutures() ws.WsClient {
	return &coin{}
}

func (c *coin) GetEndpoint(cfg ws.WsConfig) string {
	if cfg.UseSSL {
		return fmt.Sprintf("%s/%s", fmt.Sprintf(ws.Ws_format, "wss", ws.Ws_coin_futures), cfg.Name)
	}
	return fmt.Sprintf("%s/%s", fmt.Sprintf(ws.Ws_format, "ws", ws.Ws_coin_futures), cfg.Name)
}

func (c *coin) GetServices(cfg ws.WsConfig) []string {
	services := []string{}
	for _, symbol := range cfg.Symbols {
		services = append(services, fmt.Sprintf("%s@%s", symbol, cfg.Service))
	}
	return services
}

func (c *coin) ErrHandler(err error) {
	// Process websocket error what you like
}

func (c *coin) MsgHandler(msg []byte) {
	// Handle message from what service you subscribe
}
```

Use StartSubscribe function to get message.

```go
	cfg := ws.WsConfig{
		UseSSL:  true,
		Name:    "CoinM",
		Symbols: []string{"btcusd_perp"},
		Service: "aggTrade",
	}

	client := NewCoinMFutures()

	cleanup, err := ws.StartSubscribe(client, cfg)
	if err != nil {
		log.Println(err)
		cleanup()
	}
	time.Sleep(time.Minute)
	// End subscription and clean up related connections
	cleanup()
```

### Http

Create a http client

```go
client := rpc.NewHttpClient(apikey, true, nil)
```

General Request

```go
req := client.GetHttpRequest(
	rpc.SetEndpoint("dapi.binance.com/dapi/v1/depth"),
	rpc.SetMethod("get"),
	rpc.SetParams(&rpc.HttpParameter{Key: "symbol", Val: "BTCUSD_PERP"},
		&rpc.HttpParameter{Key: "limit", Val: "5"}),
)

resp, err := client1.ExecuteHttpOperation(ctx, req)
```

Private Request

```go
req := client.GetHttpRequest(
	rpc.SetEndpoint("dapi.binance.com/dapi/v1/account"),
	rpc.SetMethod("get"),
	rpc.SetPrivate(),
	rpc.SetTimestamp(),
	rpc.SetSignature(secret),
	)

resp, err := client.ExecuteHttpOperation(ctx, req)
```

Private Request by POST

```go
body := []*rpc.HttpParameter{
	{Key: "symbol", Val: "BTCUSD_PERP"},
	{Key: "side", Val: "BUY"},
	{Key: "recvWindow", Val: "600000"},
	{Key: "type", Val: "LIMIT"},
	{Key: "timeInForce", Val: "GTC"},
	{Key: "quantity", Val: "0.0035"},
	{Key: "price", Val: "28000.1"}}

req := client.GetHttpRequest(rpc.SetEndpoint("dapi.binance.com/dapi/v1/order"),
	rpc.SetPrivate(),
	rpc.SetMethod("post"),
	rpc.SetParams(body...),
	rpc.SetTimestamp(),
	rpc.SetSignature(secret))

resp, err := client.ExecuteHttpOperation(ctx, req)
```
