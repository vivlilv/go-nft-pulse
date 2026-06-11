package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vivlilv/go_nft_trader/internal/domain"
	events_opensea "github.com/vivlilv/go_nft_trader/internal/events/opensea"
	"github.com/vivlilv/go_nft_trader/internal/reducer"
	"github.com/vivlilv/go_nft_trader/internal/strategy"
	"github.com/vivlilv/go_nft_trader/web"
)

func main() {
	fmt.Println("Starting NFT Trader...")
	state := domain.NewState("item_flip", "offer", "items.json")

	stateManager := domain.NewStateManager(*state)
	server := web.NewServer(stateManager)
	strategyCoordinator := strategy.NewStrategyCoordinator()
	reducerMsgCh := strategyCoordinator.ReducerChan()

	eventsCh, done := events_opensea.ListenEvents("pudgypenguins")

	go server.Run()
	go logFinalState(stateManager)
	go runReducer(stateManager, eventsCh, reducerMsgCh)
	go strategyCoordinator.Run()

	<-done
}

func runReducer(stateManager *domain.StateManager, eventsCh <-chan domain.Event, reducerMsgCh chan<- strategy.ReducerMessage) {
	for event := range eventsCh {
		state := stateManager.GetState()
		updatedState, affectedItems := reducer.Reduce(state, event)
		stateManager.UpdateState(updatedState)

		reducerMsgCh <- strategy.ReducerMessage{ //FIXME - ? do i need to pass whole state(or strategy can ask it from stateManager)
			State:         updatedState,
			AffectedItems: affectedItems,
		}
	}
}

// when the program receives an interrupt signal, print the final state before exiting
func logFinalState(sm *domain.StateManager) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	state := sm.GetState()
	stateJSON, _ := json.MarshalIndent(state, "", "  ")
	log.Printf("final state:\n%s", stateJSON)
	os.Exit(0)
}
