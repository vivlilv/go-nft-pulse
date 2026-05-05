package reducer

import (
	"log"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

func Reduce(state domain.State, event domain.Event) domain.State {
	switch e := event.(type) {

	case *domain.ItemReceivedOfferEvent:
		key := domain.NftID(e.NftID)
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          e.MakerAddress,
			OrderHash:      e.OrderHash,
		}

		applyOffer(&state, key, offer)

		return state

	case *domain.TraitOfferEvent:
		traitsFilter := e.TraitCriteriaList
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          e.MakerAddress,
			OrderHash:      e.OrderHash,
		}

		itemsToUpdate := filterItemsByTraits(traitsFilter, &state) //TODO - Add Caching
		for _, itemKey := range itemsToUpdate {
			applyOffer(&state, itemKey, offer)
		}
		return state

	case *domain.CollectionOfferEvent:
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          e.MakerAddress,
			OrderHash:      e.OrderHash,
		}

		itemsToUpdate := selectAllItems(&state) //TODO - Add Caching
		for _, itemKey := range itemsToUpdate {
			applyOffer(&state, itemKey, offer)
		}
		return state

	case *domain.ItemSoldEvent:
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          e.MakerAddress,
			OrderHash:      e.OrderHash,
		}

		key := domain.NftID(e.NftID)
		clearOffer(&state, key, offer)
		return state

	case *domain.ItemCancelledEvent:
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          e.MakerAddress,
			OrderHash:      e.OrderHash,
		}

		key := domain.NftID(e.NftID)
		clearOffer(&state, key, offer)
		return state

	default:
		log.Printf("Reducer: unhandled event type: %T", e)
	}
	return state
}
