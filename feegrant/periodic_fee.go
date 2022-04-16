package feegrant

import (
	"fmt"
	sdk "github.com/irisnet/core-sdk-go/types"
	"time"
)

var _ FeeAllowanceI = (*PeriodicAllowance)(nil)

func (m *PeriodicAllowance) Accept(ctx sdk.Context, fee sdk.Coins, _ []sdk.Msg) (remove bool, err error) {
	blockTime := ctx.BlockTime()

	if m.Basic.Expiration != nil && blockTime.After(*m.Basic.Expiration) {
		return true, sdk.Wrap(fmt.Errorf("fee allowance expired : %s", "absolute limit"))
	}

	m.tryResetPeriod(blockTime)

	// deduct from both the current period and the max amount
	var isNeg bool
	m.PeriodCanSpend, isNeg = m.PeriodCanSpend.SafeSub(fee)
	if isNeg {
		return false, sdk.Wrap(fmt.Errorf("fee limit exceeded : %s", "period limit"))
	}

	if m.Basic.SpendLimit != nil {
		m.Basic.SpendLimit, isNeg = m.Basic.SpendLimit.SafeSub(fee)
		if isNeg {
			return false, sdk.Wrap(fmt.Errorf("fee limit exceeded : %s", "absolute limit"))
		}

		return m.Basic.SpendLimit.IsZero(), nil
	}

	return false, nil
}

// tryResetPeriod will check if the PeriodReset has been hit. If not, it is a no-op.
// If we hit the reset period, it will top up the PeriodCanSpend amount to
// min(PeriodSpendLimit, Basic.SpendLimit) so it is never more than the maximum allowed.
// It will also update the PeriodReset. If we are within one Period, it will update from the
// last PeriodReset (eg. if you always do one tx per day, it will always reset the same time)
// If we are more then one period out (eg. no activity in a week), reset is one Period from the execution of this method
func (m *PeriodicAllowance) tryResetPeriod(blockTime time.Time) {
	if blockTime.Before(m.PeriodReset) {
		return
	}

	// set PeriodCanSpend to the lesser of Basic.SpendLimit and PeriodSpendLimit
	if _, isNeg := m.Basic.SpendLimit.SafeSub(m.PeriodSpendLimit); isNeg && !m.Basic.SpendLimit.Empty() {
		m.PeriodCanSpend = m.Basic.SpendLimit
	} else {
		m.PeriodCanSpend = m.PeriodSpendLimit
	}

	// If we are within the period, step from expiration (eg. if you always do one tx per day, it will always reset the same time)
	// If we are more then one period out (eg. no activity in m week), reset is one period from this time
	m.PeriodReset = m.PeriodReset.Add(m.Period)
	if blockTime.After(m.PeriodReset) {
		m.PeriodReset = blockTime.Add(m.Period)
	}
}

func (m *PeriodicAllowance) ValidateBasic() error {
	if err := m.Basic.ValidateBasic(); err != nil {
		return err
	}

	if !m.PeriodSpendLimit.IsValid() {
		return sdk.Wrap(fmt.Errorf("invalid coins , spend amount is invalid: %v", m.PeriodSpendLimit))
	}
	if !m.PeriodSpendLimit.IsAllPositive() {
		return sdk.Wrap(fmt.Errorf("invalid coins : %s", "spend limit must be positive"))
	}
	if !m.PeriodCanSpend.IsValid() {
		return sdk.Wrap(fmt.Errorf("invalid coins , can spend amount is invalid: %v", m.PeriodCanSpend))
	}
	// We allow 0 for CanSpend
	if m.PeriodCanSpend.IsAnyNegative() {
		return sdk.Wrap(fmt.Errorf("invalid coins : %s", "can spend must not be negative"))
	}

	// ensure PeriodSpendLimit can be subtracted from total (same coin types)
	if m.Basic.SpendLimit != nil && !m.PeriodSpendLimit.DenomsSubsetOf(m.Basic.SpendLimit) {
		return sdk.Wrap(fmt.Errorf("invalid coins : %s", "period spend limit has different currency than basic spend limit"))
	}

	// check times
	if m.Period.Seconds() < 0 {
		return sdk.Wrap(fmt.Errorf("invalid duration : %s", "negative clock step"))
	}

	return nil
}
