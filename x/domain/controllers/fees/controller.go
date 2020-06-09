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

func (f feeApplier) registerDomain() sdk.Coin {
	var registerDomainFee sdk.Coin
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
		registerDomainFee.Amount.Mul(f.moduleFees.RegisterOpenDomainMultiplier)
	}
	return registerDomainFee
}

func (f feeApplier) transferDomain() sdk.Coin {
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.TransferDomainOpen
	case types.ClosedDomain:
		return f.moduleFees.TransferDomainClosed
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) renewDomain() sdk.Coin {
	if f.domain.Type == types.OpenDomain {
		return f.moduleFees.RenewOpenDomain
	}
	var accountN int64
	f.k.GetAccountsInDomain(f.ctx, f.domain.Name, func(_ []byte) bool {
		accountN++
		return true
	})
	fee := f.moduleFees.RegisterClosedAccount
	fee.Amount.MulRaw(accountN)
	return fee
}

func (f feeApplier) registerAccount() sdk.Coin {
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.RegisterOpenAccount
	case types.ClosedDomain:
		return f.moduleFees.RegisterClosedAccount
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) transferAccount() sdk.Coin {
	switch f.domain.Type {
	case types.ClosedDomain:
		return f.moduleFees.TransferClosedAccount
	case types.OpenDomain:
		return f.moduleFees.TransferOpenAccount
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) renewAccount() sdk.Coin {
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.RegisterOpenAccount
	case types.ClosedDomain:
		return f.moduleFees.RegisterClosedAccount
	}
	return f.moduleFees.DefaultFee
}

func (f feeApplier) replaceTargets() sdk.Coin {
	return f.moduleFees.ReplaceAccountTargets
}

func (f feeApplier) addCert() sdk.Coin {
	return f.moduleFees.AddAccountCertificate
}

func (f feeApplier) delCert() sdk.Coin {
	return f.moduleFees.DelAccountCertificate
}

func (f feeApplier) setMetadata() sdk.Coin {
	return f.moduleFees.SetAccountMetadata
}

func (f feeApplier) defaultFee() sdk.Coin {
	return f.moduleFees.DefaultFee
}

// GetFee returns a fee based on the provided message
func (f feeApplier) GetFee(msg sdk.Msg) sdk.Coin {
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
