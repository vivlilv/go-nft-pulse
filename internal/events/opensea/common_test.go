package events_opensea

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
)

type testExpected struct {
	EventType    string   `json:"event_type"`
	PriceWei     *big.Int `json:"priceWei"`
	Slug         string   `json:"slug"`
	NftID        string   `json:"nft_id"`
	TokenID      int      `json:"token_id"`
	MakerAddress string   `json:"maker_address"`
	OrderHash    string   `json:"orderHash"`
	PriceUsd     float64  `json:"priceUsd"`
	EndTime      int64    `json:"endTime"`
}

func TestItemListedUnmarshalJSON(t *testing.T) {
	data, err := os.ReadFile("testdata/item_listed.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var tests []struct {
		Name                 string          `json:"name"`
		Input                json.RawMessage `json:"input"`
		Expected             testExpected    `json:"expected"`
		ExpectedErr          bool            `json:"expect_err"`
		ExpectedErrToContain string          `json:"expect_err_msg,omitempty"`
	}
	if err := json.Unmarshal(data, &tests); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}
	for _, tt := range tests {
		b, _ := json.MarshalIndent(tt.Expected, "", "  ")
		fmt.Println(string(b)) //just log event for inspection

		t.Run(tt.Name, func(t *testing.T) {
			var item ItemListed
			err := json.Unmarshal([]byte(tt.Input), &item)
			if tt.ExpectedErr {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
				if tt.ExpectedErrToContain != "" && !strings.Contains(err.Error(), tt.ExpectedErrToContain) {
					t.Fatalf("Expected error containing %q, got %q", tt.ExpectedErrToContain, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if item.PriceWei.Cmp(tt.Expected.PriceWei) != 0 {
				t.Errorf("Expected PriceWei %v, got %v", tt.Expected.PriceWei, item.PriceWei)
			}
			if item.Slug != tt.Expected.Slug {
				t.Errorf("Expected Slug '%s', got '%s'", tt.Expected.Slug, item.Slug)
			}
			if item.NftID != tt.Expected.NftID {
				t.Errorf("Expected NftID '%s', got '%s'", tt.Expected.NftID, item.NftID)
			}
			if item.TokenID != tt.Expected.TokenID {
				t.Errorf("Expected TokenID %d, got %d", tt.Expected.TokenID, item.TokenID)
			}
			if item.MakerAddress != tt.Expected.MakerAddress {
				t.Errorf("Expected MakerAddress '%s', got '%s'", tt.Expected.MakerAddress, item.MakerAddress)
			}
			if item.OrderHash != tt.Expected.OrderHash {
				t.Errorf("Expected OrderHash '%s', got '%s'", tt.Expected.OrderHash, item.OrderHash)
			}
			if item.UsdPrice != tt.Expected.PriceUsd {
				t.Errorf("Expected UsdPrice '%v', got '%v'", tt.Expected.PriceUsd, item.UsdPrice)
			}
			if item.EndTime != tt.Expected.EndTime {
				t.Errorf("Expected EndTime %d, got %d", tt.Expected.EndTime, item.EndTime)
			}
		})
	}
}
