package events_opensea

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

func GetEventType(data []byte) string {
	// fmt.Printf("Raw event: %s\n", data) // << inspect here

	var msg WSMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return ""
	}

	// fmt.Printf("Received event type: %q\n", msg.EventType)
	return msg.EventType
}

func ProcessEvent(eventType string, data []byte) (domain.Event, error) {
	var e Event
	switch eventType {
	case "item_listed":
		e = &ItemListed{}
	case "item_sold":
		e = &ItemSold{}
	case "item_received_offer":
		e = &ItemReceivedOffer{}
	// case "item_received_bid":
	// 	e = &ItemReceivedBid{}
	case "item_cancelled":
		e = &ItemCancelled{}
	case "collection_offer":
		e = &CollectionOffer{}
	case "trait_offer":
		e = &TraitOffer{}
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	err := e.UnmarshalJSON(data)
	if err != nil {
		log.Printf("Unmarshaling event data: %v", err)
	}
	event, err := e.ToDomainEvent()
	if err != nil {
		log.Printf("Converting to domain event: %v", err)
	}
	b, _ := json.MarshalIndent(event, "", "  ")
	fmt.Println(string(b))
	return event, err
}

func SetupListener() *Listener {
	config := NewConfigMust()
	listener := NewEventsListener(config)
	return listener
}

func ReadMessages(listener *Listener, ch chan<- []byte, errChan chan<- error) {
	go func() {
		defer close(ch)
		for {
			_, data, err := listener.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			ch <- data
		}
	}()
}

func HandleEvents(
	listener *Listener,
	rawCh chan []byte,
	eventsCh chan<- domain.Event,
	errChanKeepAlive chan error,
	quit chan os.Signal,
) {
	go func() {
		for {
			select {
			case data := <-rawCh:
				fmt.Println(string(data))
				eventType := GetEventType(data)
				event, err := ProcessEvent(eventType, data)
				if err == nil {
					eventsCh <- event
				}
			case <-quit:
				if err := listener.Close(); err != nil { //closing the websocket.Conn will stop readMsg cycle
					fmt.Println("Error closing connection:", err)
				} else {
					fmt.Println("Closed successfully")
				}
				return
			case err := <-errChanKeepAlive:
				fmt.Println(err)
				return
			}
		}
	}()
}

func ListenEvents() (<-chan domain.Event, <-chan struct{}) {

	listener := SetupListener()

	err := listener.Subscribe("megalio-16")
	if err != nil {
		log.Fatal("during ListenEvents: ", err)
	}

	rawCh := make(chan []byte)          //raw ws events data channel
	eventsCh := make(chan domain.Event) //processed events channel
	quit := make(chan os.Signal, 1)
	errChanKeepAlive := make(chan error) //to catch errors from goroutine
	signal.Notify(quit, os.Interrupt)

	listener.KeepAlive(errChanKeepAlive)

	ReadMessages(listener, rawCh, errChanKeepAlive)
	HandleEvents(listener, rawCh, eventsCh, errChanKeepAlive, quit)

	return eventsCh, listener.Done()
}
