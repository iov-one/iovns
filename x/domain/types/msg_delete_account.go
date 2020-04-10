package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgDeleteAccount struct {
	Domain string
	Name   string
	Owner  sdk.AccAddress
}

// Route returns the name of the module
func (m MsgDeleteAccount) Route() string {
	return RouterKey
}

// Type returns the action
func (m MsgDeleteAccount) Type() string {
	return "delete_account"
}

// ValidateBasic does stateless checks on the request
func (m MsgDeleteAccount) ValidateBasic() error {
	if m.Owner == nil {
		return sdkerrors.Wrap(ErrInvalidOwner, "owner is empty")
	}
	if m.Name == "" {
		return sdkerrors.Wrap(ErrInvalidAccountName, "account name is empty")
	}
	if m.Domain == "" {
		return sdkerrors.Wrap(ErrInvalidDomainName, "domain name is empty")
	}
	// success
	return nil
}

// GetSignBytes returns an ordered json of the request
func (m MsgDeleteAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners returns the list of address that should list the request
// in this case the admin of the domain
func (m MsgDeleteAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}
