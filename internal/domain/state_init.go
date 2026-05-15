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
	Name string `json:"trait_name"`
	Type string `json:"trait_type"`
}

type StateManager struct {
	mu    sync.RWMutex
	state *State
}

func (m *StateManager) GetState() *State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

func (m *StateManager) UpdateStateItems(Items map[NftID]ItemState) {
	m.mu.Lock()
	m.state.Items = Items
	m.mu.Unlock()
}

func (m *StateManager) Snapshot() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state.Snapshot()
}

type State struct {
	TradingMode string //item flip,collection wide, trait
	Mode        string //offer || listing

	Items map[NftID]ItemState
}

func NewStateManager(s *State) *StateManager {
	return &StateManager{state: s}
}

func (s *State) Snapshot() State {
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
