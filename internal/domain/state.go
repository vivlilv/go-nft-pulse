package domain_state

import "math/big"

type nftID string
type address string
type TraitCriterion struct {
	Name string `json:"trait_name"`
	Type string `json:"trait_type"`
}

type State struct {
	TradingMode string //item flip,collection wide, trait
	Mode        string //offer || listing

	Items map[nftID]ItemState
}

type ItemState struct {
	NftID            string
	IsPending        bool //whether currently doing some operation(place offer etc)
	TokenID          int
	MaxBidAllowedWei *big.Int
	OfferStepWei     *big.Int
	TopOffer         Offer
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
