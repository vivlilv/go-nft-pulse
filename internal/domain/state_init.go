package domain

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"sync"
)

// func init() {
// 	godotenv.Load() // load .env before processing
// 	// if err := envconfig.Process("", MyWallet); err != nil {
// 	// 	log.Fatalf("failed to load wallet config: %v", err)
// 	// }
// }

type MyWalletConfig struct {
	Address Address `envconfig:"MY_ADDRESS" required:"true"`
}

var MyWallet = &MyWalletConfig{}

type NftID string
type Address string
type TraitCriterion struct {
	TraitType  string `json:"trait_type"` // category: "Background"
	TraitValue string `json:"trait_name"` // specific value: "Purple"
}

type StateManager struct {
	mu    sync.RWMutex
	state State
}

// GetState returns a snapshot copy of the current State
// safe to modify this copy
func (m *StateManager) GetState() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state.Snapshot()
}

func (m *StateManager) UpdateState(newState State) {
	m.mu.Lock()
	m.state = newState
	m.mu.Unlock()
}

type State struct {
	TradingMode string //item flip,collection wide, trait
	Mode        string //offer || listing

	Items map[NftID]ItemState
}

// NewStateManager initializes StateManager with a given State;StateManager is used to manage concurrent access to State across different goroutines
func NewStateManager(s State) *StateManager {
	return &StateManager{state: s}
}

// Snapshot creates a deep copy of the State
// cruicial because map is reference type;without deep copy, concurrent reads/writes can lead to race conditions
func (s State) Snapshot() State {
	deepCopyItems := map[NftID]ItemState{}
	for _, item := range s.Items {
		deepCopyItems[NftID(item.NftID)] = ItemState{
			Slug:             item.Slug,
			NftID:            item.NftID,
			IsPending:        item.IsPending,
			TokenID:          item.TokenID,
			MaxBidAllowedWei: item.MaxBidAllowedWei,
			OfferStepWei:     item.OfferStepWei,
			TopOffer:         item.TopOffer,
			MyOffer:          item.MyOffer,
			Traits:           item.Traits,
			ImgURL:           item.ImgURL,
		}
	}

	return State{
		TradingMode: s.TradingMode,
		Mode:        s.Mode,
		Items:       deepCopyItems,
	}
}

type ItemState struct {
	Slug             string //used for collection wide filter
	NftID            string
	IsPending        bool //whether currently doing some operation(place offer etc)
	TokenID          int
	MaxBidAllowedWei *big.Int
	OfferStepWei     *big.Int
	TopOffer         OfferBase
	MyOffer          OfferBase        //YAGNI - for now only per item bid;* if nil
	Traits           []TraitCriterion //each nft has traits;to query trait offers
	ImgURL           string
}

type OfferBase struct {
	PriceWei       *big.Int
	OrderHash      string
	ExpirationTime int64
	Maker          Address
}

type ItemOffer struct {
	NftID string
	OfferBase
}

type TraitOffer struct {
	Traits []TraitCriterion
	OfferBase
}

type CollectionOffer struct {
	OfferBase
}

func NewState(tradingMode string, mode string, path string) *State {

	items := InitializeItems(path)

	return &State{
		TradingMode: tradingMode,
		Mode:        mode,
		Items:       items,
	}
}

func InitializeItems(path string) map[NftID]ItemState {
	items := make(map[NftID]ItemState)
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
		items[NftID(item.NftID)] = item
	}

	return items
}
