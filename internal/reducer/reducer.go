package reducer

import (
	"fmt"
	"log"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

// Reduce takes the current state and an event, and returns the new state along with a list of affected NftIDs.
func Reduce(state domain.State, event domain.Event) (domain.State, []domain.NftID) {
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

		return state, []domain.NftID{key}

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
		return state, itemsToUpdate

	case domain.CollectionOfferEvent:
		offer := domain.OfferBase{
			PriceWei:       e.PriceWei,
			ExpirationTime: e.EndTime,
			Maker:          domain.Address(e.MakerAddress),
			OrderHash:      e.OrderHash,
		}
		fmt.Printf("%v", offer)

		itemsToUpdate := selectAllItemsForSlug(&state, e.Slug)
		for _, itemKey := range itemsToUpdate {
			applyOffer(&state, itemKey, offer)
		}
		return state, itemsToUpdate

	case domain.ItemSoldEvent: //FIXME - can be offer or listing sold;
		key := domain.NftID(e.NftID)
		clearOffer(&state, key, e.OrderHash)
		return state, []domain.NftID{key}

	case domain.ItemCancelledEvent: //FIXME - can be offer or listing cancel;
		key := domain.NftID(e.NftID)
		clearOffer(&state, key, e.OrderHash)
		return state, []domain.NftID{key}

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
	return state, []domain.NftID{}
}
