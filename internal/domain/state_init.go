package domain

import (
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type MyWalletConfig struct {
	Address Address `envconfig:"MY_ADDRESS" required:"true"`
}

var MyWallet = &MyWalletConfig{}

func init() {
	if err := envconfig.Process("", MyWallet); err != nil {
		log.Fatalf("failed to load wallet config: %v", err)
	}
}

type NftID string
type Address string
type TraitCriterion struct {
	Name string `json:"trait_name"`
	Type string `json:"trait_type"`
}

type State struct {
	TradingMode string //item flip,collection wide, trait
	Mode        string //offer || listing

	Items map[NftID]ItemState
}

type ItemState struct {
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
