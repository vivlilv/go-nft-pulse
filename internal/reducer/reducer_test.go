package reducer

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/vivlilv/go_nft_trader/internal/domain"
)

func TestReduce(t *testing.T) {
	//given
	myAddress := domain.MyWallet.Address
	state := domain.NewState("item_flip", "offer", "items.json") //FIXME - state should be set by setState only

	NFT_ID := "ethereum/0x8fe1a377b83921fe1429adb1b8fbfecd45de9cd8/4705"
	event0 := domain.CollectionOfferEvent{ //base top offer - somebody else
		EventType:    "collection_offer",
		PriceWei:     big.NewInt(3300000000000000),
		Slug:         "megalio-16",
		Chain:        "ethereum",
		MakerAddress: "0xRANDOM",
		OrderHash:    "0xRANDOMHASH",
		UsdPrice:     6.55265,
		EndTime:      1778585440,
	}
	event1 := domain.CollectionOfferEvent{
		EventType:    "collection_offer",
		PriceWei:     big.NewInt(5500000000000000),
		Slug:         "megalio-16",
		Chain:        "ethereum",
		MakerAddress: "0x54c6d73baf9dd8612978b694b23a309263b54cfa",
		OrderHash:    "0x166f10778be4ab4b81e5fa3f2479ab2c8e1206b4e21786483d9dcbad613659c4",
		UsdPrice:     12.55265,
		EndTime:      1778585440,
	}
	event2 := domain.CollectionOfferEvent{
		EventType:    "collection_offer",
		PriceWei:     big.NewInt(5500000000000000),
		Slug:         "megalio-16",
		Chain:        "ethereum",
		MakerAddress: "0xOTHERADDRESS",
		OrderHash:    "0xOTHERHASH",
		UsdPrice:     12.55265,
		EndTime:      1778585440,
	}
	event3 := domain.CollectionOfferEvent{
		EventType:    "collection_offer",
		PriceWei:     big.NewInt(2000000000000000),
		Slug:         "megalio-16",
		Chain:        "ethereum",
		MakerAddress: "0x54c6d73baf9dd8612978b694b23a309263b54cfa",
		OrderHash:    "0xMYHASH",
		UsdPrice:     12.55265,
		EndTime:      1778585440,
	}
	tests := []struct {
		name          string
		state         *domain.State
		event         domain.Event
		setupState    func() *domain.State
		checkExpected func(t *testing.T, state domain.State)
	}{
		{
			name:  "#1: my collection offer(highest) >>> update both top offer, my offer",
			state: state,
			event: event1,
			setupState: func() *domain.State {
				return domain.NewState("item_flip", "offer", "items.json")
			},
			checkExpected: func(t *testing.T, state domain.State) {
				item := state.Items[domain.NftID(NFT_ID)]

				if item.TopOffer.PriceWei != event1.PriceWei {
					t.Errorf("top offer priceWei diff %v %v", item.TopOffer.PriceWei, event1.PriceWei)
				}
				if item.TopOffer.OrderHash != event1.OrderHash {
					t.Errorf("top offer hash diff")
				}
				if item.TopOffer.ExpirationTime != event1.EndTime {
					t.Errorf("top offer expirationTime diff")
				}
				if item.TopOffer.Maker != domain.Address(event1.MakerAddress) {
					t.Errorf("top offer maker diff")
				}

				if item.MyOffer.PriceWei != event1.PriceWei {
					t.Errorf("my offer priceWei diff")
				}
				if item.MyOffer.OrderHash != event1.OrderHash {
					t.Errorf("my offer hash diff")
				}
				if item.MyOffer.ExpirationTime != event1.EndTime {
					t.Errorf("my offer expirationTime diff")
				}
				if item.MyOffer.Maker != myAddress {
					t.Errorf("my offer maker diff")
				}
			}},
		{
			name:  "#2: collection offer(highest) >>> update top offer only",
			state: state,
			event: event2,
			setupState: func() *domain.State {
				return domain.NewState("item_flip", "offer", "items.json")
			},
			checkExpected: func(t *testing.T, state domain.State) {
				item := state.Items[domain.NftID(NFT_ID)]

				if item.TopOffer.PriceWei != event2.PriceWei {
					t.Errorf("top offer priceWei diff")
				}
				if item.TopOffer.OrderHash != event2.OrderHash {
					t.Errorf("top offer hash diff")
				}
				if item.TopOffer.ExpirationTime != event2.EndTime {
					t.Errorf("top offer expirationTime diff")
				}
				if item.TopOffer.Maker != domain.Address(event2.MakerAddress) {
					t.Errorf("top offer maker diff")
				}

				if item.MyOffer.OrderHash != "" {
					t.Errorf("My offer hash should be null")
				}
			},
		},
		{
			name:  "#3: my collection offer(lower) >>> update my offer only",
			state: state,
			event: event3,
			setupState: func() *domain.State {
				state := domain.NewState("item_flip", "offer", "items.json")
				result := Reduce(*state, event0)
				return &result

			},
			checkExpected: func(t *testing.T, state domain.State) {
				item := state.Items[domain.NftID(NFT_ID)]

				if item.MyOffer.PriceWei != event3.PriceWei {
					t.Errorf("my offer priceWei diff")
				}
				if item.MyOffer.OrderHash != event3.OrderHash {
					t.Errorf("my offer hash diff")
				}
				if item.MyOffer.ExpirationTime != event3.EndTime {
					t.Errorf("my offer expirationTime diff")
				}
				if item.MyOffer.Maker != domain.Address(event3.MakerAddress) {
					t.Errorf("my offer maker diff")
				}

				if item.TopOffer.OrderHash == event3.OrderHash {
					t.Errorf("Top offer hash shouldn't be same as my offer")
				}
			},
		},
	}

	//when
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fresh state for each case
			state := tt.setupState()
			log1, _ := json.MarshalIndent(state, "", "  ")
			t.Logf("State setup:\n%s", log1)
			// when
			result := Reduce(*state, tt.event)

			// then
			tt.checkExpected(t, result)
			log2, _ := json.MarshalIndent(result, "", "  ")
			t.Logf("State after reduce:\n%s", log2)
		})
	}

}
