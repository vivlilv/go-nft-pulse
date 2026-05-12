package domain

import (
	"math/big"
)

type Event interface {
}

type WSMessageEvent struct {
	EventType string `json:"event"`
}

type ItemListedEvent struct {
	EventType    string
	PriceWei     *big.Int
	Slug         string
	NftID        string
	TokenID      int
	MakerAddress string
	OrderHash    string
	UsdPrice     float64
	EndTime      int64
}

type CollectionOfferEvent struct {
	EventType    string
	PriceWei     *big.Int
	Slug         string
	Chain        string
	MakerAddress string
	OrderHash    string
	UsdPrice     float64
	EndTime      int64
}

type ItemCancelledEvent struct {
	EventType      string
	PriceWei       *big.Int
	Slug           string
	Chain          string
	NftID          string
	TokenID        int
	MakerAddress   string
	OrderHash      string
	UsdPrice       float64
	ExpirationDate string
	ListingDate    string
}

type ItemReceivedOfferEvent struct {
	EventType    string
	PriceWei     *big.Int
	Slug         string
	Chain        string
	NftID        string
	TokenID      int
	MakerAddress string
	OrderHash    string
	UsdPrice     float64
	EndTime      int64
}

type ItemSoldEvent struct {
	EventType    string
	Slug         string
	Chain        string
	NftID        string
	TokenID      int
	MakerAddress string
	TakerAddress string
	OrderHash    string
	PriceWei     *big.Int
	UsdPrice     float64
	ClosingDate  string
	EndTime      int64
	TxHash       string
}

type TraitOfferEvent struct {
	EventType         string
	Slug              string
	Chain             string
	ContractAddr      string
	PriceWei          *big.Int
	MakerAddress      string
	OrderHash         string
	UsdPrice          float64
	CreatedDate       string
	ExpirationDate    string
	EndTime           int64
	TraitCriteriaList []TraitCriterion
}
