package events_opensea

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type EventsListener interface {
	KeepAlive(errChan chan<- error)
	Subscribe(collectionSlug string) error
	Close() error
}

type Listener struct {
	ctx    context.Context
	cancel func()
	conn   *websocket.Conn //websocket connection
}

func NewEventsListener(config Config) *Listener {
	apiEndpoint := fmt.Sprintf("wss://stream.openseabeta.com/socket/websocket?token=%v", config.OPENSEA_API_KEY)
	conn, _, err := websocket.DefaultDialer.Dial(apiEndpoint, nil)
	if err != nil {
		panic(err)
	}

	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)
	return &Listener{
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
	}
}

func (e *Listener) ReadMessage() (messageType int, p []byte, err error) {
	msgType, data, err := e.conn.ReadMessage()
	return msgType, data, err
}

func (e *Listener) WriteMessage(messageType int, p []byte) error {
	err := e.conn.WriteMessage(messageType, p)
	return err
}

func (e *Listener) KeepAlive(c chan<- error) {
	//FIXME - it has to run in some goroutine to not block others.
	//Also when connection is closed send signal to break the cycle
	//Now handled with errChan, change it to context.Context
	go func() {
		for {
			heartbeatMessage := `{"topic": "phoenix", "event": "heartbeat", "payload": {}, "ref": 0}`
			err := e.WriteMessage(websocket.TextMessage, []byte(heartbeatMessage))
			if err != nil {
				c <- fmt.Errorf("heartbeat keepalive:%v", err)
				return
			}
			select {
			case <-e.ctx.Done():
				return //ctx.cancel() was called - stop keepAlive
			case <-time.After(30 * time.Second):
				//if ctx wasn't closed - just do nothing and start new cycle
			}
		}
	}()
}

func (e *Listener) Subscribe(collectionSlug string) error {
	subscribeMessage := fmt.Sprintf(`{"topic": "collection:%v", "event": "phx_join", "payload": {}, "ref": 0}`, collectionSlug)
	err := e.WriteMessage(websocket.TextMessage, []byte(subscribeMessage))
	if err != nil {
		return fmt.Errorf("subscribing to collection:%v", err)
	}

	return nil
}

func (e *Listener) Close() error {
	err := e.conn.Close() //TODO - implement graceful shutdown with context
	e.cancel()
	return err
}

func (e *Listener) Done() <-chan struct{} {
	return e.ctx.Done()
}
