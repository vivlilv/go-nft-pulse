package domain_state

import (
	"encoding/json"
	"log"
	"os"
)

func NewState(tradingMode string, mode string, path string) *State {

	items := InitializeItems(path)

	return &State{
		TradingMode: tradingMode,
		Mode:        mode,
		Items:       items,
	}
}

func InitializeItems(path string) map[nftID]ItemState {
	items := make(map[nftID]ItemState)
	itemsData := []ItemState{}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = json.Unmarshal(data, &itemsData)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}
	for _, item := range itemsData {
		items[nftID(item.NftID)] = item
	}

	return items
}
