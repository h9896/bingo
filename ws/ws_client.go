package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

func StartSubscribe[T WsClient](client T, cfg WsConfig) (cleanup func(), err error) {
	conn, _, err := websocket.DefaultDialer.Dial(client.GetEndpoint(cfg), nil)
	cleanup = func() {}
	if err != nil {
		return
	}

	conn.SetPingHandler(
		func(appData string) error {
			return conn.WriteMessage(websocket.PongMessage, []byte{})
		},
	)

	quit := make(chan struct{})
	cleanup = clean(quit, nil, conn)
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-time.After(25 * time.Millisecond):
				_, message, err := conn.ReadMessage()
				if err != nil {
					client.ErrHandler(err)
					return
				}
				client.MsgHandler(message)
			}
		}
	}()
	err = sendSubReq(client.GetServices(cfg), conn)
	if err != nil {
		return
	}
	cleanup = clean(quit, client.GetServices(cfg), conn)
	return
}

func clean(quit chan struct{}, streams []string, conn *websocket.Conn) func() {
	if streams == nil {
		return func() {
			defer close(quit)
			defer conn.Close()
		}
	} else {
		return func() {
			defer close(quit)
			sendUnSubReq(streams, conn)
		}
	}
}

func sendSubReq(streams []string, conn *websocket.Conn) error {
	req := &SubReq{
		Method: Subscribe,
		Params: streams,
		Id:     time.Now().Unix(),
	}
	buff, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, buff)
	if err != nil {
		return err
	}
	return nil
}

func sendUnSubReq(streams []string, conn *websocket.Conn) error {
	defer conn.Close()
	req := &SubReq{
		Method: Unsubscribe,
		Params: streams,
		Id:     time.Now().Unix(),
	}
	buff, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, buff)
	if err != nil {
		return err
	}
	return nil
}
