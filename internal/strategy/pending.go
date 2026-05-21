package strategy

import (
	"log"
	"sync"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

type Coordinator interface {
	Add(id domain.NftID)
	Remove(id domain.NftID)
	Has(id domain.NftID) bool
	ReducerChan() chan<- ReducerMessage
	ExecutorChan() chan<- ExecutorResult
}

type StrategyCoordinator struct {
	mu             sync.RWMutex
	pendingItems   map[domain.NftID]struct{}
	reducerOutput  chan ReducerMessage
	executorResult chan ExecutorResult
}

type ReducerMessage struct {
	State         domain.State
	AffectedItems []domain.NftID
}

type ExecutorResult struct {
	ID      domain.NftID
	Success bool
}

func (s *StrategyCoordinator) ReducerChan() chan<- ReducerMessage {
	return s.reducerOutput
}

func (s *StrategyCoordinator) ExecutorChan() chan<- ExecutorResult {
	return s.executorResult
}

func (s *StrategyCoordinator) Add(id domain.NftID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pendingItems[id] = struct{}{}
}

func (s *StrategyCoordinator) Remove(id domain.NftID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pendingItems, id)
}

func (s *StrategyCoordinator) Has(id domain.NftID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.pendingItems[id]
	return ok
}

func (s *StrategyCoordinator) Run() {
	go func() {
		for {
			select {
			case msg := <-s.reducerOutput: //reducer call to run strategy
				//check for pendingItems if not already processing.
				//call strategy to produce intent
				//send intent to executor
				for _, id := range msg.AffectedItems {
					isPending := s.Has(id)
					if !isPending {
						s.Add(id)
						intent := Strategy(msg.State, id)
						//pass intent to executor via channel
						log.Printf("INTENT: %T", intent)
					}

				}

			case res := <-s.executorResult: //executor result to update pendingItems
				//recieve executorResult and update Coordinator(change pendingItems, potentially use optimistic updates)
				//TODO - implement this(placeholder for now)
				log.Printf("%v", res)
			}
		}
	}()
}

func NewStrategyCoordinator() *StrategyCoordinator {
	return &StrategyCoordinator{
		pendingItems:   make(map[domain.NftID]struct{}),
		reducerOutput:  make(chan ReducerMessage),
		executorResult: make(chan ExecutorResult),
	}
}
