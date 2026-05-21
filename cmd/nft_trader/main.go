package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vivlilv/go_nft_trader/internal/domain"
	domain_state "github.com/vivlilv/go_nft_trader/internal/domain"
	events_opensea "github.com/vivlilv/go_nft_trader/internal/events/opensea"
	"github.com/vivlilv/go_nft_trader/internal/reducer"
	"github.com/vivlilv/go_nft_trader/internal/strategy"
)

func main() {
	// logFile, err := os.OpenFile("nft_trader.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// if err != nil {
	// 	log.Fatalf("failed to open log file: %v", err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logFile)

	fmt.Println("Starting NFT Trader...")
	state := domain_state.NewState("item_flip", "offer", "items.json")
	b, _ := json.MarshalIndent(state, "", "  ")
	fmt.Println(string(b))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		stateJSON, _ := json.MarshalIndent(state, "", "  ")
		log.Printf("final state:\n%s", stateJSON)
		os.Exit(0)
	}()

	stateManager := domain.NewStateManager(state)

	eventsCh, done := events_opensea.ListenEvents("pudgypenguins")

	strategyCoordinator := strategy.NewStrategyCoordinator()
	strategyCoordinator.Run()
	reducerMsgCh := strategyCoordinator.ReducerChan()

	runReducer(stateManager, eventsCh, reducerMsgCh)
	// analytics.Run(stateManager)
	<-done
}

func runReducer(stateManager *domain.StateManager, eventsCh <-chan domain.Event, reducerMsgCh chan<- strategy.ReducerMessage) {
	go func() {
		for event := range eventsCh {
			state := stateManager.GetState()
			updatedState, affectedItems := reducer.Reduce(*state, event)
			stateManager.UpdateStateItems(updatedState.Items)

			reducerMsgCh <- strategy.ReducerMessage{
				State:         updatedState,
				AffectedItems: affectedItems,
			}
		}
	}()
}
