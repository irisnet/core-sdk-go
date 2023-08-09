package irishub

import (
	"cosmossdk.io/math"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/irisnet/core-sdk-go/types"
)

type TokenManager struct{}

func (TokenManager TokenManager) QueryToken(denom string) (sdktypes.Token, error) {
	return sdktypes.Token{}, nil
}

func (TokenManager TokenManager) SaveTokens(tokens ...sdktypes.Token) {
	return
}

func (TokenManager TokenManager) ToMinCoin(coins ...cosmostypes.DecCoin) (cosmostypes.Coins, sdktypes.Error) {
	for i := range coins {
		if coins[i].Denom == "iris" {
			coins[i].Denom = "uiris"
			coins[i].Amount = coins[i].Amount.MulInt(math.NewIntWithDecimal(1, 6))
		}
	}
	ucoins, _ := cosmostypes.DecCoins(coins).TruncateDecimal()
	return ucoins, nil
}

func (TokenManager TokenManager) ToMainCoin(coins ...cosmostypes.Coin) (cosmostypes.DecCoins, sdktypes.Error) {
	decCoins := make(cosmostypes.DecCoins, len(coins), 0)
	for _, coin := range coins {
		if coin.Denom == "uiris" {
			amtount := cosmostypes.NewDecFromInt(coin.Amount).Mul(cosmostypes.NewDecWithPrec(1, 6))
			decCoins = append(decCoins, cosmostypes.NewDecCoinFromDec("iris", amtount))
		}
	}
	return decCoins, nil
}
