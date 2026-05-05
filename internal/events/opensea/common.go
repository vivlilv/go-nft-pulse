package events_opensea

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

type Event interface {
	UnmarshalJSON(data []byte) error
}

type WSMessage struct {
	EventType string `json:"event"`
}

type ItemListed struct {
	EventType    string
	PriceWei     *big.Int
	Slug         string
	NftID        string
	TokenID      int
	MakerAddress string
	OrderHash    string
	UsdPrice     float64
	EndTime      int64
}

func (e *ItemListed) UnmarshalJSON(data []byte) error {
	type msg struct {
		Payload struct {
			Payload struct {
				BasePrice  string `json:"base_price"`
				Collection struct {
					Slug string `json:"slug"`
				}
				Item struct {
					NftID string `json:"nft_id"`
				}
				Maker struct {
					Address string `json:"address"`
				} `json:"maker"`
				OrderHash    string `json:"order_hash"`
				PaymentToken struct {
					UsdPrice string `json:"usd_price"`
				} `json:"payment_token"`
				ProtocolData struct {
					Parameters struct {
						EndTime string `json:"endTime"`
					} `json:"parameters"`
				} `json:"protocol_data"`
			} `json:"payload"`
		} `json:"payload"`
	}
	var temp msg
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	priceUsd, err := strconv.ParseFloat(temp.Payload.Payload.PaymentToken.UsdPrice, 64)
	if err != nil {
		return err
	}

	tokenID, err := GetTokenID(temp.Payload.Payload.Item.NftID)
	if err != nil {
		return fmt.Errorf("unmarshalJSON item listed;failed to extract token ID: %w", err)
	}

	priceWei, err := StringToBigInt(temp.Payload.Payload.BasePrice)
	if err != nil {
		return fmt.Errorf("unmarshalJSON item listed;failed to convert base price to big.Int: %w", err)
	}
	e.EventType = "item_listed"
	e.PriceWei = priceWei
	e.Slug = temp.Payload.Payload.Collection.Slug
	e.NftID = temp.Payload.Payload.Item.NftID
	e.TokenID = tokenID
	e.MakerAddress = temp.Payload.Payload.Maker.Address
	e.OrderHash = temp.Payload.Payload.OrderHash
	e.UsdPrice = priceUsd

	endTime, err := strconv.Atoi(temp.Payload.Payload.ProtocolData.Parameters.EndTime)
	if err != nil {
		return err
	}
	e.EndTime = int64(endTime)

	return nil
}

type CollectionOffer struct {
	EventType    string
	PriceWei     *big.Int
	Slug         string
	Chain        string
	MakerAddress string
	OrderHash    string
	UsdPrice     string
	EndTime      int64
}

func (e *CollectionOffer) UnmarshalJSON(data []byte) error {
	type msg struct {
		Payload struct {
			Payload struct {
				BasePrice  string `json:"base_price"`
				Chain      string `json:"chain"`
				Collection struct {
					Slug string `json:"slug"`
				} `json:"collection"`
				Maker struct {
					Address string `json:"address"`
				} `json:"maker"`
				OrderHash    string `json:"order_hash"`
				PaymentToken struct {
					UsdPrice string `json:"usd_price"`
				} `json:"payment_token"`
				ProtocolData struct {
					Parameters struct {
						EndTime string `json:"endTime"`
					} `json:"parameters"`
				} `json:"protocol_data"`
			} `json:"payload"`
		} `json:"payload"`
	}

	var temp msg
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	priceWei, err := StringToBigInt(temp.Payload.Payload.BasePrice)
	if err != nil {
		return fmt.Errorf("unmarshalJSON collection offer;failed to convert base price to big.Int: %w", err)
	}
	e.EventType = "collection_offer"
	e.PriceWei = priceWei
	e.Slug = temp.Payload.Payload.Collection.Slug
	e.Chain = temp.Payload.Payload.Chain
	e.MakerAddress = temp.Payload.Payload.Maker.Address
	e.OrderHash = temp.Payload.Payload.OrderHash
	e.UsdPrice = temp.Payload.Payload.PaymentToken.UsdPrice

	endTime, err := strconv.Atoi(temp.Payload.Payload.ProtocolData.Parameters.EndTime)
	if err != nil {
		return err
	}
	e.EndTime = int64(endTime)

	return nil
}

type ItemCancelled struct {
	EventType      string
	PriceWei       *big.Int
	Slug           string
	Chain          string
	NftID          string
	TokenID        int
	MakerAddress   string
	OrderHash      string
	UsdPrice       string
	ExpirationDate string
	ListingDate    string
}

func (e *ItemCancelled) UnmarshalJSON(data []byte) error {
	type msg struct {
		Payload struct {
			Payload struct {
				BasePrice  string `json:"base_price"`
				Chain      string `json:"chain"`
				Collection struct {
					Slug string `json:"slug"`
				} `json:"collection"`
				ExpirationDate string `json:"expiration_date"`
				ListingDate    string `json:"listing_date"`
				Item           struct {
					NftID string `json:"nft_id"`
				} `json:"item"`
				Maker struct {
					Address string `json:"address"`
				} `json:"maker"`
				OrderHash    string `json:"order_hash"`
				PaymentToken struct {
					UsdPrice string `json:"usd_price"`
				} `json:"payment_token"`
			} `json:"payload"`
		} `json:"payload"`
	}

	var temp msg
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	priceWei, err := StringToBigInt(temp.Payload.Payload.BasePrice)
	if err != nil {
		return fmt.Errorf("unmarshalJSON item cancelled;failed to convert base price to big.Int: %w", err)
	}

	tokenID, err := GetTokenID(temp.Payload.Payload.Item.NftID)
	if err != nil {
		return fmt.Errorf("unmarshalJSON item cancelled;failed to extract token ID: %w", err)
	}
	e.EventType = "item_cancelled"
	e.PriceWei = priceWei
	e.Slug = temp.Payload.Payload.Collection.Slug
	e.Chain = temp.Payload.Payload.Chain
	e.NftID = temp.Payload.Payload.Item.NftID
	e.TokenID = tokenID
	e.MakerAddress = temp.Payload.Payload.Maker.Address
	e.OrderHash = temp.Payload.Payload.OrderHash
	e.UsdPrice = temp.Payload.Payload.PaymentToken.UsdPrice
	e.ExpirationDate = temp.Payload.Payload.ExpirationDate
	e.ListingDate = temp.Payload.Payload.ListingDate

	return nil
}

type ItemReceivedOffer struct {
	EventType      string
	PriceWei       *big.Int
	Slug           string
	Chain          string
	NftID          string
	TokenID        int
	MakerAddress   string
	OrderHash      string
	UsdPrice       string
	EndTime        int64
	CreatedDate    string
	ExpirationDate string
}

func (e *ItemReceivedOffer) UnmarshalJSON(data []byte) error {
	type msg struct {
		Payload struct {
			Payload struct {
				BasePrice  string `json:"base_price"`
				Chain      string `json:"chain"`
				Collection struct {
					Slug string `json:"slug"`
				} `json:"collection"`
				CreatedDate    string `json:"created_date"`
				ExpirationDate string `json:"expiration_date"`
				Item           struct {
					NftID string `json:"nft_id"`
				} `json:"item"`
				Maker struct {
					Address string `json:"address"`
				} `json:"maker"`
				OrderHash    string `json:"order_hash"`
				PaymentToken struct {
					UsdPrice string `json:"usd_price"`
				} `json:"payment_token"`
				ProtocolData struct {
					Parameters struct {
						EndTime string `json:"endTime"`
					} `json:"parameters"`
				} `json:"protocol_data"`
			} `json:"payload"`
		} `json:"payload"`
	}

	var temp msg
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	priceWei, err := StringToBigInt(temp.Payload.Payload.BasePrice)
	if err != nil {
		return fmt.Errorf("unmarshalJSON item received offer;failed to convert base price to big.Int: %w", err)
	}

	tokenID, err := GetTokenID(temp.Payload.Payload.Item.NftID)
	if err != nil {
		return fmt.Errorf("unmarshalJSON item received offer;failed to extract token ID: %w", err)
	}
	e.EventType = "item_received_offer"
	e.PriceWei = priceWei
	e.Slug = temp.Payload.Payload.Collection.Slug
	e.Chain = temp.Payload.Payload.Chain
	e.NftID = temp.Payload.Payload.Item.NftID
	e.TokenID = tokenID
	e.MakerAddress = temp.Payload.Payload.Maker.Address
	e.OrderHash = temp.Payload.Payload.OrderHash
	e.UsdPrice = temp.Payload.Payload.PaymentToken.UsdPrice

	endTime, err := strconv.Atoi(temp.Payload.Payload.ProtocolData.Parameters.EndTime)
	if err != nil {
		return err
	}
	e.EndTime = int64(endTime) // or assign to a field if you add EndTime int64 to the struct

	return nil
}

type ItemSold struct {
	EventType    string
	Slug         string
	Chain        string
	NftID        string
	TokenID      int
	MakerAddress string
	TakerAddress string
	OrderHash    string
	PriceWei     *big.Int
	UsdPrice     string
	ClosingDate  string
	EndTime      int64
	TxHash       string
}

func (e *ItemSold) UnmarshalJSON(data []byte) error {
	type msg struct {
		Payload struct {
			Payload struct {
				Chain          string `json:"chain"`
				EventTimestamp string `json:"event_timestamp"`
				Collection     struct {
					Slug string `json:"slug"`
				} `json:"collection"`
				ClosingDate string `json:"closing_date"`
				Item        struct {
					NftID string `json:"nft_id"`
				} `json:"item"`
				Maker struct {
					Address string `json:"address"`
				} `json:"maker"`
				Taker struct {
					Address string `json:"address"`
				} `json:"taker"`
				OrderHash    string `json:"order_hash"`
				PaymentToken struct {
					UsdPrice string `json:"usd_price"`
				} `json:"payment_token"`
				SalePrice   string `json:"sale_price"`
				Transaction struct {
					Hash string `json:"hash"`
				} `json:"transaction"`
			} `json:"payload"`
		} `json:"payload"`
	}

	var temp msg
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	tokenID, err := GetTokenID(temp.Payload.Payload.Item.NftID)
	if err != nil {
		return fmt.Errorf("unmarshalJSON msg;failed to extract token ID: %w", err)
	}

	priceWei, err := StringToBigInt(temp.Payload.Payload.SalePrice)
	if err != nil {
		return fmt.Errorf("unmarshalJSON msg;failed to convert sale price to big.Int: %w", err)
	}

	endTime, err := strconv.Atoi(temp.Payload.Payload.EventTimestamp)
	if err != nil {
		return err
	}

	e.EndTime = int64(endTime)
	e.Slug = temp.Payload.Payload.Collection.Slug
	e.Chain = temp.Payload.Payload.Chain
	e.NftID = temp.Payload.Payload.Item.NftID
	e.TokenID = tokenID
	e.MakerAddress = temp.Payload.Payload.Maker.Address
	e.TakerAddress = temp.Payload.Payload.Taker.Address
	e.OrderHash = temp.Payload.Payload.OrderHash
	e.PriceWei = priceWei
	e.UsdPrice = temp.Payload.Payload.PaymentToken.UsdPrice
	e.ClosingDate = temp.Payload.Payload.ClosingDate
	e.TxHash = temp.Payload.Payload.Transaction.Hash

	return nil
}

type TraitCriterion struct {
	Name string `json:"trait_name"`
	Type string `json:"trait_type"`
}

type TraitOffer struct {
	EventType         string
	Slug              string
	Chain             string
	ContractAddr      string
	PriceWei          *big.Int
	MakerAddress      string
	OrderHash         string
	UsdPrice          string
	CreatedDate       string
	ExpirationDate    string
	EndTime           int64
	TraitCriteriaList []TraitCriterion
}

func (e *TraitOffer) UnmarshalJSON(data []byte) error {
	type msg struct {
		Payload struct {
			Payload struct {
				BasePrice  string `json:"base_price"`
				Chain      string `json:"chain"`
				Collection struct {
					Slug string `json:"slug"`
				} `json:"collection"`
				CreatedDate           string `json:"created_date"`
				ExpirationDate        string `json:"expiration_date"`
				AssetContractCriteria struct {
					Address string `json:"address"`
				} `json:"asset_contract_criteria"`
				Maker struct {
					Address string `json:"address"`
				} `json:"maker"`
				OrderHash    string `json:"order_hash"`
				PaymentToken struct {
					UsdPrice string `json:"usd_price"`
				} `json:"payment_token"`
				ProtocolData struct {
					Parameters struct {
						EndTime string `json:"endTime"`
					} `json:"parameters"`
				} `json:"protocol_data"`
				Quantity          int `json:"quantity"`
				TraitCriteriaList []struct {
					TraitName string `json:"trait_name"`
					TraitType string `json:"trait_type"`
				} `json:"trait_criteria_list"`
			} `json:"payload"`
		} `json:"payload"`
	}

	var temp msg
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	quantity := temp.Payload.Payload.Quantity
	if quantity != 1 {
		price, err := StringToBigInt(temp.Payload.Payload.BasePrice)
		if err != nil {
			return fmt.Errorf("unmarshalJSON trait offer;failed to convert base price to big.Int: %w", err)
		}
		e.PriceWei = CalculatePriceForQuantity(price, quantity)
	} else {
		price, err := StringToBigInt(temp.Payload.Payload.BasePrice)
		if err != nil {
			return fmt.Errorf("unmarshalJSON trait offer;failed to convert base price to big.Int: %w", err)
		}

		e.PriceWei = price
	}

	e.EventType = "trait_offer"
	e.Slug = temp.Payload.Payload.Collection.Slug
	e.Chain = temp.Payload.Payload.Chain
	e.ContractAddr = temp.Payload.Payload.AssetContractCriteria.Address
	e.MakerAddress = temp.Payload.Payload.Maker.Address
	e.OrderHash = temp.Payload.Payload.OrderHash
	e.UsdPrice = temp.Payload.Payload.PaymentToken.UsdPrice
	e.CreatedDate = temp.Payload.Payload.CreatedDate
	e.ExpirationDate = temp.Payload.Payload.ExpirationDate

	endTime, err := strconv.Atoi(temp.Payload.Payload.ProtocolData.Parameters.EndTime)
	if err != nil {
		return err
	}
	e.EndTime = int64(endTime)

	// Map trait_criteria_list
	e.TraitCriteriaList = make([]TraitCriterion, 0, len(temp.Payload.Payload.TraitCriteriaList))
	for _, tc := range temp.Payload.Payload.TraitCriteriaList {
		e.TraitCriteriaList = append(e.TraitCriteriaList, TraitCriterion{
			Name: tc.TraitName,
			Type: tc.TraitType,
		})
	}

	return nil
}

//TODO - add more events if needed
// ? ITEM_RECEIVED_BID = "item_received_bid",

// ORDER_INVALIDATE = "order_invalidate",
// ORDER_REVALIDATE = "order_revalidate",
