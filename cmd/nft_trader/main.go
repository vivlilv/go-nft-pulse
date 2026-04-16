package main

import (
	"fmt"

	events_opensea "github.com/vivlilv/go_nft_trader/internal/events/opensea"
)

func main() {
	fmt.Println("Starting NFT Trader...")
	events_opensea.ListenEvents()
}
