package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/bnsd/x/account"
	"github.com/iov-one/bnsd/x/domain/internal/types"
	"regexp"
	"time"
)

// handleMsgRegisterDomain handles the domain registration process
func handleMsgRegisterDomain(ctx sdk.Context, keeper Keeper, msg types.MsgRegisterDomain) (resp *sdk.Result, err error) {
	// check if domain exists
	if _, ok := keeper.GetDomain(ctx, msg.Name); ok {
		err = sdkerrors.Wrap(types.ErrDomainAlreadyExists, msg.Name)
		return
	}
	// if domain does not exist then check if we can register it
	// check if name is valid based on the configuration saved in the state
	if !regexp.MustCompile(keeper.ConfKeeper.GetValidDomainRegexp(ctx)).MatchString(msg.Name) {
		err = sdkerrors.Wrap(types.ErrInvalidDomainName, msg.Name)
		return
	}
	// if domain has super user then admin must be configuration owner
	if !msg.HasSuperuser && !msg.Admin.Equals(keeper.ConfKeeper.GetOwner(ctx)) {
		err = sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to register a domain without a superuser", msg.Admin)
		return
	}
	// set new domain
	domain := types.Domain{
		Name:         msg.Name,
		Admin:        msg.Admin,
		ValidUntil:   ctx.BlockTime().Add(keeper.ConfKeeper.GetDomainRenewDuration(ctx)).Unix(),
		HasSuperuser: false,
		AccountRenew: time.Duration(msg.AccountRenew) * time.Second,
		Broker:       msg.Broker,
	}
	// save domain
	keeper.SetDomain(ctx, domain)
	// generate empty name account
	acc := account.Account{
		Domain:       msg.Name,
		Name:         "",
		Owner:        msg.Admin,
		ValidUntil:   ctx.BlockTime().Add(domain.AccountRenew).Unix(),
		Targets:      nil,
		Certificates: nil,
		Broker:       nil, // TODO ??
	}
	// save account
	keeper.AccKeeper.SetAccount(ctx, acc)
	// success TODO think here, can we emit any useful event
	return &sdk.Result{
		Data:   nil,
		Log:    "",
		Events: nil,
	}, nil
}
