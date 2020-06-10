package fees

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/types"
)

// Keeper defines the behaviour of the keeper, it's here just to avoid import cycling
type Keeper interface {
	GetAccountsInDomain(ctx sdk.Context, name string, do func(k []byte) bool)
}

// Controller defines the fee controller behaviour
// exists only in order to avoid devs creating a fee
// controller without using the constructor feunction
type Controller interface {
	GetFee(msg sdk.Msg) sdk.Coin
}

func NewController(ctx sdk.Context, k Keeper, fees *configuration.Fees, domain types.Domain) Controller {
	return feeApplier{
		moduleFees: *fees,
		ctx:        ctx,
		k:          k,
		domain:     domain,
	}
}

type feeApplier struct {
	moduleFees configuration.Fees
	ctx        sdk.Context
	k          Keeper
	domain     types.Domain
}

func (f feeApplier) registerDomain() sdk.Dec {
	var registerDomainFee sdk.Dec
	level := len(f.domain.Name)
	switch level {
	case 1:
		registerDomainFee = f.moduleFees.RegisterDomain1
	case 2:
		registerDomainFee = f.moduleFees.RegisterDomain2
	case 3:
		registerDomainFee = f.moduleFees.RegisterDomain3
	case 4:
		registerDomainFee = f.moduleFees.RegisterDomain4
	case 5:
		registerDomainFee = f.moduleFees.RegisterDomain5
	default:
		registerDomainFee = f.moduleFees.RegisterDomainDefault
	}
	// if domain is open then we multiply
	if f.domain.Type == types.OpenDomain {
		registerDomainFee.Mul(f.moduleFees.RegisterOpenDomainMultiplier)
	}
	return registerDomainFee
}

func (f feeApplier) transferDomain() sdk.Dec {
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.TransferDomainOpen
	case types.ClosedDomain:
		return f.moduleFees.TransferDomainClosed
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) renewDomain() sdk.Dec {
	if f.domain.Type == types.OpenDomain {
		return f.moduleFees.RenewOpenDomain
	}
	var accountN int64
	f.k.GetAccountsInDomain(f.ctx, f.domain.Name, func(_ []byte) bool {
		accountN++
		return true
	})
	fee := f.moduleFees.RegisterClosedAccount
	fee.MulInt64(accountN)
	return fee
}

func (f feeApplier) registerAccount() sdk.Dec {
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.RegisterOpenAccount
	case types.ClosedDomain:
		return f.moduleFees.RegisterClosedAccount
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) transferAccount() sdk.Dec {
	switch f.domain.Type {
	case types.ClosedDomain:
		return f.moduleFees.TransferClosedAccount
	case types.OpenDomain:
		return f.moduleFees.TransferOpenAccount
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) renewAccount() sdk.Dec {
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.RegisterOpenAccount
	case types.ClosedDomain:
		return f.moduleFees.RegisterClosedAccount
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) replaceTargets() sdk.Dec {
	return f.moduleFees.ReplaceAccountTargets
}

func (f feeApplier) addCert() sdk.Dec {
	return f.moduleFees.AddAccountCertificate
}

func (f feeApplier) delCert() sdk.Dec {
	return f.moduleFees.DelAccountCertificate
}

func (f feeApplier) setMetadata() sdk.Dec {
	return f.moduleFees.SetAccountMetadata
}

func (f feeApplier) defaultFee() sdk.Dec {
	return f.moduleFees.DefaultFee
}

func (f feeApplier) getFeeParam(msg sdk.Msg) sdk.Dec {
	switch msg.(type) {
	case *types.MsgTransferDomain:
		return f.transferDomain()
	case *types.MsgRegisterDomain:
		return f.registerDomain()
	case *types.MsgRenewDomain:
		return f.renewDomain()
	case *types.MsgRegisterAccount:
		return f.registerAccount()
	case *types.MsgTransferAccount:
		return f.transferAccount()
	case *types.MsgRenewAccount:
		return f.renewAccount()
	case *types.MsgReplaceAccountTargets:
		return f.replaceTargets()
	case *types.MsgDeleteAccountCertificate:
		return f.delCert()
	case *types.MsgReplaceAccountMetadata:
		return f.setMetadata()
	default:
		return f.defaultFee()
	}
}

// GetFee returns a fee based on the provided message
func (f feeApplier) GetFee(msg sdk.Msg) sdk.Coin {
	// get current price
	currentPrice := f.moduleFees.FeeCoinPrice
	// get coin denom
	coinDenom := f.moduleFees.FeeCoinDenom
	// get fee parameter
	param := f.getFeeParam(msg)
	// calculate the quotient between current price and parameter
	toPay := currentPrice.Quo(param)
	var feeAmount sdk.Int
	// get fee amount
	feeAmount = toPay.TruncateInt()
	// if expected fee is lower than default fee then set the default fee as current fee
	if feeAmount.LT(f.moduleFees.DefaultFee.TruncateInt()) {
		feeAmount = f.moduleFees.DefaultFee.TruncateInt()
	}
	// generate coins to pay
	coinsToPay := sdk.NewCoin(coinDenom, feeAmount)
	return coinsToPay
}
