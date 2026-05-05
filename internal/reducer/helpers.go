package reducer

import "github.com/vivlilv/go_nft_trader/internal/domain"

func filterItemsByTraits(traits []domain.TraitCriterion, state *domain.State) []domain.NftID {
	itemsList := []domain.NftID{}

	for _, i := range state.Items {
		hasAllTraits := true
		for _, t := range traits {
			foundTrait := false
			for _, k := range i.Traits {
				if t == k {
					foundTrait = true
					break
				}
			}
			if !foundTrait {
				hasAllTraits = false
				break
			}
		}
		if hasAllTraits {
			itemsList = append(itemsList, domain.NftID(i.NftID))
		}
	}
	return itemsList
}

func selectAllItems(state *domain.State) []domain.NftID {
	itemsList := []domain.NftID{}

	for _, i := range state.Items {
		itemsList = append(itemsList, domain.NftID(i.NftID))
	}
	return itemsList
}

func applyOffer(state *domain.State, key domain.NftID, offer domain.OfferBase) {
	item := state.Items[key] // fresh read every time

	if item.TopOffer.PriceWei.Cmp(offer.PriceWei) == -1 {
		item.TopOffer = offer
	}
	if offer.Maker == domain.MyAddress {
		item.MyOffer = offer
	}

	state.Items[key] = item // single write back
}

func clearOffer(state *domain.State, key domain.NftID, offer domain.OfferBase) {
	item := state.Items[key]

	if item.TopOffer.OrderHash == offer.OrderHash {
		item.TopOffer = domain.OfferBase{} //Schedule refetch
	}
	if offer.Maker == domain.MyAddress {
		item.MyOffer = domain.OfferBase{}
	}

	state.Items[key] = item
}
