package domain

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
)

type NftID string
type address string
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
	MyOffer          *ItemOffer       //YAGNI - for now only per item bid;* if nil
	Traits           []TraitCriterion //each nft has traits;to query trait offers
}

type Offer interface {
	//TODO: add methods like GetPriceForItem(1), HandlelExpiration(), etc
}

type OfferBase struct {
	BasePrice      *big.Int
	QuantityMin    int
	QuantityMax    int
	ExpirationTime int64
	Maker          address
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
