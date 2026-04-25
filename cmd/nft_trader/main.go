package main

import (
	"encoding/json"
	"fmt"

	domain_state "github.com/vivlilv/go_nft_trader/internal/domain"
)

func main() {
	fmt.Println("Starting NFT Trader...")

	state := domain_state.NewState("item_flip", "offer", "items.json")
	b, _ := json.MarshalIndent(state, "", "  ")
	fmt.Println(string(b))
	// events_opensea.ListenEvents()
}
