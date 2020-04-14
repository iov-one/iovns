package domain

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgDeleteAccountCertificate(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteAccountCertificate) (*sdk.Result, error) {
	// get account
	account, exists := k.GetAccount(ctx, iovns.GetAccountKey(msg.Domain, msg.Name))
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found: %s", msg.Name)
	}
	// check if signer is account owner
	if !msg.Owner.Equals(account.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s cannot delete certificates from account owned by %s", msg.Owner, account.Owner)
	}
	// check if certificate exists
	var found bool
	var certIndex int
	// iterate over certs
	for i, cert := range account.Certificates {
		// if found
		if bytes.Equal(cert, msg.DeleteCertificate) {
			found = true  // set found to true
			certIndex = i // save index of cert for removal
			break
		}
	}
	// check if found
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCertificateDoesNotExist, "not found")
	}
	// remove it
	account.Certificates = append(account.Certificates[:certIndex], account.Certificates[certIndex+1:]...)
	// update account
	k.SetAccount(ctx, account)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}
