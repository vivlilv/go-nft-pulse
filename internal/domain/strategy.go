package domain

import (
	"math/big"
)

type Intent interface {
	isIntent() // unexported, just marks the type
}

type RefetchTopOfferIntent struct {
	NftID
}

func (RefetchTopOfferIntent) isIntent() {
}

type PlaceItemOfferIntent struct {
	NftID
	PriceWei *big.Int
	Duration int64
}

func (PlaceItemOfferIntent) isIntent() {}

type CancelOfferIntent struct {
	OrderHash string
}

func (CancelOfferIntent) isIntent() {}

type NOOP struct {
	Msg string
}

func (NOOP) isIntent() {}

type StrategyPanic struct {
	NftID
	State
}

func (StrategyPanic) isIntent() {

}
