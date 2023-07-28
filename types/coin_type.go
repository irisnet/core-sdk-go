package types

import (
	"strings"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types"
)

type Unit struct {
	Denom string `json:"denom"` //denom of unit
	Scale uint8  `json:"scale"` //scale of unit
}

func NewUnit(denom string, scale uint8) Unit {
	return Unit{
		Denom: denom,
		Scale: scale,
	}
}

// GetScaleFactor return 1 * 10^scale
func (u Unit) GetScaleFactor() math.Int {
	return math.NewIntWithDecimal(1, int(u.Scale))
}

type CoinType struct {
	Name     string `json:"name"`      //description name of CoinType
	MinUnit  Unit   `json:"min_unit"`  //the min unit of CoinType
	MainUnit Unit   `json:"main_unit"` //the max unit of CoinType
	Desc     string `json:"desc"`      //the description of CoinType
}

// ConvertToMainCoin return the main denom coin from args
func (ct CoinType) ConvertToMainCoin(coin types.Coin) (types.DecCoin, error) {
	if !ct.hasUnit(coin.Denom) {
		return types.DecCoin{
			Amount: types.NewDecFromInt(coin.Amount),
			Denom:  coin.Denom,
		}, nil
		//return DecCoin{}, errors.New("coinType unit (%s) not defined" + coin.Denom)
	}

	if ct.isMainUnit(coin.Denom) {
		return types.DecCoin{}, nil
	}

	// dest amount = src amount * (10^(dest scale) / 10^(src scale))
	dstScale := types.NewDecFromInt(ct.MainUnit.GetScaleFactor())
	srcScale := types.NewDecFromInt(ct.MinUnit.GetScaleFactor())
	amount := types.NewDecFromInt(coin.Amount)

	amt := amount.Mul(dstScale).Quo(srcScale)
	return types.NewDecCoinFromDec(ct.MainUnit.Denom, amt), nil
}

// ToMinCoin return the min denom coin from args
func (ct CoinType) ConvertToMinCoin(coin types.DecCoin) (newCoin types.Coin, err error) {
	if !ct.hasUnit(coin.Denom) {
		return types.Coin{
			Amount: coin.Amount.TruncateInt(),
			Denom:  coin.Denom,
		}, nil
		//return newCoin, errors.New("coinType unit (%s) not defined" + coin.Denom)
	}

	if ct.isMinUnit(coin.Denom) {
		newCoin, _ := coin.TruncateDecimal()
		return newCoin, nil
	}

	// dest amount = src amount * (10^(dest scale) / 10^(src scale))
	srcScale := types.NewDecFromInt(ct.MainUnit.GetScaleFactor())
	dstScale := types.NewDecFromInt(ct.MinUnit.GetScaleFactor())
	amount := coin.Amount

	amt := amount.Mul(dstScale).Quo(srcScale)
	return types.NewCoin(ct.MinUnit.Denom, amt.RoundInt()), nil
}

func (ct CoinType) isMainUnit(name string) bool {
	return ct.MainUnit.Denom == strings.TrimSpace(name)
}

func (ct CoinType) isMinUnit(name string) bool {
	return ct.MinUnit.Denom == strings.TrimSpace(name)
}

func (ct CoinType) hasUnit(name string) bool {
	return ct.isMainUnit(name) || ct.isMinUnit(name)
}
