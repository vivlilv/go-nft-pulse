package reducer

import (
	"github.com/vivlilv/go_nft_trader/internal/domain"
)

func filterItemsByTraits(state *domain.State, slug string, traits []domain.TraitCriterion) []domain.NftID {
	itemsList := selectAllItemsForSlug(state, slug)
	result := make([]domain.NftID, 0, len(itemsList)) //prealloc to avoid realloc space

	for _, i := range itemsList {
		item := state.Items[i]
		if itemHasAllTraits(item, traits) {
			result = append(result, domain.NftID(i))
		}
	}
	return result
}

func itemHasAllTraits(item domain.ItemState, traits []domain.TraitCriterion) bool {
	for _, t := range traits {
		found := false
		for _, i := range item.Traits {
			if t == i {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func selectAllItemsForSlug(state *domain.State, slug string) []domain.NftID {
	itemsList := []domain.NftID{}

	for _, i := range state.Items {
		if i.Slug == slug {
			itemsList = append(itemsList, domain.NftID(i.NftID))
		}
	}
	return itemsList
}

func applyOffer(state *domain.State, key domain.NftID, offer domain.OfferBase) {
	item, exists := state.Items[key] // fresh read every time
	if !exists {
		return
	}

	if item.TopOffer.PriceWei == nil || item.TopOffer.PriceWei.Cmp(offer.PriceWei) == -1 {
		item.TopOffer = offer
	}
	if offer.Maker == domain.MyWallet.Address {
		item.MyOffer = offer
	}

	state.Items[key] = item // single write back
}

func clearOffer(state *domain.State, key domain.NftID, orderHash string) {
	item, exists := state.Items[key] // fresh read every time
	if !exists {
		return
	}

	if item.TopOffer.OrderHash == orderHash {
		item.TopOffer = domain.OfferBase{} //Schedule refetch
	}
	if item.MyOffer.OrderHash == orderHash {
		item.MyOffer = domain.OfferBase{}
	}

	state.Items[key] = item
}
