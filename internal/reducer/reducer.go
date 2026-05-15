package reducer

import (
	"log"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

func Reduce(state domain.State, event domain.Event) domain.State {
	switch e := event.(type) {

	case domain.ItemReceivedOfferEvent:
		key := domain.NftID(e.NftID)
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          domain.Address(e.MakerAddress),
			OrderHash:      e.OrderHash,
		}

		applyOffer(&state, key, offer)

		return state

	case domain.TraitOfferEvent:
		traitsFilter := e.TraitCriteriaList
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          domain.Address(e.MakerAddress),
			OrderHash:      e.OrderHash,
		}

		itemsToUpdate := filterItemsByTraits(&state, e.Slug, traitsFilter)
		for _, itemKey := range itemsToUpdate {
			applyOffer(&state, itemKey, offer)
		}
		return state

	case domain.CollectionOfferEvent:
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          domain.Address(e.MakerAddress),
			OrderHash:      e.OrderHash,
		}

		itemsToUpdate := selectAllItemsForSlug(&state, e.Slug) //TODO - Add Caching
		for _, itemKey := range itemsToUpdate {
			applyOffer(&state, itemKey, offer)
		}
		return state

	case domain.ItemSoldEvent: //FIXME - can be offer or listing sold;
		key := domain.NftID(e.NftID)
		clearOffer(&state, key, e.OrderHash)
		return state

	case domain.ItemCancelledEvent: //FIXME - can be offer or listing cancel;
		key := domain.NftID(e.NftID)
		clearOffer(&state, key, e.OrderHash)
		return state

	case domain.ExpiredOfferEvent:
		switch e.OfferKind {
		case "collection_offer":
			nftIDs := selectAllItemsForSlug(&state, e.Slug)
			for _, itemKey := range nftIDs {
				clearOffer(&state, itemKey, e.OrderHash)
			}

		case "trait_offer":
			nftIDs := filterItemsByTraits(&state, e.Slug, e.Traits)
			for _, itemKey := range nftIDs {
				clearOffer(&state, itemKey, e.OrderHash)
			}
		case "item_offer":
			key := domain.NftID(e.NftID)
			clearOffer(&state, key, e.OrderHash)
		default:
			//ignore if some unhandled event type
		}

	default:
		log.Printf("Reducer: unhandled event type: %T", e)
	}
	return state
}
