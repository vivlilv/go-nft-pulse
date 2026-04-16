package events_opensea

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func StringToBigInt(s string) (*big.Int, error) {
	value := new(big.Int)
	_, ok := value.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert string '%s' to big.Int", s)
	}
	return value, nil
}

func GetTokenID(s string) (int, error) {
	if len(s) == 0 {
		return -1, fmt.Errorf("empty string provided for token ID extraction")
	}
	parts := strings.Split(s, "/")
	tokenIDStr := parts[len(parts)-1]
	tokenID, err := strconv.Atoi(tokenIDStr)
	if err != nil {
		return -1, fmt.Errorf("failed to convert token ID '%s' to integer: %w", tokenIDStr, err)
	}
	return tokenID, nil
}

func CalculatePriceForQuantity(basePrice *big.Int, quantity int) *big.Int {
	divider := big.NewInt(int64(quantity))
	return new(big.Int).Div(basePrice, divider)
}
