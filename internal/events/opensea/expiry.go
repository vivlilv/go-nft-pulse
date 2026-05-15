package events_opensea

import (
	"sync"
	"time"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

type ExpiryScheduler struct {
	eventsCh chan<- domain.Event
	timers   map[string]*time.Timer // keyed by OrderHash
	mu       sync.Mutex
}

func NewExpiryScheduler(ch chan<- domain.Event) *ExpiryScheduler {
	return &ExpiryScheduler{
		eventsCh: ch,
		timers:   make(map[string]*time.Timer),
	}
}

func (s *ExpiryScheduler) Register(endTime int64, orderHash string, offerKind string, slug string, nftID string, traits []domain.TraitCriterion) {
	//TODO - right now only implementation for offers; add listings later
	s.mu.Lock()
	defer s.mu.Unlock()

	duration := time.Until(time.Unix(endTime, 0))
	newTimer := time.AfterFunc(duration, func() {
		//fire the ExpiredEvent onto eventsCh

		newEvent := domain.ExpiredOfferEvent{
			OrderHash: orderHash,
			OfferKind: offerKind,
			Slug:      slug,
			NftID:     nftID,
			Traits:    traits,
		}

		s.eventsCh <- newEvent
	})

	//in case order hash existed before on the map - clear old timer & delete, set new
	if timer, ok := s.timers[orderHash]; ok {
		timer.Stop()
		delete(s.timers, orderHash)
	}
	s.timers[orderHash] = newTimer

}

func (s *ExpiryScheduler) Cancel(orderHash string) {
	//delete timer from map (event canceled, fulfilled, expired)
	s.mu.Lock()
	defer s.mu.Unlock()

	if timer, ok := s.timers[orderHash]; ok {
		timer.Stop()
		delete(s.timers, orderHash)
	}
}

func ScheduleExpiration(s *ExpiryScheduler, event domain.Event) {
	var endTime int64
	var orderHash string
	var offerKind string
	var nftID = ""
	var traits []domain.TraitCriterion = nil
	var slug string

	switch e := event.(type) {
	case domain.CollectionOfferEvent:
		endTime = e.EndTime
		orderHash = e.OrderHash
		offerKind = e.EventType
		slug = e.Slug

	case domain.TraitOfferEvent:
		endTime = e.EndTime
		orderHash = e.OrderHash
		offerKind = e.EventType
		traits = e.TraitCriteriaList
		slug = e.Slug

	case domain.ItemReceivedOfferEvent:
		endTime = e.EndTime
		orderHash = e.OrderHash
		offerKind = e.EventType
		nftID = e.NftID
		slug = e.Slug

	default:
		//skip event if not on the list
		return
	}

	s.Register(
		endTime,
		orderHash,
		offerKind,
		slug,
		nftID,
		traits,
	)
}
