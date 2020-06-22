package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/keeper"
	fee "github.com/iov-one/iovns/x/fee/types"
)

type FeeCalculator struct {
	fee.Calculator
	account *Account
	domain  *Domain

	ctx sdk.Context
	k   keeper.Keeper
}

func NewFeeCalculator(ctx sdk.Context, k keeper.Keeper) *FeeCalculator {
	return &FeeCalculator{
		ctx: ctx,
		k:   k,
	}
}

func (fc *FeeCalculator) WithDomain(domain Domain) *FeeCalculator {
	fc.domain = &domain
	return fc
}

func (fc *FeeCalculator) WithAccount(account Account) *FeeCalculator {
	fc.account = &account
	return fc
}

func (fc FeeCalculator) CalculateFee(msg fee.ProductMsg) (sdk.Coin, error) {
	// calculate expected fee
	f, err := msg.CalculateFee(fc)
	if err != nil {
		return sdk.Coin{}, err
	}
	// get current price
	currentPrice := fc.k.FeeKeeper.GetFeeCoinPrice(fc.ctx)
	toPay := currentPrice.Quo(f)
	var feeAmount sdk.Int
	// get fee amount
	feeAmount = toPay.TruncateInt()
	defaultFee := fc.k.FeeKeeper.GetDefaultFee(fc.ctx)
	// if expected fee is lower than default fee then set the default fee as current fee
	if feeAmount.LT(defaultFee.TruncateInt()) {
		feeAmount = defaultFee.TruncateInt()
	}
	coinDenom := fc.k.FeeKeeper.GetFeeCoinDenom(fc.ctx)
	// generate coins to pay
	coinsToPay := sdk.NewCoin(coinDenom, feeAmount)
	return coinsToPay, nil
}

func (m *MsgAddAccountCertificates) CalculateFee(calc FeeCalculator) (sdk.Dec, error) {
}

func (m *MsgDeleteAccountCertificate) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgDeleteAccount) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgDeleteDomain) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgRegisterAccount) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

/* CONTRACT
Required fee seeds
- register_domain(1-6)
- register_open_domain_multiplier
*/
func (m *MsgRegisterDomain) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	var seedID string
	level := len(calculator.domain.Name)
	switch level {
	case 1, 2, 3, 4, 5, 6:
		seedID = snakeCaseAppend(m.Type(), string(level))
	default:
		seedID = snakeCaseAppend(m.Type(), "default")
	}

	feeSeed := calculator.k.FeeKeeper.GetFeeSeed(calculator.ctx, seedID)
	var fee sdk.Dec
	// if domain is open then we multiply
	if calculator.domain.Type == OpenDomain {
		multiplier := calculator.k.FeeKeeper.GetFeeSeed(calculator.ctx, "register_open_domain_multiplier")
		fee = sdk.Dec(feeSeed).Mul(sdk.Dec(multiplier))
	}
	return fee, nil
}

func (m *MsgRenewAccount) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgRenewDomain) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgReplaceAccountTargets) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgReplaceAccountMetadata) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgTransferAccount) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}

func (m *MsgTransferDomain) CalculateFee(calculator FeeCalculator) (sdk.Dec, error) {
	panic("implement me")
}
