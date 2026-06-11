package strategy

import (
	"fmt"
	"math/big"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

type situation int

const (
	noData       = iota + 1 // mN, tN
	noMyOffer               // mN, tY
	myOfferNoTop            // mY, tN
	winning                 // mY, tY, m==t
	losing                  // mY, tY, m<t
	logPanic                // mY, tY, m>t
)

// classify the situation of the item based on my offer and the top offer
func classify(item domain.ItemState) situation {
	myOffer := item.MyOffer
	topOffer := item.TopOffer

	if myOffer.PriceWei == nil && topOffer.PriceWei == nil {
		return noData
	}
	if myOffer.PriceWei == nil {
		return noMyOffer
	}
	if topOffer.PriceWei == nil {
		return myOfferNoTop
	}
	if myOffer.PriceWei.Cmp(topOffer.PriceWei) == 0 {
		return winning
	}
	if myOffer.PriceWei.Cmp(topOffer.PriceWei) == -1 {
		return losing
	}
	if myOffer.PriceWei.Cmp(topOffer.PriceWei) == 1 {
		return logPanic //myOffer should never be higher than topOffer
	}
	return logPanic //unforseen case
}

func Strategy(snapshot domain.State, nftID domain.NftID) domain.Intent {
	// TODO - add functionality to clear ALL myOffers for current item.
	item := snapshot.Items[nftID]

	// budget guard — always first
	if item.MyOffer.PriceWei != nil && item.MyOffer.PriceWei.Cmp(item.MaxBidAllowedWei) > 0 {
		return domain.CancelOfferIntent{OrderHash: item.MyOffer.OrderHash}
	}

	switch classify(item) {
	case noData, myOfferNoTop:
		return domain.RefetchTopOfferIntent{NftID: nftID}
	case noMyOffer, losing: // same logic — place if room, else noop
		newPrice := new(big.Int).Add(item.TopOffer.PriceWei, item.OfferStepWei)
		if newPrice.Cmp(item.MaxBidAllowedWei) > 0 {
			return domain.NOOP{Msg: "Potential offer would be over maxBid"}
		}
		return domain.PlaceItemOfferIntent{
			NftID:    nftID,
			PriceWei: newPrice,
			Duration: 86400,
		}
	case winning:
		return domain.NOOP{Msg: fmt.Sprintf("My offer is on top for %v", nftID)}
	default:
		return domain.StrategyPanic{
			NftID: nftID,
			State: snapshot,
		}
	}
}
