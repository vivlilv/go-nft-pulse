package strategy

import (
	"math/big"
	"testing"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

func TestStrategy(t *testing.T) {
	//given
	state := domain.NewState("item_flip", "offer", "testdata/items.json") //FIXME - state should be set by setState only
	NFT_ID := domain.NftID("ethereum/0xbd3531da5cf5857e7cfaa92426877b022e612cf8/920")

	tests := []struct {
		name          string
		state         domain.State
		setupState    func() *domain.State
		checkExpected func(t *testing.T, intent domain.Intent)
	}{
		{
			name:  "No top offer in state - refetchTopOffer intent",
			state: state.Snapshot(),
			setupState: func() *domain.State {
				return state
			},
			checkExpected: func(t *testing.T, intent domain.Intent) {
				_, ok := intent.(domain.RefetchTopOfferIntent)
				if !ok {
					t.Errorf("expected RefetchTopOfferIntent, got %T", intent)
				}
			},
		},

		{
			name:  "Top offer exists, my offer empty -> placeOffer intent(case room for bidding exists)",
			state: state.Snapshot(),
			setupState: func() *domain.State {
				item := state.Items[NFT_ID]

				item.TopOffer = domain.OfferBase{
					PriceWei:       big.NewInt(98000000000000000),
					OrderHash:      "0xHASH",
					ExpirationTime: 0,
					Maker:          "0xMAKER",
				}

				state.Items[NFT_ID] = item
				return state
			},
			checkExpected: func(t *testing.T, intent domain.Intent) {
				_, ok := intent.(domain.PlaceItemOfferIntent)
				if !ok {
					t.Errorf("expected PlaceItemOfferIntent, got %T", intent)
				}
			},
		},
		{
			name:  "Top offer exists, my offer empty -> NOOP intent(case NO room for bidding)",
			state: state.Snapshot(),
			setupState: func() *domain.State {
				item := state.Items[NFT_ID]

				item.TopOffer = domain.OfferBase{
					PriceWei:       big.NewInt(1980000000000000000),
					OrderHash:      "0xHASH",
					ExpirationTime: 0,
					Maker:          "0xMAKER",
				}

				state.Items[NFT_ID] = item
				return state
			},
			checkExpected: func(t *testing.T, intent domain.Intent) {
				_, ok := intent.(domain.NOOP)
				if !ok {
					t.Errorf("expected NOOP, got %T", intent)
				}
			},
		},
		{
			name:  "Top offer exists, my offer exists -> NOOP intent(WINNING)",
			state: state.Snapshot(),
			setupState: func() *domain.State {
				item := state.Items[NFT_ID]

				myTopOffer := domain.OfferBase{
					PriceWei:       big.NewInt(198000000000000000),
					OrderHash:      "0xMYHASH",
					ExpirationTime: 0,
					Maker:          "0xMYMAKER",
				}
				item.MyOffer = myTopOffer
				item.TopOffer = myTopOffer

				state.Items[NFT_ID] = item
				return state
			},
			checkExpected: func(t *testing.T, intent domain.Intent) {
				_, ok := intent.(domain.NOOP)
				if !ok {
					t.Errorf("expected NOOP, got %T", intent)
				}
			},
		},
		{
			name:  "Top offer exists, my offer exists -> PlaceItemOfferIntent intent(LOSING)",
			state: state.Snapshot(),
			setupState: func() *domain.State {
				item := state.Items[NFT_ID]

				item.TopOffer = domain.OfferBase{
					PriceWei:       big.NewInt(198000000000000000),
					OrderHash:      "0xHASH",
					ExpirationTime: 0,
					Maker:          "0xMAKER",
				}
				item.MyOffer = domain.OfferBase{
					PriceWei:       big.NewInt(197000000000000000),
					OrderHash:      "0xMYHASH",
					ExpirationTime: 0,
					Maker:          "0xMYMAKER",
				}

				state.Items[NFT_ID] = item
				return state
			},
			checkExpected: func(t *testing.T, intent domain.Intent) {
				_, ok := intent.(domain.PlaceItemOfferIntent)
				if !ok {
					t.Errorf("expected NOOP, got %T", intent)
				}
			},
		},
	}
	//when
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			state := tt.setupState()
			// when
			intent := Strategy(*state, NFT_ID)

			// then
			tt.checkExpected(t, intent)
		})
	}
}
