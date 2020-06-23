package feecalculator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	fee "github.com/iov-one/iovns/x/fee/types"
)

type FeeCalculator struct {
	fee.Calculator
	account types.Account
	domain  types.Domain

	ctx sdk.Context
	k   keeper.Keeper
}

var _ fee.Calculator = (*FeeCalculator)(nil)

func NewFeeCalculator(ctx sdk.Context, k keeper.Keeper) *FeeCalculator {
	return &FeeCalculator{
		ctx: ctx,
		k:   k,
	}
}

func (fc *FeeCalculator) WithDomain(domain types.Domain) *FeeCalculator {
	fc.domain = domain
	return fc
}

func (fc *FeeCalculator) WithAccount(account types.Account) *FeeCalculator {
	fc.account = account
	return fc
}

func (fc FeeCalculator) CalculateFee(msg fee.ProductMsg) (sdk.Coin, error) {
	var f sdk.Dec
	var err error
	switch msg := msg.(type) {
	// domain handlers
	case *types.MsgRegisterDomain:
		f, err = fc.calculateRegisterDomain(msg)
	default:
		f, err = sdk.Dec{}, nil
	}

	if err != nil {
		return sdk.Coin{}, err
	}
	feeParams := fc.k.FeeKeeper.GetFeeParams(fc.ctx)
	// get current price
	currentPrice := feeParams.FeeCoinPrice
	toPay := currentPrice.Quo(f)
	var feeAmount sdk.Int
	// get fee amount
	feeAmount = toPay.TruncateInt()
	defaultFee := feeParams.FeeDefault
	// if expected fee is lower than default fee then set the default fee as current fee
	if feeAmount.LT(defaultFee.TruncateInt()) {
		feeAmount = defaultFee.TruncateInt()
	}
	// generate coins to pay
	coinsToPay := sdk.NewCoin(feeParams.FeeCoinDenom, feeAmount)
	return coinsToPay, nil
}

func (fc *FeeCalculator) calculateAddAccountCertificates(msg *types.MsgAddAccountCertificates) (sdk.Dec, error) {
	panic("implement me")
}

func (fc *FeeCalculator) calculateRegisterDomain(msg *types.MsgRegisterDomain) (sdk.Dec, error) {
	panic("implement me")
}
