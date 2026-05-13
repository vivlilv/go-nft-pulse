package reducer

import (
	"github.com/vivlilv/go_nft_trader/internal/domain"
)

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
