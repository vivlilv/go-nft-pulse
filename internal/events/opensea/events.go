package events_opensea

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

// Opensea changed the format from json obj to array:
// parseMessage extracts event type and payload from Phoenix array format:
// [join_ref, ref, topic, event, payload]
func parseMessage(data []byte) (eventType string, payload []byte, err error) {
	var parts []json.RawMessage
	if err := json.Unmarshal(data, &parts); err != nil || len(parts) < 5 {
		return "", nil, fmt.Errorf("unexpected message format")
	}

	if err := json.Unmarshal(parts[3], &eventType); err != nil {
		return "", nil, fmt.Errorf("extracting event type: %w", err)
	}
	// wrap to preserve the double Payload.Payload nesting in unmarshalers
	//FIXME this additional wrapper is because all the methods rely on older architecture
	return eventType, fmt.Appendf(nil, `{"payload":%s}`, parts[4]), nil
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
	case "item_received_bid":
		e = &ItemReceivedOffer{} //TODO item received bid should have its own struct
	case "item_cancelled":
		e = &ItemCancelled{}
	case "collection_offer":
		e = &CollectionOffer{}
	case "trait_offer":
		e = &TraitOffer{}
	default:
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
	log.Println(string(b))
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
	expiryScheduler *ExpiryScheduler,
) {
	go func() {
		for {
			select {
			case data := <-rawCh:
				// fmt.Println(string(data))
				eventType, data, err := parseMessage(data)
				if err != nil {
					fmt.Printf("Parsing message: %v\n", err)
					continue //skip msg if parsing fails
				}
				event, err := ProcessEvent(eventType, data)
				if err == nil {
					ScheduleExpiration(expiryScheduler, event)
					eventsCh <- event
				}
			case <-quit:
				//closing the websocket.Conn will stop readMsg cycle
				if err := listener.Close(); err != nil {
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

func ListenEvents(collectionSlug string) (<-chan domain.Event, <-chan struct{}) {
	listener := SetupListener()

	err := listener.Subscribe(collectionSlug)
	if err != nil {
		log.Fatal("during ListenEvents: ", err)
	}

	rawCh := make(chan []byte)          //raw ws events data channel
	eventsCh := make(chan domain.Event) //processed events channel
	quit := make(chan os.Signal, 1)
	errChanKeepAlive := make(chan error) //to catch errors from goroutine
	signal.Notify(quit, os.Interrupt)
	scheduler := NewExpiryScheduler(eventsCh)
	listener.KeepAlive(errChanKeepAlive)

	ReadMessages(listener, rawCh, errChanKeepAlive)
	HandleEvents(listener, rawCh, eventsCh, errChanKeepAlive, quit, scheduler)

	return eventsCh, listener.Done()
}
