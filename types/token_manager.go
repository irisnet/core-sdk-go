package types

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
)

var _ TokenManager = DefaultTokenManager{}

type DefaultTokenManager struct{}

func (TokenManager DefaultTokenManager) QueryToken(denom string) (Token, error) {
	return Token{}, nil
}

func (TokenManager DefaultTokenManager) SaveTokens(tokens ...Token) {
	return
}

func (TokenManager DefaultTokenManager) ToMinCoin(coins ...types.DecCoin) (types.Coins, Error) {
	for i := range coins {
		if coins[i].Denom == "iris" {
			coins[i].Denom = "uiris"
			coins[i].Amount = coins[i].Amount.MulInt(math.NewIntWithDecimal(1, 6))
		}
	}
	ucoins, _ := types.DecCoins(coins).TruncateDecimal()
	return ucoins, nil
}

func (TokenManager DefaultTokenManager) ToMainCoin(coins ...types.Coin) (types.DecCoins, Error) {
	decCoins := make(types.DecCoins, len(coins), 0)
	for _, coin := range coins {
		if coin.Denom == "uiris" {
			amtount := types.NewDecFromInt(coin.Amount).Mul(types.NewDecWithPrec(1, 6))
			decCoins = append(decCoins, types.NewDecCoinFromDec("iris", amtount))
		}
	}
	return decCoins, nil
}
