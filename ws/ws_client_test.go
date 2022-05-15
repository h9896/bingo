package ws

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/h9896/bingo/events"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func echo(data []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, _, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, data)
			if err != nil {
				break
			}
		}
	}
}

func TestWsClient(t *testing.T) {
	data := []byte(`{
		"e": "aggTrade",
		"s": "BTCUSD_PERP",
		"p": "30000.2",
		"q": "234"
	}`)

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo(data)))
	defer s.Close()

	cfg := WsConfig{
		Name: s.URL,
	}

	client := newTestClient(t)
	cleanup, err := StartSubscribe(client, cfg)
	if err != nil {
		log.Println(err)
		cleanup()
	}
	time.Sleep(2 * time.Second)

	cleanup()
}

type testClient struct {
	aggregate *events.AggregateMsg
	t         *testing.T
	wg        sync.WaitGroup
}

func newTestClient(t *testing.T) WsClient {
	return &testClient{
		aggregate: &events.AggregateMsg{},
		t:         t,
	}
}

func (c *testClient) GetEndpoint(cfg WsConfig) string {
	// Convert http://127.0.0.1 to ws://127.0.0.
	return "ws" + strings.TrimPrefix(cfg.Name, "http")
}

func (c *testClient) GetServices(cfg WsConfig) []string {
	services := []string{}
	for _, symbol := range cfg.Symbols {
		services = append(services, fmt.Sprintf("%s@%s", symbol, cfg.Service))
	}
	return services
}

func (c *testClient) ErrHandler(err error) {
	log.Println(err)
}

func (c *testClient) MsgHandler(msg []byte) {
	c.aggregate.Unmarshal(msg)
	assert.EqualValues(c.t, "aggTrade", c.aggregate.EventType)
	assert.EqualValues(c.t, "30000.2", c.aggregate.Price)
	assert.EqualValues(c.t, "234", c.aggregate.Quantity)
}
